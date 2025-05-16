package api

// AISDKError is a custom error class for AI SDK related errors.
type AISDKError struct {
	// Name is the name of the error
	Name string

	// Message is the error message
	Message string

	// Cause is the underlying cause of the error, if any
	Cause any
}

// Error implements the error interface
func (e *AISDKError) Error() string {
	return e.Message
}

// NewAISDKError creates an AI SDK Error.
// Parameters:
//   - name: The name of the error.
//   - message: The error message.
//   - cause: The underlying cause of the error.
func NewAISDKError(name string, message string, cause any) *AISDKError {
	return &AISDKError{
		Name:    name,
		Message: message,
		Cause:   cause,
	}
}
