package toolbox

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

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

type stackTracer interface {
	StackTrace() errors.StackTrace
}

var stackFileRegex = regexp.MustCompile(`([^\(\)]+)\..+`)

// HasStack returns where the error was thrown if possible.
func HasStack(err error) (location string, ok bool) {
	if foundErr := findBehavior(err, func(err error) bool { _, ok := err.(stackTracer); return ok }); foundErr != nil {
		stackTrace := foundErr.(stackTracer).StackTrace()
		if len(stackTrace) > 0 {
			if stack := strings.SplitN(fmt.Sprintf("%+v", stackTrace[0]), "\n\t", 2); len(stack) == 2 {
				if matches := stackFileRegex.FindStringSubmatch(stack[0]); len(matches) == 2 {
					location = fmt.Sprintf("%s/%v", matches[1], stackTrace[0])
				}
			}
			return location, true
		}
	}
	return location, false
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

// HasKeyValues returns embedded key values from the error.
func HasKeyValues(err error) (keyvals []interface{}, ok bool) {
	if foundErr := findBehavior(err, func(err error) bool { _, ok := err.(keyValuer); return ok }); foundErr != nil {
		return foundErr.(keyValuer).KeyValues(), true
	}
	return nil, false
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

// IsErrNotFound indicates if some requested resource was not found.
func IsErrNotFound(err error) bool {
	if foundErr := findBehavior(err, func(err error) bool { _, ok := err.(errNotFound); return ok }); foundErr != nil {
		return true
	}
	return false
}

type errValidation interface {
	IsErrValidation()
}

type errValidationBehavior struct{}

func (err *errValidationBehavior) IsErrValidation() {}

// WithErrValidation wraps an error with a behavior indicating that some user parameters were invalid.
func WithErrValidation(err error) error {
	return struct {
		error
		*causerBehavior
		*errValidationBehavior
	}{
		err,
		&causerBehavior{cause: err},
		&errValidationBehavior{},
	}
}

// IsErrValidation indicates if some requested resource was not found.
func IsErrValidation(err error) bool {
	if foundErr := findBehavior(err, func(err error) bool { _, ok := err.(errValidation); return ok }); foundErr != nil {
		return true
	}
	return false
}

type errRetriable interface {
	IsErrRetriable()
}

type errRetriableBehavior struct{}

func (err *errRetriableBehavior) IsErrRetriable() {}

// WithErrRetriable wraps an error with a behavior indicating that the failed operation should be retried.
func WithErrRetriable(err error) error {
	return struct {
		error
		*causerBehavior
		*errRetriableBehavior
	}{
		err,
		&causerBehavior{cause: err},
		&errRetriableBehavior{},
	}
}

// IsErrRetriable indicates if some failed operation should be retried.
func IsErrRetriable(err error) bool {
	if foundErr := findBehavior(err, func(err error) bool { _, ok := err.(errRetriable); return ok }); foundErr != nil {
		return true
	}
	return false
}
