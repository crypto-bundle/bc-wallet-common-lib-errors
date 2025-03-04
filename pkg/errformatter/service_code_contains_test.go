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

func TestServiceCodeContains_ErrorWithCode(t *testing.T) {
	t.Run("default formatter ErrorWithCode - error with expected code as test value", func(t *testing.T) {
		const (
			expectedResult      = "test error"
			expectedTextForWrap = "test error"
			expectedCode        = 456789
		)
		var errorForWrap = errors.New(expectedTextForWrap)

		svc := NewErrorFormatter()

		err := svc.ErrorWithCode(errorForWrap, expectedCode)
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
	})

	t.Run("default formatter ErrorWithCode - re-wrap of error", func(t *testing.T) {
		var (
			expectedResult      = "test error"
			expectedTextForWrap = "test error"
			unExpectedCode      = 1234
			expectedCode        = 5550
		)
		var errorForWrap = errors.New(expectedTextForWrap)

		svc := NewErrorFormatter()

		err := svc.ErrorWithCode(errorForWrap, unExpectedCode)
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

		reWrappedErr := svc.ErrorWithCode(unwrappedErr, expectedCode)

		if code := svc.ErrorGetCode(reWrappedErr); code != expectedCode {
			t.Errorf("error code not equal with expected. current: %d, expected: %d",
				code, expectedCode)
		}
	})
}
