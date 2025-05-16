package api

import (
	"net/http"
	"net/url"
)

// APICallError represents an error that occurred during an API call
type APICallError struct {
	*AISDKError

	// URL of the API endpoint that was called
	URL *url.URL

	// Request contains the original HTTP request
	Request *http.Request

	// StatusCode is the HTTP status code of the response if any
	StatusCode int

	// Response contains the original HTTP response
	Response *http.Response

	// Data contains additional error data, if any
	Data any
}

// IsRetryable indicates whether the request can be retried
// Returns true for status codes: 408 (timeout), 409 (conflict), 429 (too many requests), or 5xx (server errors)
func (e *APICallError) IsRetryable() bool {
	return e.StatusCode == http.StatusRequestTimeout || e.StatusCode == http.StatusConflict || e.StatusCode == http.StatusTooManyRequests || e.StatusCode >= 500
}

// TODO:
// - Consider providing a constructor that takes a request and response,
//   and initializes the fields from the request and response.
// - Better approach to handling Data
// - Should http.Request be a shallow copy with headers removed? (it might
//   otherwise expose sensitive information)
