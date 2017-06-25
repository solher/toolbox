package toolbox

import (
	"net/http"

	"context"

	"github.com/go-kit/kit/log"
)

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
