<!-- File is generated by "github.com/hedhyw/semerr"; DO NOT EDIT. -->

# semerr

![Version](https://img.shields.io/github/v/tag/hedhyw/semerr)
![Build Status](https://github.com/hedhyw/semerr/actions/workflows/check.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedhyw/semerr)](https://goreportcard.com/report/github.com/hedhyw/semerr)
[![Coverage Status](https://coveralls.io/repos/github/hedhyw/semerr/badge.svg?branch=main)](https://coveralls.io/github/hedhyw/semerr?branch=main)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hedhyw/semerr)](https://pkg.go.dev/github.com/hedhyw/semerr?tab=doc)

Package `semerr` helps to work with errors in Golang.

<img alr="Go Bug" src="https://raw.githubusercontent.com/ashleymcnamara/gophers/master/GO_BUG.png" width="100px">

## Status errors

Those errors are based on HTTP status names, but they are designed to be
transport-independent. For example `semerr.NewNotFoundError(err)` indicates
that something is not found
(and it is possible to extract HTTP status -> `404` and gRPC status -> `5` if required).

Small example:
```go
// Repository layer.

type RedisUserRepo struct {}

func (r RedisUserRepo) Get(ctx context.Context, id string) (entity.User, error) {
    u, err := r.client.Get(id)

    switch {
    case err == nil:
        return u, nil
    case errors.Is(err, redis.ErrNil):
        return entity.User{}, semerr.NewNotFoundError(err)
    default:
        return entity.User{}, fmt.Errorf("getting user: %w", err)
    }
}

// Domain layer.

func (c *Core) CreateOrder(ctx context.Context, order entity.Order) (err error)
    user, err := c.userRepo.GetCurrentUser(ctx)
    switch {
    case err == nil:
        // OK. Go on.
    case errors.As(err, &semerr.NotFoundError{}):
        // Repository can have any implementation and we should NOT know about
        // `sql.ErrNoRows`, `redis.Nil`, `mongo.NoKey`, so we just compare the `err` to
        // `semerr.NotFoundError`.
        //
        // We still can check `errors.Is(err, redis.Nil)` if we want,
        // because the `err` is just wrapped without any modifications!
        //
        // Also we can change meaning by rewrapping the `err`. Check the next line:
        return fmt.Errorf("getting user: %w", semerr.NewUnauthorizedError(err))
    default:
        return fmt.Errorf("getting user: %w" ,err)
    }
    
    // ...
}

// Transport layer.

func (s *Server) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    /* ... */

    err := s.core.CreateOrder(ctx, order)
    if err != nil {
        // Respond with the correct status.
        w.WriteHeader(httperr.Code(err))

        // It is better to organize a helper for `err` responding.

        return
    }

    w.WriteHeader(http.StatusOK)
}
```

## Mechanics

```go
errOriginal := errors.New("some error")
errWrapped := semerr.NewBadRequestError(errOriginal) // The text will be the same.

fmt.Println(errWrapped) // "some error"
fmt.Println(httperr.Code(errWrapped)) // http.StatusBadRequest
fmt.Println(grpcerr.Code(errWrapped)) // codes.InvalidArgument
fmt.Println(errors.Is(err, errOriginal)) // true
fmt.Println(semerr.NewBadRequestError(nil)) // nil
fmt.Println(httperr.Wrap(errOriginal, http.StatusBadRequest)) // = semerr.NewBadRequestError(errOriginal)
```

## Const error

An error that can be defined as `const`.

```go
var errMutable error = errors.New("mutable error") // Do not like this?
const errImmutable semerr.Error = "immutable error" // So use this.
```

## Also see
```go
err := errors.New("some error")

// It indicates that the server did not receive a complete request
// message within the time that it was prepared to wait.
// HTTP: Request Timeout (408); GRPC: Canceled (1).
err = semerr.NewStatusRequestTimeoutError(err)

// It indicates that the server encountered an unexpected
// condition that prevented it from fulfilling the request.
// HTTP: Internal Server Error (500); GRPC: Unknown (2).
err = semerr.NewInternalServerError(err)

// It indicates that the server cannot or will not process the
// request due to something that is perceived to be a client error.
// HTTP: Bad Request (400); GRPC: InvalidArgument (3).
err = semerr.NewBadRequestError(err)

// It indicates indicates that the origin server is refusing
// to service the request because the content is in a format
// not supported by this method on the target resource.
// HTTP: Unsupported Media Type (415); GRPC: InvalidArgument (3).
err = semerr.NewUnsupportedMediaTypeError(err)

// It indicates that the server, while acting as a gateway or
// proxy, did not receive a timely response from an upstream
// server it needed to access in order to complete the request.
// HTTP: Gateway Timeout (504); GRPC: DeadlineExceeded (4).
err = semerr.NewStatusGatewayTimeoutError(err)

// It indicates that the origin server did not find a current
// representation for the target resource or is not willing to
// disclose that one exists.
// HTTP: Not Found (404); GRPC: NotFound (5).
err = semerr.NewNotFoundError(err)

// It indicates that the request could not be completed due to
// a conflict with the current state of the target resource.
// HTTP: Conflict (409); GRPC: AlreadyExists (6).
err = semerr.NewConflictError(err)

// It indicates that the server understood the request but
// refuses to fulfill it.
// HTTP: Forbidden (403); GRPC: PermissionDenied (7).
err = semerr.NewForbiddenError(err)

// It indicates the user has sent too many requests in a given
// amount of time.
// HTTP: Too Many Requests (429); GRPC: ResourceExhausted (8).
err = semerr.NewTooManyRequestsError(err)

// It indicates that the server is refusing to process
// a request because the request content is larger than
// the server 
// HTTP: Request Entity Too Large (413); GRPC: OutOfRange (11).
err = semerr.NewRequestEntityTooLargeError(err)

// It indicates that the server does not support
// the functionality required to fulfill the request.
// HTTP: Not Implemented (501); GRPC: Unimplemented (12).
err = semerr.NewUnimplementedError(err)

// It indicates that the server is not ready to handle
// the request.
// HTTP: Service Unavailable (503); GRPC: Unavailable (14).
err = semerr.NewServiceUnavailableError(err)

// It indicates that the request has not been applied because
// it lacks valid authentication credentials for the target
// resource.
// HTTP: Unauthorized (401); GRPC: Unauthenticated (16).
err = semerr.NewUnauthorizedError(err)
```

## Contributing

Pull requests are welcomed. If you want to add a new meaning error then
edit the file
[internal/cmd/generator/errors.yaml](internal/cmd/generator/errors.yaml)
and generate a new code, for this run `make`.
