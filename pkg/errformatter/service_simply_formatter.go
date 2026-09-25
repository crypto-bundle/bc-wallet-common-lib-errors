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

var _ ErrorFormatterService = (*service)(nil)

// service is the basic error formatter. It delegates all formatting
// to the package-level helper functions (Error, Errorf, etc.), so
// the produced messages contain no enrichment beyond the wrapped error text.
type service struct{}

// NewErrorBasicFormatter returns a basic error formatter that delegates to
// standard library calls. It produces plain formatted messages without any
// enrichment such as scope prefixes or value-based metadata.
// The returned service implements ErrorFormatterService.
func NewErrorBasicFormatter() *service {
	return &service{}
}

// ErrCodeIsOneOf is a short alias for ErrorCodeIsOneOf.
// Implements ErrorFormatterService.ErrCodeIsOneOf.
func (s *service) ErrCodeIsOneOf(err error, codes ...int) (int, bool) {
	return s.ErrorCodeIsOneOf(err, codes...)
}

// ErrorCodeIsOneOf reports whether err carries a code present in codes.
// Returns the matched code (or -1).
// Implements ErrorFormatterService.ErrorCodeIsOneOf.
func (s *service) ErrorCodeIsOneOf(err error, codes ...int) (int, bool) {
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

// NewErrorWithCode creates a new error carrying the given details and code,
// no wrapper. Panics if code <= 0.
// Implements ErrorFormatterService.NewErrorWithCode.
func (s *service) NewErrorWithCode(text string, code int) error {
	return ValuedNewError([]Value{
		NewValue(KindDetails, text),
		NewValue(KindCode, code),
	})
}

// ErrGetCode is a short alias for ErrorGetCode.
// Implements ErrorFormatterService.ErrGetCode.
func (s *service) ErrGetCode(err error) int {
	return s.ErrorGetCode(err)
}

// ErrorGetCode extracts the embedded error code from err. Returns -1
// when the error does not carry a code.
// Implements ErrorFormatterService.ErrorGetCode.
func (s *service) ErrorGetCode(err error) int {
	return ValuedErrorGetCode(err)
}

// ErrWithCode is a short alias for ErrorWithCode.
// Implements ErrorFormatterService.ErrWithCode.
func (s *service) ErrWithCode(err error, code int) error {
	return s.ErrorWithCode(err, code)
}

// ErrorWithCode creates an error wrapping err that carries the given
// positive code. If code <= 0 the call panics.
// Implements ErrorFormatterService.ErrorWithCode.
func (s *service) ErrorWithCode(err error, code int) error {
	if code <= 0 {
		panic("errfmt: code must be positive value")
	}

	return ValuedErrorOnly(err, NewValue(KindCode, code))
}

// ErrNoWrap is a short alias for ErrorNoWrap.
// Implements ErrorFormatterService.ErrNoWrap.
func (s *service) ErrNoWrap(err error) error {
	return s.ErrorNoWrap(err)
}

// ErrorNoWrap wraps err with formatting details but does NOT chain err
// as Inner(). Callers who need a clean output use this method.
// Implements ErrorFormatterService.ErrorNoWrap.
func (s *service) ErrorNoWrap(err error) error {
	return ErrorNoWrap(err)
}

// ErrorOnly wraps err with the provided detail strings.
// Implements ErrorFormatterService.ErrorOnly.
func (s *service) ErrorOnly(err error, details ...string) error {
	return ErrorOnly(err, details...)
}

// Error wraps err with details, applying additional formatting specific
// to the concrete service strategy.
// Implements ErrorFormatterService.Error.
func (s *service) Error(err error, details ...string) error {
	return Error(err, details...)
}

// Errorf formats err with a printf-style string and arguments.
// Implements ErrorFormatterService.Errorf.
func (s *service) Errorf(err error, format string, args ...interface{}) error {
	return Errorf(err, format, args...)
}

// NewError creates a new error (no wrapper) from detail strings.
// Implements ErrorFormatterService.NewError.
func (s *service) NewError(details ...string) error {
	return NewError(details...)
}

// NewErrorf creates a new error (no wrapper) from a format string.
// Implements ErrorFormatterService.NewErrorf.
func (s *service) NewErrorf(format string, args ...interface{}) error {
	return NewErrorf(format, args...)
}
