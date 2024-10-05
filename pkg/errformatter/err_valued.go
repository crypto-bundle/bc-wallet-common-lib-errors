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
	values  [MaxKindValue]Value
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

func (e *valuedError) ScopeIs(scope string) bool {
	return e.scopeIsEqualWith(scope)
}

func (e *valuedError) scopeIsEqualWith(scope string) bool {
	if !e.settled.Has(ValueScopeIsSet) {
		return false
	}

	currentScope := e.values[KindScope].getScope()

	return currentScope == scope
}

func (e *valuedError) getScope() string {
	if !e.settled.Has(ValueScopeIsSet) {
		return ""
	}

	return e.values[KindScope].getScope()
}

func (e *valuedError) getCode() int {
	if !e.settled.Has(ValueCodeIsSet) {
		return ValueCodeMissing
	}

	return e.values[KindCode].getCode()
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
		scope := e.values[KindScope].getScope()
		details := e.values[KindDetails].getDetails()

		e.Err = fmt.Errorf("%s: %w -> %s", scope, err, strings.Join(details, ", "))

	case e.settled.Has(ValueDetailsIsSet) && !e.settled.Has(ValueScopeIsSet):
		details := e.values[KindDetails].getDetails()

		e.Err = fmt.Errorf("%w -> %s", err, strings.Join(details, ", "))

	case !e.settled.Has(ValueDetailsIsSet) && e.settled.Has(ValueScopeIsSet):
		scope := e.values[KindScope].getScope()

		e.Err = fmt.Errorf("%s: %w", scope, err)
	default:
		e.Err = fmt.Errorf("%w", err)
	}

	return e
}

func (e *valuedError) reWrap(value Value) *valuedError {
	if value.Kind() != KindScope {
		return e.setValue(value).setError(e.Err)
	}

	if !e.scopeIsEqualWith(value.getScope()) {
		return e.setValue(value).setError(e.Err)
	}

	return e
}

func (e *valuedError) reWrapByValues(values ...Value) *valuedError {
	if !e.settled.Has(ValueScopeIsSet) {
		return e.setValues(values...).setError(e.Err)
	}

	var scopeValue Value

	for i := range values {
		value := values[i]

		if value.Kind() != KindScope {
			_ = e.setValue(value)

			continue
		}

		scopeValue = value
	}

	if scopeValue.Kind() != KindScope {
		return e.setError(e.Err)
	}

	if !e.scopeIsEqualWith(scopeValue.getScope()) { // in case of re-wrap with another scope
		return e.setValue(scopeValue).setError(e.Err)
	}

	return e.setError(e.Unwrap())
}

func ValuedErrorGetCode(err error) int {
	var vErr *valuedError

	if !errors.As(err, &vErr) {
		return ValueCodeMissing
	}

	return vErr.getCode()
}

// ValuedErrorOnly combines given error with given Value, all Value type values must contain pre-reserved Kind...
func ValuedErrorOnly(err error, value Value) *valuedError {
	if err == nil {
		return nil
	}

	var vErr *valuedError
	if errors.As(err, &vErr) {
		return vErr.reWrap(value)
	}

	vErr = &valuedError{}

	return vErr.setValue(value).setError(err)
}

// MultiValuedErrorOnly combines given error with given Value list, all Value type values must contain pre-reserved Kind...
func MultiValuedErrorOnly(err error, value ...Value) *valuedError {
	if err == nil {
		return nil
	}

	var vErr *valuedError
	if errors.As(err, &vErr) {
		return vErr.reWrapByValues(value...)
	}

	vErr = &valuedError{}

	return vErr.setValues(value...).setError(err)
}

// ValuedError combines given error with details and finishes with caller func name, printf formatting...
func ValuedError(err error, values []Value, details ...string) *valuedError {
	values = append(values, NewValue(KindDetails, details))

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

	var vErr *valuedError
	if errors.As(err, &vErr) {
		vErr.Err = ErrorOnly(vErr.Err, fmt.Sprintf(format, args...))

		return vErr.setValues(values...)
	}

	vErr = &valuedError{
		Err:     ErrorOnly(err, fmt.Sprintf(format, args...)),
		values:  [4]Value{},
		settled: 0,
	}

	return vErr.setValues(values...)
}

// ValuedNewError combines given error with details and finishes with caller func name, printf formatting...
//
//nolint:err113
func ValuedNewError(values []Value, details ...string) *valuedError {
	var vErr valuedError

	newErr := fmt.Errorf("%s", strings.Join(details, ", "))

	return vErr.setValues(values...).setError(newErr)
}

// ValuedNewErrorf combines given error with details and finishes with caller func name, printf formatting...
//
//nolint:err113
func ValuedNewErrorf(values []Value, format string, args ...interface{}) *valuedError {
	var vErr valuedError

	newErr := fmt.Errorf("%s",
		strings.Join([]string{fmt.Sprintf(format, args...)}, ", "),
	)

	return vErr.setValues(values...).setError(newErr)
}
