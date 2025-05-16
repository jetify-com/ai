package api

import "fmt"

// JSONParseError indicates a failure in parsing JSON
type JSONParseError struct {
	*AISDKError

	// Text is the string that failed to parse as JSON
	Text string
}

// NewJSONParseError creates a new JSONParseError instance
// Parameters:
//   - text: The text that failed to parse as JSON
//   - cause: The underlying parsing error
func NewJSONParseError(text string, cause any) *JSONParseError {
	message := fmt.Sprintf("JSON parsing failed: Text: %s.\nError message: %v", text, cause)
	return &JSONParseError{
		AISDKError: NewAISDKError("AI_JSONParseError", message, cause),
		Text:       text,
	}
}
