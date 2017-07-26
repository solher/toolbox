package toolbox

import (
	"context"

	"github.com/go-kit/kit/log"
)

// LoggerWithStack wraps next and adds stacktrace to log entries when available.
func LoggerWithStack(next log.Logger) log.Logger {
	return &stackLogger{
		next: next,
	}
}

type stackLogger struct {
	next log.Logger
}

func (l *stackLogger) Log(keyvals ...interface{}) error {
	for i := 0; i < len(keyvals); i += 2 {
		if keyvals[i] == "err" {
			if err, ok := keyvals[i+1].(error); ok {
				if function, location := GetStack(err); location != "" {
					keyvals = append(keyvals, "location", location, "function", function)
				}
			}
		}
	}
	return l.next.Log(keyvals...)
}

// LoggerWithRequestContext wraps next and adds key values to log entries when available.
func LoggerWithRequestContext(ctx context.Context, next log.Logger) log.Logger {
	return &reqContextLogger{
		ctx:  ctx,
		next: next,
	}
}

type reqContextLogger struct {
	ctx  context.Context
	next log.Logger
}

func (l *reqContextLogger) Log(keyvals ...interface{}) error {
	if method, path, err := GetRequestContext(l.ctx); err == nil {
		keyvals = append([]interface{}{"method", method, "path", path}, keyvals...)
	}
	return l.next.Log(keyvals...)
}

// LoggerWithKeyValues wraps next and adds key values to log entries when available.
func LoggerWithKeyValues(next log.Logger) log.Logger {
	return &keyvalsLogger{
		next: next,
	}
}

type keyvalsLogger struct {
	next log.Logger
}

func (l *keyvalsLogger) Log(keyvals ...interface{}) error {
	for i := 0; i < len(keyvals); i += 2 {
		if keyvals[i] == "err" {
			if err, ok := keyvals[i+1].(error); ok {
				if newKeyvals, ok := HasKeyValues(err); ok {
					keyvals = append(keyvals, newKeyvals...)
				}
			}
		}
	}
	return l.next.Log(keyvals...)
}
