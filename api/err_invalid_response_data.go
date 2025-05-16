package api

import (
	"encoding/json"
	"fmt"
)

// InvalidResponseDataError indicates that the server returned a response with invalid data content.
// This should be returned by providers when they cannot parse the response from the API.
type InvalidResponseDataError struct {
	*AISDKError

	// Data is the invalid response data that caused the error
	Data any
}

// NewInvalidResponseDataError creates a new InvalidResponseDataError instance
// Parameters:
//   - data: The invalid response data
//   - message: The error message (optional, will be auto-generated if empty)
func NewInvalidResponseDataError(data any, message string) *InvalidResponseDataError {
	if message == "" {
		dataJSON, _ := json.Marshal(data)
		message = fmt.Sprintf("Invalid response data: %s", string(dataJSON))
	}
	return &InvalidResponseDataError{
		AISDKError: NewAISDKError("AI_InvalidResponseDataError", message, nil),
		Data:       data,
	}
}
