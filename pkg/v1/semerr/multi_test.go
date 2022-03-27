package semerr_test

import (
	"errors"
	"testing"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

func TestNewMultiError(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		err := semerr.NewMultiError()
		if err != nil {
			t.Fatal(err)
		}

		err = &semerr.MultiErr{}
		if errors.Unwrap(err) != nil {
			t.Fatal(errors.Unwrap(err))
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		err := semerr.NewMultiError(nil)
		if err != nil {
			t.Fatal(err)
		}

		err = semerr.NewMultiError(nil, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("err and nil", func(t *testing.T) {
		t.Parallel()

		err := semerr.Error(t.Name())

		merr := semerr.NewMultiError(err, nil)
		if merr != err {
			t.Fatal(err)
		}

		merr = semerr.NewMultiError(nil, err)
		if merr != err {
			t.Fatal(err)
		}
	})

	t.Run("two errs", func(t *testing.T) {
		t.Parallel()

		errFrist := semerr.Error("first error")
		errSecond := semerr.Error("second error")

		err := semerr.NewMultiError(errFrist, errSecond)
		switch {
		case !errors.Is(err, errFrist):
			t.Fatal(err)
		case errors.Is(err, errSecond):
			t.Fatal(err)
		case err.Error() != errFrist.Error()+"; "+errSecond.Error():
			t.Fatal(err)
		}
	})
}
