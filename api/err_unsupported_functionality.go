package api

import "fmt"

// UnsupportedFunctionalityError indicates that a requested functionality is not supported
type UnsupportedFunctionalityError struct {
	*AISDKError

	// Functionality is the name of the unsupported functionality
	Functionality string
}

// NewUnsupportedFunctionalityError creates a new UnsupportedFunctionalityError instance
// Parameters:
//   - functionality: The name of the unsupported functionality
//   - message: The error message (optional, will be auto-generated if empty)
func NewUnsupportedFunctionalityError(functionality string, message string) *UnsupportedFunctionalityError {
	if message == "" {
		message = fmt.Sprintf("'%s' functionality not supported.", functionality)
	}
	return &UnsupportedFunctionalityError{
		AISDKError:    NewAISDKError("AI_UnsupportedFunctionalityError", message, nil),
		Functionality: functionality,
	}
}
