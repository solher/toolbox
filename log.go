package toolbox

import "github.com/go-kit/kit/log"

// LoggerWithStack wraps next and adds stacktrace to log entries when available.
func LoggerWithStack(next log.Logger) log.Logger {
	return &logger{
		next: next,
	}
}

type logger struct {
	next log.Logger
}

func (l *logger) Log(keyvals ...interface{}) error {
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
