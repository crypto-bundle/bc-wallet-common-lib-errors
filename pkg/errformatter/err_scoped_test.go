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

func TestScopedErrorFormatting(t *testing.T) {
	t.Run("error only - scoped error with details", func(t *testing.T) {
		const expectedResult = "test_scope: test error -> abc"

		err := ScopedErrorOnly(errors.New("test error"), "test_scope", "abc")
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}
	})

	t.Run("scoped error - details and func", func(t *testing.T) {
		const expectedResult = "test_scope: test error -> abcd, efg, [TestScopedErrorFormatting.func2]"

		err := ScopedError(errors.New("test error"), "test_scope", "abcd", "efg")
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}
	})

	t.Run("new scoped error - 2 details args and func", func(t *testing.T) {
		const expectedResult = "test_scope: error detail, error detail2, [TestScopedErrorFormatting.func3]"

		err := NewScopedError("test_scope", "error detail", "error detail2")
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}
	})

	t.Run("new scoped formatted error - 2 fmt args and func", func(t *testing.T) {
		const expectedResult = "test_scope: test fmt_arg1 fmt_arg2, [TestScopedErrorFormatting.func4]"

		err := NewScopedErrorf("test %s %s", "test_scope", "fmt_arg1", "fmt_arg2")
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}
	})

	t.Run("scoped formatted error - two fmt args and func name", func(t *testing.T) {
		const expectedResult = "test_scope: test error -> test fmt_arg1 fmt_arg2, [TestScopedErrorFormatting.func5]"

		err := ScopedErrorf(errors.New("test error"), "test_scope", "test %s %s",
			"fmt_arg1", "fmt_arg2")
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}
	})
}
