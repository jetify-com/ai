package api

// InvalidArgumentError indicates that a function argument is invalid
type InvalidArgumentError struct {
	*AISDKError

	// Argument is the name of the invalid argument
	Argument string
}

// NewInvalidArgumentError creates a new InvalidArgumentError instance
// Parameters:
//   - message: The error message
//   - argument: The name of the invalid argument
//   - cause: The underlying cause of the error (optional)
func NewInvalidArgumentError(message, argument string, cause any) *InvalidArgumentError {
	return &InvalidArgumentError{
		AISDKError: NewAISDKError("AI_InvalidArgumentError", message, cause),
		Argument:   argument,
	}
}
