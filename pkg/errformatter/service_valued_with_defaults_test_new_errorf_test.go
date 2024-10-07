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

func TestServiceValuedWithDefaults_NewErrorf(t *testing.T) {
	t.Run("valued NewErrorf - formatted error with kindCode value", func(t *testing.T) {
		const (
			expectedCode   = 404
			expectedResult = "100501 is error value"
		)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindCode,
				any: expectedCode,
			},
		}...)

		err := svc.NewErrorf("%d %s", 100501, "is error value")
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

	t.Run("valued NewErrorf - formatted error with 2 kindCode values and details", func(t *testing.T) {
		const (
			expectedCode   = 404
			expectedResult = "100502 is error value -> detail_in_value_1, detail_in_value_2"
		)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindCode,
				any: 100555,
			},
			{
				num: KindDetails,
				any: []string{"detail_in_value_1", "detail_in_value_2"},
			},
			{
				num: KindCode,
				any: expectedCode,
			},
		}...)

		err := svc.NewErrorf("%d %s", 100502, "is error value")
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

	t.Run("valued NewErrorf - formatted error with 2 kindCode values, details and scope", func(t *testing.T) {
		const (
			expectedCode   = 405
			expectedResult = "err_scope: 100502 is error value -> detail_in_value_1, detail_in_value_2"
		)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindCode,
				any: 100556,
			},
			{
				num: KindDetails,
				any: []string{"detail_in_value_1", "detail_in_value_2"},
			},
			{
				num: KindScope,
				any: "err_scope",
			},
			{
				num: KindCode,
				any: expectedCode,
			},
		}...)

		err := svc.NewErrorf("%d %s", 100502, "is error value")
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

func TestServiceValuedWithDefaults_NewErrorf_ReWrap(t *testing.T) {
	t.Run("valued NewErrorf - formatted error with kindCode value and re-wrap by same formatter svc",
		func(t *testing.T) {
			const (
				expectedCode         = 404
				expectedResult       = "100501 is error value"
				reWrapExpectedResult = "100501 is error value -> re_wrap_detail_1"
			)

			svc := NewValuesErrorFormatter([]Value{
				{
					num: KindCode,
					any: expectedCode,
				},
			}...)

			err := svc.NewErrorf("%d %s", 100501, "is error value")
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

			reWrapErr := svc.ErrorOnly(err, "re_wrap_detail_1")
			if err.Error() != reWrapExpectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					err.Error(), reWrapExpectedResult)
			}

			unWrappedReWrapErr := errors.Unwrap(reWrapErr)
			if unWrappedReWrapErr.Error() != expectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					unWrappedReWrapErr.Error(), expectedResult)
			}
		})

	t.Run("valued NewErrorf - formatted error with 2 kindCode values and details and re-wrap by another formatter srvice",
		func(t *testing.T) {
			var (
				expectedCode         = 404
				expectedResult       = "100502 is error value -> detail_in_value_1, detail_in_value_2"
				reWrapExpectedResult = "re_wrap_scope: 100502 is error value -> " +
					"detail_in_value_1, detail_in_value_2 -> re_wrap_detail_in_value_1"
			)

			svc := NewValuesErrorFormatter([]Value{
				{
					num: KindCode,
					any: 100555,
				},
				{
					num: KindDetails,
					any: []string{"detail_in_value_1", "detail_in_value_2"},
				},
				{
					num: KindCode,
					any: expectedCode,
				},
			}...)

			err := svc.NewErrorf("%d %s", 100502, "is error value")
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

			reWrapSvc := NewValuesErrorFormatter([]Value{
				{
					num: KindScope,
					any: "re_wrap_scope",
				},
			}...)

			reWrapErr := reWrapSvc.ErrorOnly(err, "re_wrap_detail_in_value_1")
			if reWrapErr.Error() != reWrapExpectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					reWrapErr.Error(), reWrapExpectedResult)
			}

			unWrappedReWrapErr := errors.Unwrap(reWrapErr)
			if unWrappedReWrapErr.Error() != expectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					unWrappedReWrapErr.Error(), expectedResult)
			}
		})

	t.Run("valued NewErrorf - formatted error with 2 kindCode values, details and scope and re-wrap by another formatter service",
		func(t *testing.T) {
			var (
				expectedCode         = 405
				expectedResult       = "err_scope: 100502 is error value -> detail_in_value_1"
				reWrapExpectedResult = "re_wrap_scope: err_scope: 100502 is error value -> " +
					"detail_in_value_1 -> re_wrap_detail_in_value_1"
			)

			svc := NewValuesErrorFormatter([]Value{
				{
					num: KindCode,
					any: 100556,
				},
				{
					num: KindDetails,
					any: []string{"detail_in_value_1"},
				},
				{
					num: KindScope,
					any: "err_scope",
				},
				{
					num: KindCode,
					any: expectedCode,
				},
			}...)

			err := svc.NewErrorf("%d %s", 100502, "is error value")
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

			reWrapSvc := NewValuesErrorFormatter([]Value{
				{
					num: KindScope,
					any: "re_wrap_scope",
				},
			}...)

			reWrapErr := reWrapSvc.ErrorOnly(err, "re_wrap_detail_in_value_1")
			if reWrapErr.Error() != reWrapExpectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					reWrapErr.Error(), reWrapExpectedResult)
			}

			unWrappedReWrapErr := errors.Unwrap(reWrapErr)
			if unWrappedReWrapErr.Error() != expectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					unWrappedReWrapErr.Error(), expectedResult)
			}
		})
}
