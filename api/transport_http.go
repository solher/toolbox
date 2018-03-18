package api

// HTTPError defines a standard format for HTTP errors.
type HTTPError struct {
	// The status code.
	Status int `json:"status"`
	// The description of the HTTP error.
	Description string `json:"description"`
	// The token uniquely identifying the HTTP error.
	ErrorCode string `json:"errorCode"`
	// Additional infos.
	Params map[string]interface{} `json:"params,omitempty"`
}

// DebugHTTPError defines a standard format for HTTP errors with additional debug info.
type DebugHTTPError struct {
	HTTPError
	// The err message.
	Err string `json:"err"`
	// The location where the error was thrown.
	Location string `json:"location,omitempty"`
}

var (
	// HTTPInternal indicates an unexpected internal error.
	HTTPInternal = HTTPError{
		Status:      500,
		Description: "An internal error occured. Please retry later.",
		ErrorCode:   "INTERNAL_ERROR",
		Params:      make(map[string]interface{}),
	}
	// HTTPUnavailable indicates that the desired service is unavailable.
	HTTPUnavailable = HTTPError{
		Status:      503,
		Description: "The service is currently unavailable. Please retry later.",
		ErrorCode:   "SERVICE_UNAVAILABLE",
		Params:      make(map[string]interface{}),
	}
	// HTTPBodyDecoding indicates that the request body could not be decoded (bad syntax).
	HTTPBodyDecoding = HTTPError{
		Status:      400,
		Description: "Could not decode the JSON request.",
		ErrorCode:   "BODY_DECODING_ERROR",
		Params:      make(map[string]interface{}),
	}
	// HTTPQueryParam indicates that an expected query parameter is missing.
	HTTPQueryParam = HTTPError{
		Status:      400,
		Description: "Missing query parameter.",
		ErrorCode:   "QUERY_PARAM_ERROR",
		Params:      make(map[string]interface{}),
	}
	// HTTPValidation indicates that some received parameters are invalid.
	HTTPValidation = HTTPError{
		Status:      400,
		Description: "The parameters validation failed.",
		ErrorCode:   "VALIDATION_ERROR",
		Params:      make(map[string]interface{}),
	}
	// HTTPUnauthorized indicates that the user does not have a valid associated session.
	HTTPUnauthorized = HTTPError{
		Status:      401,
		Description: "Authorization Required.",
		ErrorCode:   "AUTHORIZATION_REQUIRED",
		Params:      make(map[string]interface{}),
	}
	// HTTPForbidden indicates that the user has a valid session but he is missing some permissions.
	HTTPForbidden = HTTPError{
		Status:      403,
		Description: "The specified resource was not found or you don't have sufficient permissions.",
		ErrorCode:   "FORBIDDEN",
		Params:      make(map[string]interface{}),
	}
	// HTTPNotFound indicates that the requested resource was not found.
	HTTPNotFound = HTTPError{
		Status:      404,
		Description: "The specified resource was not found.",
		ErrorCode:   "NOT_FOUND",
		Params:      make(map[string]interface{}),
	}
)
