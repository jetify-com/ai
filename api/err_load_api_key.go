package api

// LoadAPIKeyError indicates a failure in loading an API key
type LoadAPIKeyError struct {
	*AISDKError
}

// NewLoadAPIKeyError creates a new LoadAPIKeyError instance
// Parameters:
//   - message: The error message describing why the API key failed to load
func NewLoadAPIKeyError(message string) *LoadAPIKeyError {
	return &LoadAPIKeyError{
		AISDKError: NewAISDKError("AI_LoadAPIKeyError", message, nil),
	}
}
