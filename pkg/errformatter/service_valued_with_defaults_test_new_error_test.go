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
	"testing"
)

func TestServiceValuedWithDefaults_NewError(t *testing.T) {
	runName1 := "valued error - error with 3 scope values for overwrite and details value for overwrite"
	t.Run(runName1, func(t *testing.T) {
		var (
			expectedResult = "valued_err_scope: some error text -> " +
				"err_detail_true_1, err_detail_true_2"
		)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindScope,
				any: "wrong_err_scope_for_overwrite",
			},
			{
				num: KindScope,
				any: "wrong_err_scope_2_for_overwrite",
			},
			{
				num: KindScope,
				any: "valued_err_scope",
			},
			{
				num: KindDetails,
				any: []string{"err_detail_true_1", "err_detail_true_2"},
			},
		}...)

		err := svc.NewError("some error text")
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}
	})

	runName2 := "valued error - error with one scope value, one details value and one code value"
	t.Run(runName2, func(t *testing.T) {
		var (
			expectedResult = "valued_err_scope: detail_info_as_value_1, detail_info_as_value_2 -> " +
				"detail_in_value_1, detail_in_value_2"
			expectedCode = 9999
		)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindScope,
				any: "valued_err_scope",
			},
			{
				num: KindDetails,
				any: []string{"detail_in_value_1", "detail_in_value_2"},
			},
			{
				num: KindCode,
				any: 100503,
			},
			{
				num: KindCode,
				any: expectedCode,
			},
		}...)

		err := svc.NewError("detail_info_as_value_1", "detail_info_as_value_2")
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

func TestServiceValuedWithDefaults_NewError_ReWrap(t *testing.T) {
	runName1 := "valued newError - error with 3 scope values for overwrite and details value for overwrite plus " +
		"re-wrap by same fmt service"
	t.Run(runName1, func(t *testing.T) {
		var (
			expectedResult = "valued_err_scope: some error text -> " +
				"err_detail_true_1"

			expectedReWrapedResult = "valued_err_scope: some error text -> " +
				"err_detail_true_1, err_detail_true_2"
		)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindScope,
				any: "wrong_err_scope_for_overwrite",
			},
			{
				num: KindScope,
				any: "wrong_err_scope_2_for_overwrite",
			},
			{
				num: KindScope,
				any: "valued_err_scope",
			},
			{
				num: KindDetails,
				any: []string{"err_detail_true_1"},
			},
		}...)

		err := svc.NewError("some error text")
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}

		reWrapedErr := svc.ErrorOnly(err, "err_detail_true_2")
		if reWrapedErr.Error() != expectedReWrapedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				reWrapedErr.Error(), expectedReWrapedResult)
		}
	})

	runName2 := "valued error - error with one scope value, one details value and one code value with re-wrap by same service"
	t.Run(runName2, func(t *testing.T) {
		var (
			expectedResult = "valued_err_scope: detail_info_as_value_1, detail_info_as_value_2 -> " +
				"detail_in_value_1, detail_in_value_2"
			expectedCode            = 9999
			expectedReWrappedResult = "valued_err_scope: detail_info_as_value_1, detail_info_as_value_2 -> " +
				"detail_in_value_1, detail_in_value_2"
			reWrappedExpectedCode = 8765
		)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindScope,
				any: "valued_err_scope",
			},
			{
				num: KindDetails,
				any: []string{"detail_in_value_1", "detail_in_value_2"},
			},
			{
				num: KindCode,
				any: 100503,
			},
			{
				num: KindCode,
				any: expectedCode,
			},
		}...)

		err := svc.NewError("detail_info_as_value_1", "detail_info_as_value_2")
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

		reWrappedErr := svc.ErrorWithCode(err, reWrappedExpectedCode)
		if reWrappedErr.Error() != expectedReWrappedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				reWrappedErr.Error(), expectedReWrappedResult)
		}

		if code := svc.ErrorGetCode(err); code != reWrappedExpectedCode {
			t.Errorf("error code not equal with expected. current: %d, expected: %d",
				code, reWrappedExpectedCode)
		}
	})

	runName3 := "valued error - error with one scope value, one details value and one code value with re-wrap by other service"
	t.Run(runName3, func(t *testing.T) {
		var (
			expectedResult = "valued_err_scope: detail_info_as_value_1, detail_info_as_value_2 -> " +
				"detail_in_value_1"
			expectedCode            = 9999
			expectedReWrappedResult = "re_wrap_scope: valued_err_scope: detail_info_as_value_1, detail_info_as_value_2 -> " +
				"detail_in_value_1 -> re_wrap_detail_in_value_1"
		)

		svc := NewValuesErrorFormatter([]Value{
			{
				num: KindScope,
				any: "valued_err_scope",
			},
			{
				num: KindDetails,
				any: []string{"detail_in_value_1"},
			},
			{
				num: KindCode,
				any: 100503,
			},
			{
				num: KindCode,
				any: expectedCode,
			},
		}...)

		err := svc.NewError("detail_info_as_value_1", "detail_info_as_value_2")
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

		reWrappedErr := reWrapSvc.Error(err, "re_wrap_detail_in_value_1")
		if reWrappedErr.Error() != expectedReWrappedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				reWrappedErr.Error(), expectedReWrappedResult)
		}
	})
}
