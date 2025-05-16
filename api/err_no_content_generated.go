package api

// NoContentGeneratedError is returned when the AI provider fails to generate any content
type NoContentGeneratedError struct {
	*AISDKError
}

// NewNoContentGeneratedError creates a new NoContentGeneratedError instance
// Parameters:
//   - message: The error message (optional, defaults to "No content generated.")
func NewNoContentGeneratedError(message string) *NoContentGeneratedError {
	if message == "" {
		message = "No content generated."
	}
	return &NoContentGeneratedError{
		AISDKError: NewAISDKError("AI_NoContentGeneratedError", message, nil),
	}
}
