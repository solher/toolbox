package graphql

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/solher/toolbox"
	"github.com/vektah/gqlparser/gqlerror"
)

var (
	// ErrInternal indicates an unexpected internal error.
	ErrInternal = gqlerror.Error{
		Message: "An internal error occured. Please retry later.",
		Extensions: map[string]interface{}{
			"errorCode": "INTERNAL_ERROR",
		},
	}
	// ErrUnavailable indicates that the desired service is unavailable.
	ErrUnavailable = gqlerror.Error{
		Message: "The service is currently unavailable. Please retry later.",
		Extensions: map[string]interface{}{
			"errorCode": "SERVICE_UNAVAILABLE",
		},
	}
	// ErrValidation indicates that some received parameters are invalid.
	ErrValidation = gqlerror.Error{
		Message: "The parameters validation failed.",
		Extensions: map[string]interface{}{
			"errorCode": "VALIDATION_ERROR",
		},
	}
	// ErrNotFound indicates that the requested resource was not found.
	ErrNotFound = gqlerror.Error{
		Message: "The specified resource was not found.",
		Extensions: map[string]interface{}{
			"errorCode": "NOT_FOUND",
		},
	}
)

// ErrorGenerator generates a GraphQL error.
type ErrorGenerator func(ctx context.Context, gqlErr gqlerror.Error, e error) error

// NewErrorGenerator returns an GraphQL error generator.
func NewErrorGenerator(logger log.Logger, debug bool) ErrorGenerator {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return func(ctx context.Context, gqlErr gqlerror.Error, e error) error {
		if e == nil {
			e = errors.New("null")
		}

		if gqlErrCode := gqlErr.Extensions["errorCode"].(string); debug || gqlErrCode == "INTERNAL_ERROR" || gqlErrCode == "SERVICE_UNAVAILABLE" {
			toolbox.LoggerWithRequestContext(ctx, logger).Log("code", gqlErrCode, "err", e)
		}

		gqlErr.Extensions["err"] = e.Error()
		if debug {
			location, _ := toolbox.HasStack(e)
			gqlErr.Extensions["location"] = location
		}

		return &gqlErr
	}
}
