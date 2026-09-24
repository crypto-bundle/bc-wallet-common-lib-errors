package errformatter

import (
	"errors"
	"testing"
)

func TestErrorFormatterBuilder(t *testing.T) {

	t.Run("error only - scoped error formatter", func(t *testing.T) {
		blrd := formatterBuilder{}
		formatterSvc := blrd.MakeScoped("test_scope")

		const expectedResult = "test_scope: test error"

		err := formatterSvc.ErrorOnly(errors.New("test error"))
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}
	})

	t.Run("error only - simply error formatter", func(t *testing.T) {
		blrd := formatterBuilder{}
		formatterSvc := blrd.MakeSimply()

		const expectedResult = "test error -> efg"

		err := formatterSvc.Error(errors.New("test error"), "efg")
		if err.Error() != expectedResult {
			t.Errorf("error text not equal with expected. current: %s, expected: %s",
				err.Error(), expectedResult)
		}
	})

	t.Run("error - valued error formatter", func(t *testing.T) {
		blrd := formatterBuilder{}
		formatterSvc := blrd.MakeValued([]Value{
			{
				num: KindScope,
				any: "wrong_err_scope",
			},
			{
				num: KindScope,
				any: "valued_err_scope",
			},
			{
				num: KindCode,
				any: 199991,
			},
		}...)

		const (
			expectedResult      = "valued_err_scope: test error"
			expectedTextForWrap = "test error"
			expectedCode        = 456789
		)
		var errorForWrap = errors.New(expectedTextForWrap)

		err := formatterSvc.ErrorWithCode(errorForWrap, expectedCode)
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

		if code := formatterSvc.ErrorGetCode(err); code != expectedCode {
			t.Errorf("error code not equal with expected. current: %d, expected: %d",
				code, expectedCode)
		}

		if code := ValuedErrorGetCode(err); code != expectedCode {
			t.Errorf("error code not equal with expected. current: %d, expected: %d",
				code, expectedCode)
		}
	})
}
