<!-- File is generated by "github.com/hedhyw/semerr"; DO NOT EDIT. -->

# semerr

![Version](https://img.shields.io/github/v/tag/hedhyw/semerr)
![Build Status](https://github.com/hedhyw/semerr/actions/workflows/check.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedhyw/semerr)](https://goreportcard.com/report/github.com/hedhyw/semerr)
[![Coverage Status](https://coveralls.io/repos/github/hedhyw/semerr/badge.svg?branch=main)](https://coveralls.io/github/hedhyw/semerr?branch=main)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hedhyw/semerr)](https://pkg.go.dev/github.com/hedhyw/semerr?tab=doc)

Package `semerr` helps to work with errors in Golang.

<img alr="Go Bug" src="https://raw.githubusercontent.com/ashleymcnamara/gophers/master/GO_BUG.png" width="100px">

## Const error

An error that can be defined as `const`.

```go
var errMutable error = errors.New("mutable error") // Do not like this?
const errImmutable semerr.Error = "immutable error" // So use this.
```

## TemporaryError checking

The function `semerr.IsTemporaryError` checks that error has Temporary
method and it returns true.

```go
semerr.IsTemporaryError(context.DeadlineExceeded) // true
semerr.IsTemporaryError(context.Canceled) // false
```

## Errors with meaning

```go
errOriginal := errors.New("some error")
errWrapped := semerr.NewBadRequestError(errOriginal) // The text will be the same.

fmt.Println(errWrapped) // "some error"
fmt.Println(httperr.Code(errWrapped)) // http.StatusBadRequest
fmt.Println(grpcerr.Code(errWrapped)) // codes.InvalidArgument
fmt.Println(errors.Is(err, errOriginal)) // true
fmt.Println(semerr.NewBadRequestError(nil)) // nil
```

Also see:
```go
err := errors.New("some error")

// It indicates that the server did not receive a complete request
// message within the time that it was prepared to wait.
err = semerr.NewStatusRequestTimeoutError(err)

// It indicates that the server encountered an unexpected
// condition that prevented it from fulfilling the request.
err = semerr.NewInternalServerError(err)

// It indicates that the server cannot or will not process the
// request due to something that is perceived to be a client error.
err = semerr.NewBadRequestError(err)

// It indicates indicates that the origin server is refusing
// to service the request because the content is in a format
// not supported by this method on the target resource.
err = semerr.NewUnsupportedMediaTypeError(err)

// It indicates that the server, while acting as a gateway or
// proxy, did not receive a timely response from an upstream
// server it needed to access in order to complete the request.
err = semerr.NewStatusGatewayTimeoutError(err)

// It indicates that the origin server did not find a current
// representation for the target resource or is not willing to
// disclose that one exists.
err = semerr.NewNotFoundError(err)

// It indicates that the request could not be completed due to
// a conflict with the current state of the target resource.
err = semerr.NewConflictError(err)

// It indicates that the server understood the request but
// refuses to fulfill it.
err = semerr.NewForbiddenError(err)

// It indicates the user has sent too many requests in a given
// amount of time.
err = semerr.NewTooManyRequestsError(err)

// It indicates that the server is refusing to process
// a request because the request content is larger than
// the server 
err = semerr.NewRequestEntityTooLargeError(err)

// It indicates that the server does not support
// the functionality required to fulfill the request.
err = semerr.NewUnimplementedError(err)

// It indicates that the server is not ready to handle
// the request.
err = semerr.NewServiceUnavailableError(err)

// It indicates that the request has not been applied because
// it lacks valid authentication credentials for the target
// resource.
err = semerr.NewUnauthorizedError(err)
```

## Contributing

Pull requests are welcomed. If you want to add a new meaning error then
edit the file
[internal/cmd/generator/errors.json](internal/cmd/generator/errors.json)
and generate a new code, for this run `make`.
