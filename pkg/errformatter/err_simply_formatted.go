// MIT NON-AI License
//
// Copyright (c) 2022-2026 Aleksei Kotelnikov(gudron2s@gmail.com)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
// In addition, the following restrictions apply:
//
// 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
// including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
// modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
//
// 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
// including but not limited to artificial intelligence, natural language processing, or data mining.
//
// 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
// for any damages resulting from such use.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package errformatter

import (
	"fmt"
	"runtime"
	"strings"
)

// ErrorNoWrap returns err unchanged. If err is nil it returns nil.
// No wrapping or formatting is applied.
func ErrorNoWrap(err error) error {
	if err == nil {
		return nil
	}

	return err
}

// ErrorOnly combines err with detail strings using "%w -> detail1, detail2" format.
// If err is nil it returns nil; if details is empty it returns err as-is.
func ErrorOnly(err error, details ...string) error {
	if err == nil {
		return nil
	}

	if len(details) == 0 {
		return err
	}

	return fmt.Errorf("%w -> %s", err, strings.Join(details, ", "))
}

// Error wraps err with detail strings and finishes with the caller function name.
// Currently identical to ErrorOnly — caller-location injection is not yet implemented.
func Error(err error, details ...string) error {
	return ErrorOnly(err, details...)
}

// NewError creates a new error from detail strings joined by ", " separator.
// The returned error does not wrap an existing error.
//
//nolint:err113 // direct error creation is intentional here.
func NewError(details ...string) error {
	return fmt.Errorf("%s", strings.Join(details, ", "))
}

// NewErrorf creates a new error from a printf-style format string and arguments.
// The formatted result is wrapped in a slice before joining by ", ".
// The returned error does not wrap an existing error.
//
//nolint:err113 // direct error creation is intentional here.
func NewErrorf(format string, args ...any) error {
	return fmt.Errorf(
		"%s",
		strings.Join([]string{fmt.Sprintf(format, args...)}, ", "),
	)
}

// Errorf wraps err with a printf-formatted detail string using "%w -> detail" format.
func Errorf(err error, format string, args ...any) error {
	return ErrorOnly(err, fmt.Sprintf(format, args...))
}

// getFuncName returns the calling function's package and name via reflection.
// Used internally to inject caller context into error messages. Not exported deliberately.
//
//nolint:unused // planned for future caller-location injection.
func getFuncName() string {
	pc := CallerStackSkip
	if runtime.FuncForPC(uintptr(pc)) == nil {
		return ""
	}

	name := runtime.FuncForPC(uintptr(pc)).Name()
	parts := strings.Split(name, ".")

	//nolint:mnd // minimum parts to split function name by "."
	if len(parts) < 2 {
		return ""
	}

	return fmt.Sprintf("[%s.%s]", parts[len(parts)-2], parts[len(parts)-1])
}
