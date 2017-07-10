package toolbox

// HasErrKeyValues returns embedded key values from the error.
func HasErrKeyValues(err error) (keyvals []interface{}, ok bool) {
	errKeyvals, ok := err.(errKeyValues)
	if !ok {
		return nil, false
	}
	return errKeyvals.ErrKeyValues(), true
}

type errKeyValues interface {
	error
	ErrKeyValues() []interface{}
}

type errKeyValuesBehavior struct {
	keyvals []interface{}
}

func (err errKeyValuesBehavior) ErrKeyValues() []interface{} { return err.keyvals }

// WithErrKeyValues wraps an error with some key values.
func WithErrKeyValues(err error, keyvals ...interface{}) error {
	return struct {
		error
		errKeyValuesBehavior
	}{
		err,
		errKeyValuesBehavior{keyvals: keyvals},
	}
}

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
