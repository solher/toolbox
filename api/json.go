package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"

	"github.com/pkg/errors"
	"github.com/solher/toolbox/middlewares"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// NewJSON returns a new JSON instance.
func NewJSON(debug bool) *JSON {
	return &JSON{
		debug: debug,
	}
}

// JSON returns a new JSON encoder that logs errors.
// If debug is set, it sets the stacktrace into the HTTP body responses.
type JSON struct {
	debug bool
}

// RenderError renders a HTTPError and logs it if it's a 500.
func (j *JSON) RenderError(ctx context.Context, w http.ResponseWriter, httpError HTTPError, e error) {
	stack := ""
	if e, ok := e.(stackTracer); ok && len(e.StackTrace()) > 0 {
		stack = fmt.Sprintf("%+v", e.StackTrace()[0])
	}

	if httpError.Status < 500 && httpError.Status >= 600 {
		if logger, err := middlewares.GetLogger(ctx); err != nil {
			if stack != "" {
				log.With(logger, "stack", stack)
			}
			logger.Log("status", httpError.Status, "err", e)
		}
	}

	if j.debug {
		j.renderJSON(w, httpError.Status, &DebugHTTPError{HTTPError: httpError, Err: e.Error(), Stack: stack})
	} else {
		j.renderJSON(w, httpError.Status, &httpError)
	}
}

// Render renders an object to JSON.
func (j *JSON) Render(ctx context.Context, w http.ResponseWriter, status int, object interface{}) {
	if object == nil {
		w.WriteHeader(status)
	} else {
		j.renderJSON(w, status, object)
	}
}

func (j *JSON) renderJSON(w http.ResponseWriter, status int, object interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(object)
}
