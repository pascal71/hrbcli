package api

import (
	"encoding/json"
	"fmt"
)

// APIError represents a Harbor API error
type APIError struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Details != nil {
		details, _ := json.Marshal(e.Details)
		return fmt.Sprintf("Harbor API error (code: %d): %s - %s", e.Code, e.Message, string(details))
	}
	return fmt.Sprintf("Harbor API error (code: %d): %s", e.Code, e.Message)
}

// IsNotFound returns true if the error is a 404
func (e *APIError) IsNotFound() bool {
	return e.Code == 404
}

// IsUnauthorized returns true if the error is a 401
func (e *APIError) IsUnauthorized() bool {
	return e.Code == 401
}

// IsForbidden returns true if the error is a 403
func (e *APIError) IsForbidden() bool {
	return e.Code == 403
}

// IsConflict returns true if the error is a 409
func (e *APIError) IsConflict() bool {
	return e.Code == 409
}

// IsBadRequest returns true if the error is a 400
func (e *APIError) IsBadRequest() bool {
	return e.Code == 400
}

// IsServerError returns true if the error is a 5xx
func (e *APIError) IsServerError() bool {
	return e.Code >= 500 && e.Code < 600
}
