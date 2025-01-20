package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/log"

	"github.com/solher/toolbox"
)

// NewJSON returns a new JSON instance.
func NewJSON(logger log.Logger, debug bool) *JSON {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &JSON{
		logger: logger,
		debug:  debug,
	}
}

// JSON returns a new JSON encoder that logs errors.
// If debug is set, it sets the stacktrace into the HTTP body responses.
type JSON struct {
	logger log.Logger
	debug  bool
}

// RenderError renders a HTTPError and logs it if it's a 500.
func (j *JSON) RenderError(ctx context.Context, w http.ResponseWriter, httpError HTTPError, e error) {
	if e == nil {
		e = errors.New("null")
	}

	// We log errors and we export them to Sentry.
	if j.debug || (httpError.Status >= 500 && httpError.Status < 600) {
		logger := toolbox.LoggerWithRequestContext(ctx, j.logger)
		logger = toolbox.LoggerWithSentry(ctx, logger)
		logger.Log("status", httpError.Status, "err", e)
	}

	if j.debug {
		location, _ := toolbox.HasStack(e)
		j.renderJSON(w, httpError.Status, &DebugHTTPError{
			HTTPError: httpError,
			Err:       e.Error(),
			Location:  location,
		})
	} else {
		j.renderJSON(w, httpError.Status, &DebugHTTPError{
			HTTPError: httpError,
			Err:       e.Error(),
		})
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
