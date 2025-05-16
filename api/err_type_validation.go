package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

// TypeValidationError represents a type validation failure
type TypeValidationError struct {
	*AISDKError

	// Value is the value that failed validation
	Value any
}

// NewTypeValidationError creates a new TypeValidationError instance
// Parameters:
//   - value: The value that failed validation
//   - cause: The original error or cause of the validation failure
func NewTypeValidationError(value any, cause any) *TypeValidationError {
	valueJSON, _ := json.Marshal(value)
	message := fmt.Sprintf("Type validation failed: Value: %s.\nError message: %v", string(valueJSON), cause)
	return &TypeValidationError{
		AISDKError: NewAISDKError("AI_TypeValidationError", message, cause),
		Value:      value,
	}
}

// WrapTypeValidationError wraps an error into a TypeValidationError.
// If the cause is already a TypeValidationError with the same value, it returns the cause.
// Otherwise, it creates a new TypeValidationError.
// Parameters:
//   - value: The value that failed validation
//   - cause: The original error or cause of the validation failure
//
// Returns a TypeValidationError instance
func WrapTypeValidationError(value any, cause error) *TypeValidationError {
	var existingErr *TypeValidationError
	if errors.As(cause, &existingErr) && existingErr.Value == value {
		return existingErr
	}
	return NewTypeValidationError(value, cause)
}
