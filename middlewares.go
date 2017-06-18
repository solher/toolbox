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

type loggerInjector struct {
	logger log.Logger
}

// NewLoggerInjector returns a new LoggerInjector middleware.
func NewLoggerInjector(logger log.Logger) func(next http.Handler) http.Handler {
	l := &loggerInjector{
		logger: logger,
	}
	return l.middleware
}

func (l *loggerInjector) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.With(l.logger, "method", r.Method, "path", r.URL.Path)
		ctx := context.WithValue(r.Context(), contextLogger, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
