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

type ErrorScoped struct {
	Err   error
	scope string
}

// Error to string converter
func (e ErrorScoped) Error() string {
	return fmt.Sprintf("%s", e.Err.Error())
}

// Unwrap returns previous error
func (e ErrorScoped) Unwrap() error {
	return errors.Unwrap(e.Err)
}

// ScopedErrorOnly combines given error with details, WITHOUT function name...
func ScopedErrorOnly(err error, scope string, details ...string) *ErrorScoped {
	if err == nil {
		return nil
	}

	if len(details) == 0 {
		return &ErrorScoped{
			scope: scope,
			Err:   fmt.Errorf("%s: %w", scope, err),
		}
	}

	return &ErrorScoped{
		scope: scope,
		Err:   fmt.Errorf("%s: %w -> %s", scope, err, strings.Join(details, ", ")),
	}
}

// ScopedError combines given error with details and finishes with caller func name...
func ScopedError(err error, scope string, details ...string) *ErrorScoped {
	return ScopedErrorOnly(err, scope, append(details, getFuncName())...)
}

// NewScopedError returns error by combining given details and finishes with caller func name...
//
//nolint:err113
func NewScopedError(scope string, details ...string) *ErrorScoped {
	return &ErrorScoped{
		Err: fmt.Errorf("%s: %s", scope,
			strings.Join(append(details, getFuncName()), ", ")),
		scope: scope,
	}
}

// NewScopedErrorf returns error by combining given details and finishes with caller func name, printf formatting...
//
//nolint:err113
func NewScopedErrorf(format string, scope string, args ...interface{}) *ErrorScoped {
	return &ErrorScoped{
		Err: fmt.Errorf(
			"%s: %s", scope,
			strings.Join(append([]string{fmt.Sprintf(format, args...)}, getFuncName()), ", "),
		),
		scope: scope,
	}
}

// ScopedErrorf combines given error with details and finishes with caller func name, printf formatting...
func ScopedErrorf(err error, scope string,
	format string,
	args ...interface{},
) *ErrorScoped {
	return ScopedErrorOnly(err, scope, fmt.Sprintf(format, args...), getFuncName())
}
