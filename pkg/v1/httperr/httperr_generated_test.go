// Code generated by "github.com/hedhyw/semerr"; DO NOT EDIT.

package httperr_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/hedhyw/semerr/pkg/v1/httperr"
	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

func TestCode(t *testing.T) {
	t.Parallel()

	const err = semerr.Error("some error")

	testCases := []struct {
		Err  error
		Code int
	}{
		{
			Err:  nil,
			Code: http.StatusOK,
		},
		{
			Err:  err,
			Code: http.StatusInternalServerError,
		},
		{
			Err:  semerr.NewStatusRequestTimeoutError(err),
			Code: 408,
		},
		{
			Err:  semerr.NewInternalServerError(err),
			Code: 500,
		},
		{
			Err:  semerr.NewBadRequestError(err),
			Code: 400,
		},
		{
			Err:  semerr.NewUnsupportedMediaTypeError(err),
			Code: 415,
		},
		{
			Err:  semerr.NewStatusGatewayTimeoutError(err),
			Code: 504,
		},
		{
			Err:  semerr.NewNotFoundError(err),
			Code: 404,
		},
		{
			Err:  semerr.NewConflictError(err),
			Code: 409,
		},
		{
			Err:  semerr.NewForbiddenError(err),
			Code: 403,
		},
		{
			Err:  semerr.NewTooManyRequestsError(err),
			Code: 429,
		},
		{
			Err:  semerr.NewRequestEntityTooLargeError(err),
			Code: 413,
		},
		{
			Err:  semerr.NewUnimplementedError(err),
			Code: 501,
		},
		{
			Err:  semerr.NewServiceUnavailableError(err),
			Code: 503,
		},
		{
			Err:  semerr.NewUnauthorizedError(err),
			Code: 401,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprint(tc.Err), func(t *testing.T) {
			t.Parallel()

			err := tc.Err
			gotCode := httperr.Code(err)
			if tc.Code != gotCode {
				t.Fatal("exp", tc.Code, "got", gotCode)
			}

			if err != nil {
				err = fmt.Errorf("wrapped: 1: %w", err)
				err = fmt.Errorf("wrapped: 2: %w", err)
				gotCode = httperr.Code(err)
				if tc.Code != gotCode {
					t.Fatal("exp", tc.Code, "got", gotCode)
				}
			}
		})
	}
}

func TestWrap(t *testing.T) {
	t.Parallel()

	const err = semerr.Error("some error")

	testCases := []struct {
		Code  int
		Check func(err error) bool
	}{
		{
			Check: func(actualErr error) bool {
				return err == actualErr
			},
			Code: -1,
		},
		{
			Check: func(err error) bool {
				return errors.As(err, &semerr.StatusRequestTimeoutError{})
			},
			Code: 408,
		},
		{
			Check: func(err error) bool {
				return errors.As(err, &semerr.InternalServerError{})
			},
			Code: 500,
		},
		{
			Check: func(err error) bool {
				return errors.As(err, &semerr.BadRequestError{})
			},
			Code: 400,
		},
		{
			Check: func(err error) bool {
				return errors.As(err, &semerr.StatusGatewayTimeoutError{})
			},
			Code: 504,
		},
		{
			Check: func(err error) bool {
				return errors.As(err, &semerr.NotFoundError{})
			},
			Code: 404,
		},
		{
			Check: func(err error) bool {
				return errors.As(err, &semerr.ConflictError{})
			},
			Code: 409,
		},
		{
			Check: func(err error) bool {
				return errors.As(err, &semerr.ForbiddenError{})
			},
			Code: 403,
		},
		{
			Check: func(err error) bool {
				return errors.As(err, &semerr.TooManyRequestsError{})
			},
			Code: 429,
		},
		{
			Check: func(err error) bool {
				return errors.As(err, &semerr.RequestEntityTooLargeError{})
			},
			Code: 413,
		},
		{
			Check: func(err error) bool {
				return errors.As(err, &semerr.UnimplementedError{})
			},
			Code: 501,
		},
		{
			Check: func(err error) bool {
				return errors.As(err, &semerr.ServiceUnavailableError{})
			},
			Code: 503,
		},
		{
			Check: func(err error) bool {
				return errors.As(err, &semerr.UnauthorizedError{})
			},
			Code: 401,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprint(tc.Code), func(t *testing.T) {
			t.Parallel()

			if err := httperr.Wrap(err, tc.Code); !tc.Check(err) {
				t.Fatalf("%T", err)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	t.Parallel()

	const err = semerr.Error("some error")

	gotCode := httperr.Code(errors.Join(
		fmt.Errorf("regular: %w", err),
		fmt.Errorf("bad request: %w", semerr.NewBadRequestError(err)),
		semerr.NewNotFoundError(fmt.Errorf("not found: %w", err)),
	))

	if gotCode != http.StatusBadRequest {
		t.Fatal("exp", http.StatusBadRequest, "got", gotCode)
	}
}
