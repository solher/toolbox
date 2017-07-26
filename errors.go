package toolbox

type causer interface {
	Cause() error
}

type causerBehavior struct {
	cause error
}

func (e *causerBehavior) Cause() error { return e.cause }

func findBehavior(err error, found func(err error) bool) error {
	for err != nil {
		if !found(err) {
			if cause, ok := err.(causer); ok {
				err = cause.Cause()
				continue
			} else {
				break
			}
		}
		return err
	}
	return nil
}

// HasKeyValues returns embedded key values from the error.
func HasKeyValues(err error) (keyvals []interface{}, ok bool) {
	if foundErr := findBehavior(err, func(err error) bool { _, ok := err.(keyValuer); return ok }); foundErr != nil {
		return foundErr.(keyValuer).KeyValues(), true
	}
	return nil, false
}

type keyValuer interface {
	KeyValues() []interface{}
}

type keyValuerBehavior struct {
	keyvals []interface{}
}

func (err *keyValuerBehavior) KeyValues() []interface{} { return err.keyvals }

// WithKeyValues wraps an error with some key values.
func WithKeyValues(err error, keyvals ...interface{}) error {
	return struct {
		error
		*causerBehavior
		*keyValuerBehavior
	}{
		err,
		&causerBehavior{cause: err},
		&keyValuerBehavior{keyvals: keyvals},
	}
}

// IsErrNotFound indicates if some requested resource was not found.
func IsErrNotFound(err error) bool {
	if foundErr := findBehavior(err, func(err error) bool { _, ok := err.(errNotFound); return ok }); foundErr != nil {
		return true
	}
	return false
}

type errNotFound interface {
	IsErrNotFound()
}

type errNotFoundBehavior struct{}

func (err *errNotFoundBehavior) IsErrNotFound() {}

// WithErrNotFound wraps an error with a behavior indicating that a requested resource was not found.
func WithErrNotFound(err error) error {
	return struct {
		error
		*causerBehavior
		*errNotFoundBehavior
	}{
		err,
		&causerBehavior{cause: err},
		&errNotFoundBehavior{},
	}
}
