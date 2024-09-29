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

type ValuedError struct {
	Err     error
	values  [...]Value
	settled Bits
}

// Error to string converter...
func (e ValuedError) Error() string {
	return e.Err.Error()
}

// Unwrap returns previous error...
func (e ValuedError) Unwrap() error {
	return errors.Unwrap(e.Err)
}

func (e *ValuedError) getCode() int {
	if !e.settled.Has(ValueCodeIsSet) {
		return ValueMissing
	}

	code, isCasted := e.values[KindCode].any.(int)
	if !isCasted {
		return ValueMissing
	}

	return code
}

func (e *ValuedError) setValues(values ...Value) *ValuedError {
	for i := range values {
		_ = e.setValue(values[i])
	}

	return e
}

func (e *ValuedError) setValue(value Value) *ValuedError {
	switch {
	case value.num == KindDetails && e.settled.Has(ValueScopeIsSet):
		scope := e.values[KindScope].any.(string)
		details := value.any.([]string)

		e.values[value.num] = value
		e.settled.Set(ValueDetailsIsSet)
		e.Err = fmt.Errorf("%s: %w -> %s", scope, e, strings.Join(details, ", "))

	case value.num == KindDetails && !e.settled.Has(ValueScopeIsSet):
		details := value.any.([]string)

		e.values[value.num] = value
		e.settled.Set(ValueDetailsIsSet)
		e.Err = fmt.Errorf("%w -> %s", e, strings.Join(details, ", "))

	case value.num == KindScope && e.settled.Has(ValueDetailsIsSet):
		details := e.values[KindDetails].any.([]string)
		scope := value.any.(string)

		e.values[value.num] = value
		e.settled.Set(ValueScopeIsSet)
		e.Err = fmt.Errorf("%s: %w -> %s", scope, e, strings.Join(details, ", "))

	case value.num == KindScope && !e.settled.Has(ValueDetailsIsSet):
		scope := value.any.(string)

		e.values[value.num] = value
		e.settled.Set(ValueScopeIsSet)
		e.Err = fmt.Errorf("%s: %w", scope, e)

	case value.num == KindCode:
		e.settled.Set(ValueCodeIsSet)
		e.values[KindCode] = value

	case value.num == KindPublicCode:
		e.settled.Set(ValuePublicCodeIsSet)
		e.values[KindPublicCode] = value
	}

	return e
}

func ValuedErrorGetCode(err error) int {
	var vErr ValuedError
	if !errors.As(err, &vErr) {
		return ValueMissing
	}

	return vErr.getCode()
}

// ValuedErrorOnly combines given error with given Value, all Value type values must contain pre-reserved Kind...
func ValuedErrorOnly(err error, value Value) *ValuedError {
	if err == nil {
		return nil
	}

	var vErr ValuedError
	if errors.As(err, &vErr) {
		return vErr.setValue(value)
	}

	vErr.values = make([...]Value, 0, len(kindStrings)-1)
	vErr.Err = nil
	vErr.settled = 0

	return vErr.setValue(value)
}

// MultiValuedErrorOnly combines given error with given Value list, all Value type values must contain pre-reserved Kind...
func MultiValuedErrorOnly(err error, value ...Value) *ValuedError {
	if err == nil {
		return nil
	}

	var vErr ValuedError
	if errors.As(err, &vErr) {
		return vErr.setValues(value...)
	}

	vErr.values = make([...]Value, 0, len(kindStrings)-1)
	vErr.Err = nil
	vErr.settled = 0

	return vErr.setValues(value...)
}
