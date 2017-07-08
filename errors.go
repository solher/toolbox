package toolbox

// IsErrNotFound indicates if some requested resource was not found.
func IsErrNotFound(err error) bool {
	_, ok := err.(errNotFound)
	return ok
}

type errNotFound interface {
	error
	IsErrNotFound()
}

type errNotFoundBehavior struct{}

func (err errNotFoundBehavior) IsErrNotFound() {}

// WithErrNotFound wraps an error with a behavior indicating that a requested resource was not found.
func WithErrNotFound(err error) error {
	return struct {
		error
		errNotFoundBehavior
	}{
		err,
		errNotFoundBehavior{},
	}
}
