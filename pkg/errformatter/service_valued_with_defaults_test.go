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
	"testing"
)

func TestServiceValuedWithDefaults(t *testing.T) {
	t.Run("service valued error only - error with scope value and two details", func(t *testing.T) {
		const (
			expectedResult = "valued_err_scope: test error -> detail_1, detail_2"
			expectedCode   = 404
		)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindScope,
				any: "valued_err_scope",
			},
			{
				num: KindCode,
				any: expectedCode,
			},
		}...)

		err := svc.ErrorOnly(errors.New("test error"), "detail_1", "detail_2")
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}

		if code := svc.ErrorGetCode(err); code != expectedCode {
			t.Errorf("error code not equal with expected. current: %d, expected: %d",
				code, expectedCode)
		}

		if code := ValuedErrorGetCode(err); code != expectedCode {
			t.Errorf("error code not equal with expected. current: %d, expected: %d",
				code, expectedCode)
		}
	})

	t.Run("service valued error with code - error with scope value and code value rewrite", func(t *testing.T) {
		const (
			expectedResult = "valued_err_scope: test error"
			expectedCode   = 404
		)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindScope,
				any: "valued_err_scope",
			},
			{
				num: KindCode,
				any: 100500,
			},
		}...)

		err := svc.ErrorWithCode(errors.New("test error"), expectedCode)
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}

		if code := svc.ErrorGetCode(err); code != expectedCode {
			t.Errorf("error code not equal with expected. current: %d, expected: %d",
				code, expectedCode)
		}

		if code := ValuedErrorGetCode(err); code != expectedCode {
			t.Errorf("error code not equal with expected. current: %d, expected: %d",
				code, expectedCode)
		}
	})

	t.Run("valued error with code - error with scope value, details value and code value rewrite", func(t *testing.T) {
		const (
			expectedResult = "valued_err_scope: test error -> detail_1, detail_2"
			expectedCode   = 404
		)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindScope,
				any: "valued_err_scope",
			},
			{
				num: KindCode,
				any: 100500,
			},
			{
				num: KindDetails,
				any: []string{"detail_1", "detail_2"},
			},
		}...)

		err := svc.ErrorWithCode(errors.New("test error"), expectedCode)
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}

		if code := svc.ErrorGetCode(err); code != expectedCode {
			t.Errorf("error code not equal with expected. current: %d, expected: %d",
				code, expectedCode)
		}

		if code := ValuedErrorGetCode(err); code != expectedCode {
			t.Errorf("error code not equal with expected. current: %d, expected: %d",
				code, expectedCode)
		}
	})
}
