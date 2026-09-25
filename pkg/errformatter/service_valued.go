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

// serviceValued attaches Value-typed metadata to errors. Details, scope, code,
// and public-code information flows through the Valued* helper functions
// defined in the package.
type serviceValued struct{}

var _ ErrorFormatterService = (*serviceValued)(nil)

// ErrCodeIsOneOf is a short alias for ErrorCodeIsOneOf.
// Implements ErrorFormatterService.ErrCodeIsOneOf.
func (s *serviceValued) ErrCodeIsOneOf(err error, codes ...int) (int, bool) {
	return s.ErrorCodeIsOneOf(err, codes...)
}

// ErrorCodeIsOneOf reports whether err carries a code present in codes.
// Returns the matched code (or -1).
// Implements ErrorFormatterService.ErrorCodeIsOneOf.
func (s *serviceValued) ErrorCodeIsOneOf(err error, codes ...int) (int, bool) {
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

// NewErrorWithCode creates a new error carrying details and code, no wrapper.
// Panics if code <= 0.
// Implements ErrorFormatterService.NewErrorWithCode.
func (s *serviceValued) NewErrorWithCode(text string, code int) error {
	return ValuedNewError([]Value{
		NewValue(KindDetails, text),
		NewValue(KindCode, code),
	})
}

// ErrGetCode is a short alias for ErrorGetCode.
// Implements ErrorFormatterService.ErrGetCode.
func (s *serviceValued) ErrGetCode(err error) int {
	return s.ErrorGetCode(err)
}

// ErrorGetCode extracts the embedded error code from err. Returns -1
// when the error does not carry a code.
// Implements ErrorFormatterService.ErrorGetCode.
func (s *serviceValued) ErrorGetCode(err error) int {
	return ValuedErrorGetCode(err)
}

// ErrWithCode is a short alias for ErrorWithCode.
// Implements ErrorFormatterService.ErrWithCode.
func (s *serviceValued) ErrWithCode(err error, code int) error {
	return s.ErrorWithCode(err, code)
}

// ErrNoWrap is a short alias for ErrorNoWrap.
// Implements ErrorFormatterService.ErrNoWrap.
func (s *serviceValued) ErrNoWrap(err error) error {
	return s.ErrorNoWrap(err)
}

// ErrorNoWrap wraps err with formatting details but does NOT chain err
// as Inner(). Callers who need a clean output use this method.
// Implements ErrorFormatterService.ErrorNoWrap.
func (s *serviceValued) ErrorNoWrap(err error) error {
	return ErrorNoWrap(err)
}

// ErrorWithCode creates an error wrapping err that carries the given
// positive code. If code <= 0 the call panics.
// Implements ErrorFormatterService.ErrorWithCode.
func (s *serviceValued) ErrorWithCode(err error, code int) error {
	if code <= 0 {
		panic("errfmt: code must be positive value")
	}

	return ValuedErrorOnly(err, NewValue(KindCode, code))
}

// ErrorOnly wraps err with detail strings wrapped in a KindDetails Value.
// Implements ErrorFormatterService.ErrorOnly.
func (s *serviceValued) ErrorOnly(err error, details ...string) error {
	return ValuedErrorOnly(err, NewValue(KindDetails, details))
}

// Errorf formats err with a printf-style string and arguments.
// Implements ErrorFormatterService.Errorf.
func (s *serviceValued) Errorf(err error, format string, args ...interface{}) error {
	return ValuedErrorf(err, nil, format, args...)
}

// Error wraps err with detail strings wrapped in a KindDetails Value.
// Implements ErrorFormatterService.Error.
func (s *serviceValued) Error(err error, details ...string) error {
	return ValuedError(err, nil, details...)
}

// NewError creates a new error (no wrapper) from detail strings.
// Implements ErrorFormatterService.NewError.
func (s *serviceValued) NewError(details ...string) error {
	return ValuedNewError(nil, details...)
}

// NewErrorf creates a new error (no wrapper) from a format string.
// Implements ErrorFormatterService.NewErrorf.
func (s *serviceValued) NewErrorf(format string, args ...interface{}) error {
	return ValuedNewErrorf(nil, format, args...)
}

func newValuedErrFmtService(values ...Value) ErrorFormatterService {
	if len(values) > 0 {
		return &serviceValuedWithDefaults{
			serviceValued: &serviceValued{},
			defaultValues: values,
		}
	}

	return &serviceValued{}
}

// NewValuesErrorFormatter returns an ErrorFormatterService that attaches
// the provided values to every formatted error. If values is non-empty
// the result is a serviceValuedWithDefaults; otherwise a plain serviceValued.
func NewValuesErrorFormatter(values ...Value) ErrorFormatterService {
	return newValuedErrFmtService(values...)
}

// NewValuedErrorFormatter returns an ErrorFormatterService that carries all four
// value kinds simultaneously — KindScope, KindCode, KindDetails, and KindPublicCode —
// making it a universal error-formatting solution. When values is non-empty the result
// is a serviceValuedWithDefaults so every error inherits those defaults; otherwise it
// returns a plain serviceValued suitable for one-off enriched errors.
func NewValuedErrorFormatter(values ...Value) ErrorFormatterService {
	return newValuedErrFmtService(values...)
}
