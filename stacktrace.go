package toolbox

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// GetStack returns where the error was thrown if possible.
func GetStack(err error) (function, location string) {
	if err, ok := err.(stackTracer); ok && len(err.StackTrace()) > 0 {
		if stack := strings.SplitN(fmt.Sprintf("%+v", err.StackTrace()[0]), "\n\t", 2); len(stack) == 2 {
			function, location = stack[0], stack[1]
		}
	}
	return function, location
}
