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

func TestServiceValuedWithDefaults_Errorf(t *testing.T) {
	t.Run("service valued errorf - formatted error with kindCode value", func(t *testing.T) {
		const (
			expectedCode        = 404
			expectedResult      = "test error -> 100501 is error value"
			expectedTextForWrap = "test error"
		)
		var errorForWrap = errors.New(expectedTextForWrap)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindCode,
				any: expectedCode,
			},
		}...)

		err := svc.Errorf(errorForWrap, "%d %s", 100501, "is error value")
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}

		unwrappedErr := errors.Unwrap(err)

		if !errors.Is(unwrappedErr, errorForWrap) {
			t.Errorf("error text not equal with expected. current: %e, expected: %e",
				unwrappedErr, errorForWrap)
		}

		if unwrappedErr.Error() != expectedTextForWrap {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				unwrappedErr.Error(), expectedTextForWrap)
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

	t.Run("service valued errorf - formatted error with 2 kindCode values for overwrite", func(t *testing.T) {
		const (
			expectedCode        = 404
			expectedResult      = "test error -> 100502 is error value"
			expectedTextForWrap = "test error"
		)
		var errorForWrap = errors.New(expectedTextForWrap)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindCode,
				any: 100555,
			},
			{
				num: KindCode,
				any: expectedCode,
			},
		}...)

		err := svc.Errorf(errorForWrap, "%d %s", 100502, "is error value")
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}

		unwrappedErr := errors.Unwrap(err)

		if !errors.Is(unwrappedErr, errorForWrap) {
			t.Errorf("error text not equal with expected. current: %e, expected: %e",
				unwrappedErr, errorForWrap)
		}

		if unwrappedErr.Error() != expectedTextForWrap {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				unwrappedErr.Error(), expectedTextForWrap)
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

func TestServiceValuedWithDefaults_Errorf_ReWrap(t *testing.T) {
	t.Run("service valued errorf - formatted error with kindCode value and re-wrap by same fmt service",
		func(t *testing.T) {
			const (
				expectedCode         = 404
				expectedResult       = "test error -> 100501 is error value"
				expectedTextForWrap  = "test error"
				expectedReWrapResult = "test error -> 100501 is error value -> some_details_for_re_wrap_1"
			)
			var errorForWrap = errors.New(expectedTextForWrap)

			svc := NewValuesErrorFormatter([]Value{
				{
					num: KindCode,
					any: expectedCode,
				},
			}...)

			err := svc.Errorf(errorForWrap, "%d %s", 100501, "is error value")
			if err.Error() != expectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					err.Error(), expectedResult)
			}

			unwrappedErr := errors.Unwrap(err)

			if !errors.Is(unwrappedErr, errorForWrap) {
				t.Errorf("error text not equal with expected. current: %e, expected: %e",
					unwrappedErr, errorForWrap)
			}

			if unwrappedErr.Error() != expectedTextForWrap {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					unwrappedErr.Error(), expectedTextForWrap)
			}

			if code := svc.ErrorGetCode(err); code != expectedCode {
				t.Errorf("error code not equal with expected. current: %d, expected: %d",
					code, expectedCode)
			}

			if code := ValuedErrorGetCode(err); code != expectedCode {
				t.Errorf("error code not equal with expected. current: %d, expected: %d",
					code, expectedCode)
			}

			reWrappedErr := svc.Error(err, "some_details_for_re_wrap_1")
			if reWrappedErr.Error() != expectedReWrapResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					reWrappedErr.Error(), expectedReWrapResult)
			}
		})

	t.Run("service valued errorf - formatted error with 2 kindCode values for overwrite and re-wrap by another service",
		func(t *testing.T) {
			const (
				expectedCode            = 404
				expectedResult          = "test error -> 100502 is error value"
				reWrappedExpectedResult = "re_wrap_scope: test error -> 100502 is error value -> re_wrap_details_1"
				expectedTextForWrap     = "test error"
			)
			var errorForWrap = errors.New(expectedTextForWrap)

			svc := NewValuesErrorFormatter([]Value{
				{
					num: KindCode,
					any: 100555,
				},
				{
					num: KindCode,
					any: expectedCode,
				},
			}...)

			err := svc.Errorf(errorForWrap, "%d %s", 100502, "is error value")
			if err.Error() != expectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					err.Error(), expectedResult)
			}

			unwrappedErr := errors.Unwrap(err)

			if !errors.Is(unwrappedErr, errorForWrap) {
				t.Errorf("error text not equal with expected. current: %e, expected: %e",
					unwrappedErr, errorForWrap)
			}

			if unwrappedErr.Error() != expectedTextForWrap {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					unwrappedErr.Error(), expectedTextForWrap)
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

			reWrappedErr := reWrapSvc.Error(err, "re_wrap_details_1")
			if reWrappedErr.Error() != reWrappedExpectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					reWrappedErr.Error(), reWrappedExpectedResult)
			}

			unwrappedReWrappedErr := errors.Unwrap(err)
			if unwrappedReWrappedErr.Error() != expectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					unwrappedReWrappedErr.Error(), expectedResult)
			}
		})

	t.Run("service valued errorf - formatted error with 2 kindCode values for overwrite and re-wrap by another service with errorf call",
		func(t *testing.T) {
			var (
				expectedCode              = 404
				expectedResult            = "test error -> 100502 is error value"
				reWrappedExpectedResult   = "test error -> 100502 is error value -> error with formatted_text_13_value"
				reReWrappedExpectedResult = "re_wrap_scope: test error -> 100502 is error value -> " +
					"some_re_re_wrap_detail_1"
				expectedTextForWrap = "test error"
			)

			var errorForWrap = errors.New(expectedTextForWrap)

			svc := NewValuesErrorFormatter([]Value{
				{
					num: KindCode,
					any: 100555,
				},
				{
					num: KindCode,
					any: expectedCode,
				},
			}...)

			err := svc.Errorf(errorForWrap, "%d %s", 100502, "is error value")
			if err.Error() != expectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					err.Error(), expectedResult)
			}

			unwrappedErr := errors.Unwrap(err)

			if !errors.Is(unwrappedErr, errorForWrap) {
				t.Errorf("error text not equal with expected. current: %e, expected: %e",
					unwrappedErr, errorForWrap)
			}

			if unwrappedErr.Error() != expectedTextForWrap {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					unwrappedErr.Error(), expectedTextForWrap)
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

			reWrappedErr := reWrapSvc.Errorf(err, "error with formatted_text_%d_%s", 13, "value")
			if reWrappedErr.Error() != reWrappedExpectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					reWrappedErr.Error(), reWrappedExpectedResult)
			}

			unwrappedReWrappedErr := errors.Unwrap(err)
			if unwrappedReWrappedErr.Error() != expectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					unwrappedReWrappedErr.Error(), expectedResult)
			}

			reReWrappedErr := reWrapSvc.Error(reWrappedErr, "some_re_re_wrap_detail_1")
			if reReWrappedErr.Error() != reReWrappedExpectedResult {
				t.Errorf("error text not equal with expected. current: %s, expected: %s",
					reReWrappedErr.Error(), reReWrappedExpectedResult)
			}
		})
}
