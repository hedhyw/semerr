package multierr

type MultiError interface {
	error
	Unwrap() []error
}
