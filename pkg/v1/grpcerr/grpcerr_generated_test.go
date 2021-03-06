// Code generated by "github.com/hedhyw/semerr"; DO NOT EDIT.

package grpcerr_test

import (
	"fmt"
	"testing"

	"github.com/hedhyw/semerr/pkg/v1/grpcerr"
	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCode(t *testing.T) {
	t.Parallel()

	const err = semerr.Error("some error")

	testCases := []struct {
		Err  error
		Code codes.Code
	}{
		{
			Err:  nil,
			Code: codes.OK,
		},
		{
			Err:  err,
			Code: codes.Unknown,
		},
		{
			Err:  status.Error(codes.AlreadyExists, "already found"),
			Code: codes.AlreadyExists,
		},
		{
			Err:  semerr.NewStatusRequestTimeoutError(err),
			Code: 1,
		},
		{
			Err:  semerr.NewInternalServerError(err),
			Code: 2,
		},
		{
			Err:  semerr.NewBadRequestError(err),
			Code: 3,
		},
		{
			Err:  semerr.NewUnsupportedMediaTypeError(err),
			Code: 3,
		},
		{
			Err:  semerr.NewStatusGatewayTimeoutError(err),
			Code: 4,
		},
		{
			Err:  semerr.NewNotFoundError(err),
			Code: 5,
		},
		{
			Err:  semerr.NewConflictError(err),
			Code: 6,
		},
		{
			Err:  semerr.NewForbiddenError(err),
			Code: 7,
		},
		{
			Err:  semerr.NewTooManyRequestsError(err),
			Code: 8,
		},
		{
			Err:  semerr.NewRequestEntityTooLargeError(err),
			Code: 11,
		},
		{
			Err:  semerr.NewUnimplementedError(err),
			Code: 12,
		},
		{
			Err:  semerr.NewServiceUnavailableError(err),
			Code: 14,
		},
		{
			Err:  semerr.NewUnauthorizedError(err),
			Code: 16,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprint(tc.Err), func(t *testing.T) {
			t.Parallel()

			err := tc.Err
			gotCode := grpcerr.Code(err)
			if tc.Code != gotCode {
				t.Fatal("exp", tc.Code, "got", gotCode)
			}

			if err != nil {
				err = fmt.Errorf("wrapped: 1: %w", err)
				err = fmt.Errorf("wrapped: 2: %w", err)
				gotCode = grpcerr.Code(err)
				if tc.Code != gotCode {
					t.Fatal("exp", tc.Code, "got", gotCode)
				}
			}
		})
	}
}
