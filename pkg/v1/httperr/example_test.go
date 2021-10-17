package httperr_test

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hedhyw/semerr/pkg/v1/httperr"
	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

func ExampleCode() {
	err := semerr.NewBadRequestError(errors.New("bad request error"))

	fmt.Println(
		"Is err bad request?",
		httperr.Code(err) == http.StatusBadRequest,
	)
	// Output: Is err bad request? true
}
