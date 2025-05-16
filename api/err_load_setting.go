package api

// LoadSettingError indicates a failure in loading a setting
type LoadSettingError struct {
	*AISDKError
}

// NewLoadSettingError creates a new LoadSettingError instance
// Parameters:
//   - message: The error message describing why the setting failed to load
func NewLoadSettingError(message string) *LoadSettingError {
	return &LoadSettingError{
		AISDKError: NewAISDKError("AI_LoadSettingError", message, nil),
	}
}
