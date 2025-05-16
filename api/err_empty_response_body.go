package api

// EmptyResponseBodyError indicates that the response body is empty
type EmptyResponseBodyError struct {
	*AISDKError
}

// NewEmptyResponseBodyError creates a new EmptyResponseBodyError instance
// Parameters:
//   - message: The error message (optional, defaults to "Empty response body")
func NewEmptyResponseBodyError(message string) *EmptyResponseBodyError {
	if message == "" {
		message = "Empty response body"
	}
	return &EmptyResponseBodyError{
		AISDKError: NewAISDKError("AI_EmptyResponseBodyError", message, nil),
	}
}
