# AGENTS.md

Guidance for AI coding agents working in this repository.

## What this project is

`semerr` (module `github.com/hedhyw/semerr`) is a Go library for semantic,
transport-independent errors. Errors are named after HTTP status names
(e.g. `NotFoundError`, `BadRequestError`) but carry no transport dependency:
the corresponding HTTP status code or gRPC code is extracted only at the
transport layer via helper packages. Wrapping preserves the original error
chain (`errors.Is` / `errors.As` / `errors.Join` all work).

## Public packages (import paths)

- `github.com/hedhyw/semerr/pkg/v1/semerr` — core error types and constructors:
  - `semerr.New<Name>Error(err)` wraps `err` with a semantic meaning
    (returns `nil` if `err` is `nil`).
  - `semerr.Error` — a string-based error type usable as `const`:
    `const errFoo semerr.Error = "foo"`.
  - `semerr.IsTemporaryError` and `semerr.NewMultiError` exist but are
    **deprecated** — do not use them in new code or examples
    (use `errors.Join` instead of `NewMultiError`).
- `github.com/hedhyw/semerr/pkg/v1/httperr` — HTTP mapping:
  - `httperr.Code(err)` returns the HTTP status code for an error
    (500 for unknown errors, 200 for `nil`).
  - `httperr.Wrap(err, statusCode)` wraps an error by HTTP status code.
- `github.com/hedhyw/semerr/pkg/v1/grpcerr` — gRPC mapping:
  - `grpcerr.Code(err)` returns the `google.golang.org/grpc/codes` code.
  - `grpcerr.Wrap(err, code)` wraps an error by gRPC code.

Typical usage pattern: the repository layer converts driver-specific errors
(`sql.ErrNoRows`, `redis.Nil`, …) into semantic errors; the domain layer
checks them with `errors.As(err, &semerr.NotFoundError{})`; the transport
layer responds with `httperr.Code(err)` or `grpcerr.Code(err)`.

## Code generation — read before editing

Most of this repository is generated. The single source of truth is:

- `internal/cmd/generator/errors.yaml` — declares every error: its name,
  description, HTTP status, and gRPC code.

The generator (`internal/cmd/generator/generator.go`, run via
`go generate ./semerr.go`) renders `*.tmpl` templates into:

- `pkg/v1/semerr/*_generated.go` and tests
- `pkg/v1/httperr/*_generated.go` and tests
- `pkg/v1/grpcerr/*_generated.go` and tests
- `README.md` (from `README.md.tmpl`)

Rules:

- **Never hand-edit generated files** (`*_generated*.go`, `README.md` —
  they carry a `DO NOT EDIT` header). Edit `errors.yaml` or the `*.tmpl`
  templates instead, then run `make generate`.
- To add a new error type: add an entry to `errors.yaml`, run
  `make generate`, and commit the regenerated files together with the YAML
  change. CI fails if generated output is stale (`git diff --exit-code`).

## Layout

```
semerr.go                       # go:generate entry point only
pkg/v1/{semerr,httperr,grpcerr} # public API (v1)
internal/cmd/generator/         # code generator + errors.yaml
internal/pkg/multierr/          # helper for joined (multi) errors
```

## Commands

```sh
make generate   # regenerate code and README from errors.yaml/templates
make lint       # golangci-lint (installed into ./bin on first run)
make test       # tests with coverage across pkg/...
make            # generate + lint + test
```

Go version: see `go.mod`. After changing dependencies run `go mod tidy`.

## Conventions

- Public API lives under `pkg/v1/...`; breaking changes would require a new
  major version directory, so keep the v1 API backward compatible.
- Constructors must return `nil` when given a `nil` error — preserve this
  invariant in any new code.
- Wrapped errors must keep the original message and error chain intact
  (wrapping adds meaning, never modifies the underlying error).
- Every exported error/constructor has doc comments and tests — generated
  from templates, so changes to comment style belong in the templates.
