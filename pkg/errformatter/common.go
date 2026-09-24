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

// errorCodeContainable defines the error-code-aware operations.
//
// Methods in this interface attach codes to new errors, extract codes from
// wrapped errors, or test whether an error carries one of several given codes.
type errorCodeContainable interface {
	// ErrorWithCode creates a new formatted error wrapping err that carries
	// the given positive code. If code <= 0 the call panics.
	ErrorWithCode(err error, code int) error
	// ErrWithCode is a short alias for ErrorWithCode.
	ErrWithCode(err error, code int) error
	// NewErrorWithCode creates a new error with text and code, no wrapper.
	NewErrorWithCode(text string, code int) error
	// ErrorGetCode extracts the embedded error code from err. Returns -1
	// when the error does not carry a code.
	ErrorGetCode(err error) int
	// ErrGetCode is a short alias for ErrorGetCode.
	ErrGetCode(err error) int
	// ErrorCodeIsOneOf reports whether err carries a code present in codes.
	// Returns the matched code (or -1).
	ErrorCodeIsOneOf(err error, codes ...int) (int, bool)
	// ErrCodeIsOneOf is a short alias for ErrorCodeIsOneOf.
	ErrCodeIsOneOf(err error, codes ...int) (int, bool)
}

// ErrorFormatterService formats errors into strings while preserving wrap
// semantics and optional enriched data (details, scope, codes, public codes).
//
// Every concrete service implements this full interface. Use the constructor
// functions or the ErrorFormatterBuidler to obtain an instance:
//
//	svc := errfmt.NewErrorBasicFormatter()            // simply formatted
//	svc := errfmt.NewScopedErrorFormatter("wallet")   // scope-prefixed
//	svc := errfmt.NewValuesErrorFormatter(vs...)      // value-based
//
// ErrorNoWrap variants return the formatted message as an error without
// wrapping the original error value. All other methods perform pseudo-wrapping
// (concatenating a formatted header before the wrapped error).
type ErrorFormatterService interface {
	errorCodeContainable

	// ErrorNoWrap wraps err with formatting details but does NOT chain
	// err as Inner() — callers who need a clean output use this method
	// to avoid linter warnings from unhandled wrap chains.
	// Implements ErrorFormatterService.ErrorNoWrap.
	ErrorNoWrap(err error) error

	// ErrNoWrap is a short alias for ErrorNoWrap.
	// Implements ErrorFormatterService.ErrNoWrap.
	ErrNoWrap(err error) error

	// ErrorOnly wraps err with the provided detail strings.
	// Implements ErrorFormatterService.ErrorOnly.
	ErrorOnly(err error, details ...string) error

	// Error wraps err with details, applying additional formatting specific
	// to the concrete service strategy.
	// Implements ErrorFormatterService.Error.
	Error(err error, details ...string) error

	// Errorf formats err with a printf-style string and arguments.
	// Implements ErrorFormatterService.Errorf.
	Errorf(err error, format string, args ...interface{}) error

	// NewError creates a new error (no wrapper) from detail strings.
	// Implements ErrorFormatterService.NewError.
	NewError(details ...string) error

	// NewErrorf creates a new error (no wrapper) from a format string.
	// Implements ErrorFormatterService.NewErrorf.
	NewErrorf(format string, args ...interface{}) error
}

// ErrorFormatterBuidler constructs configured ErrorFormatterService instances.
//
// It offers four construction paths, each producing a different formatting strategy:
//
//	MakeScoped("wallet")     → scope-prefixed formatter
//	MakeValued(vs...)        → value-based formatter (with values or defaults)
//	MakeSimply()             → basic formatter delegating to stdlib
//	Make()                   → default / simply-formatted (same as MakeSimply)
//
// Note: the interface name contains a historical typo ("Buidler" instead of "Builder").
// It is preserved for backward compatibility. See pkg/errformatter/doc.go for usage examples.
type ErrorFormatterBuidler interface {
	// MakeScoped returns a scope-prefixed formatter.
	// Implements ErrorFormatterBuidler.MakeScoped.
	MakeScoped(scope string) ErrorFormatterService

	// MakeValued returns a value-based formatter initialized with vals.
	// If vals is empty the result is equivalent to a plain valued formatter.
	// Implements ErrorFormatterBuidler.MakeValued.
	MakeValued(values ...Value) ErrorFormatterService

	// MakeSimply returns the basic formatter that delegates to standard library calls.
	// Implements ErrorFormatterBuidler.MakeSimply.
	MakeSimply() ErrorFormatterService

	// Make returns the default formatter (equivalent to MakeSimply).
	// Implements ErrorFormatterBuidler.Make.
	Make() ErrorFormatterService
}
