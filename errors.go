package edgecdnx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// APIError represents a non-2xx response from the EdgeCDN-X API.
type APIError struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int
	// Message is the error message returned by the API.
	Message string `json:"error"`
	// Reason is an optional machine-readable reason (Kubernetes status reason).
	Reason string `json:"reason,omitempty"`
	// Details carries optional structured error details.
	Details json.RawMessage `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("edgecdnx: API error %d (%s): %s", e.StatusCode, e.Reason, e.Message)
	}
	return fmt.Sprintf("edgecdnx: API error %d: %s", e.StatusCode, e.Message)
}

func newAPIError(resp *http.Response) error {
	apiErr := &APIError{StatusCode: resp.StatusCode}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err == nil && len(body) > 0 {
		if jsonErr := json.Unmarshal(body, apiErr); jsonErr != nil || apiErr.Message == "" {
			if apiErr.Message == "" {
				apiErr.Message = string(body)
			}
		}
	}
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(resp.StatusCode)
	}
	return apiErr
}

func hasStatus(err error, status int) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == status
}

// IsBadRequest reports whether err is an API error with status 400.
func IsBadRequest(err error) bool { return hasStatus(err, http.StatusBadRequest) }

// IsUnauthorized reports whether err is an API error with status 401.
func IsUnauthorized(err error) bool { return hasStatus(err, http.StatusUnauthorized) }

// IsForbidden reports whether err is an API error with status 403.
func IsForbidden(err error) bool { return hasStatus(err, http.StatusForbidden) }

// IsNotFound reports whether err is an API error with status 404.
func IsNotFound(err error) bool { return hasStatus(err, http.StatusNotFound) }

// IsConflict reports whether err is an API error with status 409.
func IsConflict(err error) bool { return hasStatus(err, http.StatusConflict) }
