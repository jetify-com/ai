package api

// ObjectGenerationMode specifies how the model should generate structured objects.
type ObjectGenerationMode string

const (
	// ObjectGenerationModeNone indicates no specific object generation mode (empty string)
	ObjectGenerationModeNone ObjectGenerationMode = ""

	// ObjectGenerationModeJSON indicates the model should generate JSON directly
	ObjectGenerationModeJSON ObjectGenerationMode = "json"

	// ObjectGenerationModeTool indicates the model should use tool calls to generate objects
	ObjectGenerationModeTool ObjectGenerationMode = "tool"
)
