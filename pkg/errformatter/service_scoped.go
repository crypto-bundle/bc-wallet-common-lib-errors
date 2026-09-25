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

// serviceScoped formats errors with a leading "[scope] " prefix.
// The scope string is set at construction time and is attached to every
// new error produced through scoped helpers via MultiValuedErrorOnly
// and ScopedError* convenience functions.
type serviceScoped struct {
	scope string
}

// ErrCodeIsOneOf is a short alias for ErrorCodeIsOneOf.
// Implements ErrorFormatterService.ErrCodeIsOneOf.
func (s *serviceScoped) ErrCodeIsOneOf(err error, codes ...int) (int, bool) {
	return s.ErrorCodeIsOneOf(err, codes...)
}

// ErrorCodeIsOneOf reports whether err carries a code present in codes.
// Returns the matched code (or -1).
// Implements ErrorFormatterService.ErrorCodeIsOneOf.
func (s *serviceScoped) ErrorCodeIsOneOf(err error, codes ...int) (int, bool) {
	errCode := s.ErrorGetCode(err)
	if errCode == -1 {
		return -1, false
	}

	for _, targetCode := range codes {
		if targetCode == errCode {
			return targetCode, true
		}
	}

	return -1, false
}

// NewErrorWithCode creates a new error carrying details and code,
// no wrapper. Panics if code <= 0.
// Implements ErrorFormatterService.NewErrorWithCode.
func (s *serviceScoped) NewErrorWithCode(text string, code int) error {
	return ValuedNewError([]Value{
		NewValue(KindDetails, text),
		NewValue(KindCode, code),
		NewValue(KindScope, s.scope),
	})
}

// ErrGetCode is a short alias for ErrorGetCode.
// Implements ErrorFormatterService.ErrGetCode.
func (s *serviceScoped) ErrGetCode(err error) int {
	return s.ErrorGetCode(err)
}

// ErrorGetCode extracts the embedded error code from err. Returns -1
// when the error does not carry a code.
// Implements ErrorFormatterService.ErrorGetCode.
func (s *serviceScoped) ErrorGetCode(err error) int {
	return ValuedErrorGetCode(err)
}

// ErrWithCode is a short alias for ErrorWithCode.
// Implements ErrorFormatterService.ErrWithCode.
func (s *serviceScoped) ErrWithCode(err error, code int) error {
	return s.ErrorWithCode(err, code)
}

// ErrorWithCode creates an error wrapping err that carries the given
// positive code. If code <= 0 the call panics.
// Implements ErrorFormatterService.ErrorWithCode.
func (s *serviceScoped) ErrorWithCode(err error, code int) error {
	if code <= 0 {
		panic("errfmt: code must be positive value")
	}

	return MultiValuedErrorOnly(err,
		NewValue(KindCode, code),
		NewValue(KindScope, s.scope))
}

// ErrNoWrap is a short alias for ErrorNoWrap.
// Implements ErrorFormatterService.ErrNoWrap.
func (s *serviceScoped) ErrNoWrap(err error) error {
	return ErrorNoWrap(err)
}

// ErrorNoWrap wraps err with formatting details but does NOT chain err
// as Inner(). Callers who need a clean output use this method.
// Implements ErrorFormatterService.ErrorNoWrap.
func (s *serviceScoped) ErrorNoWrap(err error) error {
	return ErrorNoWrap(err)
}

// ErrorOnly wraps err with the provided detail strings.
// Implements ErrorFormatterService.ErrorOnly.
func (s *serviceScoped) ErrorOnly(err error, details ...string) error {
	return ScopedErrorOnly(err, s.scope, details...)
}

// Error wraps err with details, applying additional formatting specific
// to the concrete service strategy.
// Implements ErrorFormatterService.Error.
func (s *serviceScoped) Error(err error, details ...string) error {
	return ScopedError(err, s.scope, details...)
}

// Errorf formats err with a printf-style string and arguments.
// Implements ErrorFormatterService.Errorf.
func (s *serviceScoped) Errorf(err error, format string, args ...interface{}) error {
	return ScopedErrorf(err, s.scope, format, args...)
}

// NewError creates a new error (no wrapper) from detail strings.
// Implements ErrorFormatterService.NewError.
func (s *serviceScoped) NewError(details ...string) error {
	return NewScopedError(s.scope, details...)
}

// NewErrorf creates a new error (no wrapper) from a format string.
// Implements ErrorFormatterService.NewErrorf.
func (s *serviceScoped) NewErrorf(format string, args ...interface{}) error {
	return NewScopedErrorf(format, s.scope, args...)
}

// NewScopedErrorFormatter creates a new ErrorFormatterService that prefixes
// every formatted message with "[<scope>] ". The returned *serviceScoped
// implements the full ErrorFormatterService interface.
func NewScopedErrorFormatter(scope string) *serviceScoped {
	return &serviceScoped{
		scope: scope,
	}
}
