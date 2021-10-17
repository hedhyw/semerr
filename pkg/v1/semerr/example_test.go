package semerr_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

// This example demonstrates returning of multiple errors.
func ExampleNewMultiError() {
	const (
		errFirst  semerr.Error = "first error"
		errSecond semerr.Error = "second error"
	)

	// All nil errors will be skipped.
	errMulti := semerr.NewMultiError(errFirst, errSecond, nil)

	// The first error is the main.
	fmt.Printf("errMulti is %q\n", errMulti.Error())
	fmt.Println("Is first? ", errors.Is(errMulti, errFirst))
	fmt.Println("Is second?", errors.Is(errMulti, errSecond))

	// Output:
	// errMulti is "first error; second error"
	// Is first?  true
	// Is second? false
}

// This example demonstrates using of const error.
func ExampleError() {
	const errImmutable semerr.Error = "immutable error"

	// semerr.Error implements error interface.
	var _ error = errImmutable

	fmt.Println(errImmutable.Error())
	// Output: immutable error
}

// This example demonstrates with meaningfull errors.
func ExampleBadRequestError() {
	const errExample semerr.Error = "example error"
	errBadReq := semerr.NewBadRequestError(errExample)

	fmt.Println(
		"Is the text the same?",
		semerr.NewBadRequestError(errExample).Error() == errExample.Error(),
	)

	fmt.Println(
		"Is the err original?",
		errors.Is(errBadReq, errExample),
	)

	fmt.Println(
		"Is nil handled?",
		semerr.NewBadRequestError(nil) == nil,
	)

	// Output:
	// Is the text the same? true
	// Is the err original? true
	// Is nil handled? true
}

func ExampleIsTemporaryError() {
	fmt.Println(
		"Is context.DeadlineExceeded tmp?",
		semerr.IsTemporaryError(context.DeadlineExceeded),
	)

	fmt.Println(
		"Is context.Canceled tmp?",
		semerr.IsTemporaryError(context.Canceled),
	)

	// Output:
	// Is context.DeadlineExceeded tmp? true
	// Is context.Canceled tmp? false
}
