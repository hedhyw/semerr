// Code generated by "github.com/hedhyw/semerr"; DO NOT EDIT.

package grpcerr

import (
	"errors"

	"github.com/hedhyw/semerr/internal/pkg/multierr"
	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Code returns grpc status code for err. In case of joined errors
// it returns the first code found in the chain.
func Code(err error) codes.Code {
	return codeRecursion(err)
}

func codeRecursion(err error) codes.Code {
	switch err.(type) {
	case nil:
		return codes.OK
	case semerr.StatusRequestTimeoutError:
		return 1
	case semerr.InternalServerError:
		return 2
	case semerr.BadRequestError:
		return 3
	case semerr.UnsupportedMediaTypeError:
		return 3
	case semerr.StatusGatewayTimeoutError:
		return 4
	case semerr.NotFoundError:
		return 5
	case semerr.ConflictError:
		return 6
	case semerr.ForbiddenError:
		return 7
	case semerr.TooManyRequestsError:
		return 8
	case semerr.RequestEntityTooLargeError:
		return 11
	case semerr.UnimplementedError:
		return 12
	case semerr.ServiceUnavailableError:
		return 14
	case semerr.UnauthorizedError:
		return 16
	}

	grpcCode := getGRPCErrorCode(err)
	if grpcCode != codes.Unknown {
		return grpcCode
	}

	if err := errors.Unwrap(err); err != nil {
		return codeRecursion(err)
	}

	if multiErr, ok := err.(multierr.MultiError); ok && multiErr != nil {
		for _, err := range multiErr.Unwrap() {
			code := codeRecursion(err)

			if code != codes.Unknown {
				return code
			}
		}
	}

	return codes.Unknown
}

func getGRPCErrorCode(err error) codes.Code {
	var errGRPC interface {
		GRPCStatus() *status.Status
		error
	}

	if errors.As(err, &errGRPC) {
		status := errGRPC.GRPCStatus()

		if status == nil {
			return codes.Unknown
		}

		code := status.Code()
		if code != codes.OK && code != codes.Unknown {
			return code
		}
	}

	return codes.Unknown
}

// Wrap wraps the `err` with an error corresponding to the `code`.
// If there is no `err` for this code then the `err` will be returned
// without wrapping.
func Wrap(err error, code codes.Code) error {
	switch code {
	case 1:
		return semerr.NewStatusRequestTimeoutError(err)
	case 2:
		return semerr.NewInternalServerError(err)
	case 3:
		return semerr.NewBadRequestError(err)
	case 4:
		return semerr.NewStatusGatewayTimeoutError(err)
	case 5:
		return semerr.NewNotFoundError(err)
	case 6:
		return semerr.NewConflictError(err)
	case 7:
		return semerr.NewForbiddenError(err)
	case 8:
		return semerr.NewTooManyRequestsError(err)
	case 11:
		return semerr.NewRequestEntityTooLargeError(err)
	case 12:
		return semerr.NewUnimplementedError(err)
	case 14:
		return semerr.NewServiceUnavailableError(err)
	case 16:
		return semerr.NewUnauthorizedError(err)
	default:
		return err
	}
}
