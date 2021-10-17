package grpcerr_test

import (
	"errors"
	"fmt"

	"github.com/hedhyw/semerr/pkg/v1/grpcerr"
	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"google.golang.org/grpc/codes"
)

func ExampleCode() {
	err := semerr.NewBadRequestError(errors.New("bad request error"))

	fmt.Println(
		"Is err invalid argument?",
		grpcerr.Code(err) == codes.InvalidArgument,
	)
	// Output: Is err invalid argument? true
}
