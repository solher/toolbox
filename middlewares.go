package toolbox

import (
	"net/http"

	"context"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
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

// Logger injects a logger in each HTTP request context.
type Logger struct {
	logger log.Logger
}

// NewLogger returns a new logger instance.
func NewLogger(logger log.Logger) func(next http.Handler) http.Handler {
	l := &Logger{
		logger: logger,
	}
	return l.middleware
}

func (l *Logger) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.With(l.logger, "path", r.URL.Path, "method", r.Method)
		ctx := context.WithValue(r.Context(), contextLogger, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
