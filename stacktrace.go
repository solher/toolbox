package toolbox

import (
	"fmt"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// GetStack returns where the error was thrown if possible.
func GetStack(err error) (stack string) {
	if err, ok := err.(stackTracer); ok && len(err.StackTrace()) > 0 {
		stack = fmt.Sprintf("%+v", err.StackTrace()[0])
	}
	return stack
}
