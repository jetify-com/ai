package api

import "fmt"

// InvalidPromptError indicates that a prompt is invalid.
// This error should be returned by providers when they cannot process a prompt.
type InvalidPromptError struct {
	*AISDKError

	// Prompt is the invalid prompt that caused the error
	Prompt any
}

// NewInvalidPromptError creates a new InvalidPromptError instance
// Parameters:
//   - prompt: The invalid prompt
//   - message: The error message describing why the prompt is invalid
//   - cause: The underlying cause of the error (optional)
func NewInvalidPromptError(prompt any, message string, cause any) *InvalidPromptError {
	fullMessage := fmt.Sprintf("Invalid prompt: %s", message)
	return &InvalidPromptError{
		AISDKError: NewAISDKError("AI_InvalidPromptError", fullMessage, cause),
		Prompt:     prompt,
	}
}
