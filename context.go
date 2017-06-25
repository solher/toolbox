package toolbox

import (
	"context"

	"github.com/pkg/errors"

	"github.com/go-kit/kit/log"
)

type key string

const (
	contextLogger key = "toolboxLogger"
)

// GetLogger returns a logger from a context.
func GetLogger(ctx context.Context) (log.Logger, error) {
	if ctx == nil {
		return log.NewNopLogger(), errors.New("nil context")
	}
	logger, ok := ctx.Value(contextLogger).(log.Logger)
	if !ok {
		return log.NewNopLogger(), errors.New("unexpected type")
	}
	if logger == nil {
		return log.NewNopLogger(), errors.New("nil value in context")
	}
	return logger, nil
}
