package api

import (
	"context"
	"encoding/xml"
	"errors"
	"net/http"

	"github.com/go-kit/log"

	"github.com/solher/toolbox"
)

// NewXML returns a new XML instance.
func NewXML(logger log.Logger, debug bool) *XML {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &XML{
		logger: logger,
		debug:  debug,
	}
}

// XML returns a new XML encoder that logs errors.
// If debug is set, it sets the stacktrace into the HTTP body responses.
type XML struct {
	logger log.Logger
	debug  bool
}

// RenderError renders a HTTPError and logs it if it's a 500.
func (x *XML) RenderError(ctx context.Context, w http.ResponseWriter, httpError HTTPError, e error) {
	if e == nil {
		e = errors.New("null")
	}

	// We log errors and we export them to Sentry.
	if x.debug || (httpError.Status >= 500 && httpError.Status < 600) {
		logger := toolbox.LoggerWithRequestContext(ctx, x.logger)
		logger = toolbox.LoggerWithSentry(ctx, logger)
		logger.Log("status", httpError.Status, "err", e)
	}

	if x.debug {
		location, _ := toolbox.HasStack(e)
		x.renderXML(w, httpError.Status, &DebugHTTPError{
			HTTPError: httpError,
			Err:       e.Error(),
			Location:  location,
		})
	} else {
		x.renderXML(w, httpError.Status, &DebugHTTPError{
			HTTPError: httpError,
			Err:       e.Error(),
		})
	}
}

// Render renders an object to XML.
func (x *XML) Render(ctx context.Context, w http.ResponseWriter, status int, object interface{}) {
	if object == nil {
		w.WriteHeader(status)
	} else {
		x.renderXML(w, status, object)
	}
}

func (x *XML) renderXML(w http.ResponseWriter, status int, object interface{}) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(object)
}
