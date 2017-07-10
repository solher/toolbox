package toolbox

import (
	"context"
	"net/http"
)

type key string

const (
	reqContextMethod key = "toolbox_req_context_method"
	reqContextPath   key = "toolbox_req_context_path"
)

type requestContext struct{}

// NewRequestContext returns a new RequestContext middleware.
func NewRequestContext() func(next http.Handler) http.Handler {
	l := &requestContext{}
	return l.middleware
}

func (r *requestContext) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), reqContextMethod, r.Method)
		ctx = context.WithValue(ctx, reqContextPath, r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
