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

// formatterBuilder implements ErrorFormatterBuidler.
// It is a lightweight structural type — all logic lives in the
// factory functions that each MakeXxx method delegates to.
type formatterBuilder struct{}

var _ ErrorFormatterBuidler = (*formatterBuilder)(nil)

// MakeScoped returns a scope-prefixed formatter configured with the given
// scope. The scope is prepended as "[scope] " to all messages produced
// through the returned ErrorFormatterService.
// Implements ErrorFormatterBuidler.MakeScoped.
func (r *formatterBuilder) MakeScoped(scope string) ErrorFormatterService {
	return NewScopedErrorFormatter(scope)
}

// MakeValued returns a value-based formatter initialized with vals.
// If vals is empty the result is equivalent to a plain valued formatter
// (serviceValuedWithoutDefaults).
// Implements ErrorFormatterBuidler.MakeValued.
func (r *formatterBuilder) MakeValued(values ...Value) ErrorFormatterService {
	return NewValuesErrorFormatter(values...)
}

// MakeSimply returns the basic formatter that delegates to standard library
// calls, producing plain formatted messages without enrichment.
// Implements ErrorFormatterBuidler.MakeSimply.
func (r *formatterBuilder) MakeSimply() ErrorFormatterService {
	return NewErrorBasicFormatter()
}

// Make returns the default formatter — equivalent to MakeSimply.
// Implements ErrorFormatterBuidler.Make.
func (r *formatterBuilder) Make() ErrorFormatterService {
	return NewErrorFormatter()
}
