/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package errformatter

import (
	"errors"
	"fmt"
	"strings"
)

type valuedError struct {
	Err     error
	values  []Value
	settled Bits
}

// Error to string converter...
func (e valuedError) Error() string {
	return e.Err.Error()
}

// Unwrap returns previous error...
func (e valuedError) Unwrap() error {
	return errors.Unwrap(e.Err)
}

func (e *valuedError) getCode() int {
	if !e.settled.Has(ValueCodeIsSet) {
		return ValueMissing
	}

	code, isCasted := e.values[KindCode].any.(int)
	if !isCasted {
		return ValueMissing
	}

	return code
}

func (e *valuedError) setValues(values ...Value) *valuedError {
	for i := range values {
		_ = e.setValue(values[i])
	}

	return e
}

func (e *valuedError) setValue(value Value) *valuedError {
	e.values[value.num] = value
	e.settled.Set(value.num.Bits())

	return e
}

func (e *valuedError) setError(err error) *valuedError {
	switch {
	case e.settled.Has(ValueDetailsIsSet) && e.settled.Has(ValueScopeIsSet):
		scope := e.values[KindScope].any.(string)
		details := e.values[KindDetails].any.([]string)

		e.Err = fmt.Errorf("%s: %w -> %s", scope, err, strings.Join(details, ", "))

	case e.settled.Has(ValueDetailsIsSet) && !e.settled.Has(ValueScopeIsSet):
		details := e.values[KindDetails].any.([]string)

		e.Err = fmt.Errorf("%w -> %s", err, strings.Join(details, ", "))

	case !e.settled.Has(ValueDetailsIsSet) && e.settled.Has(ValueScopeIsSet):
		scope := e.values[KindScope].any.(string)

		e.Err = fmt.Errorf("%s: %w", scope, err)
	}

	return e
}

func ValuedErrorGetCode(err error) int {
	var vErr valuedError
	if !errors.As(err, &vErr) {
		return ValueMissing
	}

	return vErr.getCode()
}

// ValuedErrorOnly combines given error with given Value, all Value type values must contain pre-reserved Kind...
func ValuedErrorOnly(err error, value Value) *valuedError {
	if err == nil {
		return nil
	}

	var vErr valuedError
	if errors.As(err, &vErr) {
		return vErr.setValue(value)
	}

	vErr.values = make([]Value, kindCount)
	vErr.Err = nil
	vErr.settled = 0

	return vErr.setValue(value).setError(err)
}

// MultiValuedErrorOnly combines given error with given Value list, all Value type values must contain pre-reserved Kind...
func MultiValuedErrorOnly(err error, value ...Value) *valuedError {
	if err == nil {
		return nil
	}

	var vErr valuedError
	if errors.As(err, &vErr) {
		return vErr.setValues(value...)
	}

	vErr.values = make([]Value, kindCount)
	vErr.Err = nil
	vErr.settled = 0

	return vErr.setValues(value...).setError(err)
}

// ValuedError combines given error with details and finishes with caller func name, printf formatting...
func ValuedError(err error, values []Value, details ...string) *valuedError {
	values = append(values, Value{
		num: KindDetails,
		any: append(details, getFuncName()),
	})

	return MultiValuedErrorOnly(err, values...).setError(err)
}

// ValuedErrorf combines given error with details and finishes with caller func name, printf formatting...
func ValuedErrorf(err error,
	values []Value,
	format string,
	args ...interface{},
) *valuedError {
	if err == nil {
		return nil
	}

	var vErr = valuedError{
		Err:     ErrorOnly(err, fmt.Sprintf(format, args...), getFuncName()),
		values:  make([]Value, 0, len(kindStrings)-1),
		settled: 0,
	}

	if errors.As(err, &vErr) {
		_ = vErr.setValues(values...)
	}

	return &vErr
}

// ValuedNewError combines given error with details and finishes with caller func name, printf formatting...
func ValuedNewError(values []Value, details ...string) *valuedError {
	details = append(details, getFuncName())

	vErr := &valuedError{
		Err:     fmt.Errorf("%s", strings.Join(details, ", ")),
		values:  make([]Value, kindCount),
		settled: 0,
	}

	return vErr.setValues(append(values, Value{
		num: KindDetails,
		any: details,
	})...)
}

// ValuedNewErrorf combines given error with details and finishes with caller func name, printf formatting...
func ValuedNewErrorf(values []Value, format string, args ...interface{}) *valuedError {
	vErr := &valuedError{
		Err: fmt.Errorf("%s",
			strings.Join(append([]string{fmt.Sprintf(format, args...)}, getFuncName()), ", "),
		),
		values:  make([]Value, kindCount),
		settled: 0,
	}

	return vErr.setValues(values...)
}
