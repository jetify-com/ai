package api

// TODO: This schema package is pretty small.
// It might be best to just in line it into our AI SDK.
import "github.com/sashabaranov/go-openai/jsonschema"

// ToolChoice specifies how tools should be selected by the model.
type ToolChoice struct {
	// Type indicates how tools should be selected:
	// - "auto": tool selection is automatic (can be no tool)
	// - "none": no tool must be selected
	// - "required": one of the available tools must be selected
	// - "tool": a specific tool must be selected
	Type string `json:"type"`

	// ToolName specifies which tool to use when Type is "tool"
	ToolName string `json:"tool_name,omitzero"`
}

// ToolDefinition represents a tool that can be used in a language model call.
// It can either be a user-defined tool or a built-in provider-defined tool.
type ToolDefinition interface {
	// ToolType is the type of the tool. Either "function" or "provider-defined".
	ToolType() string
}

// FunctionTool represents a tool that has a name, description, and set of input arguments.
// Note: this is not the user-facing tool definition. The AI SDK methods will
// map the user-facing tool definitions to this format.
type FunctionTool struct {
	// Name is the unique identifier for this tool within this model call
	Name string `json:"name"`

	// Description explains the tool's purpose. The language model uses this to understand
	// the tool's purpose and provide better completion suggestions.
	Description string `json:"description,omitzero"`

	// InputSchema defines the expected inputs. The language model uses this to understand
	// the tool's input requirements and provide matching suggestions.
	// InputSchema should be defined using a JSON schema.
	InputSchema *jsonschema.Definition `json:"input_schema,omitzero"`
}

var _ ToolDefinition = &FunctionTool{}

// Type is the type of the tool (always "function")
func (t FunctionTool) ToolType() string { return "function" }

// ProviderDefinedTool represents a tool that has built-in support by the provider.
// Provider implementations will usually predefine these.
type ProviderDefinedTool interface {
	// ToolType is the type of the tool. Always "provider-defined" for provider-defined tools.
	ToolType() string

	// ID is the tool identifier in format "<provider-name>.<tool-name>"
	ID() string

	// Name returns the unique name used to identify this tool within the model's messages.
	// This is the name that will be used by the language model as the value of the ToolName field
	// in ToolCall blocks.
	Name() string

	// TODO: Consider adding a Validate method that checks if the arguments are valid and returns an error otherwise.
	// This would be used to validate the tool call arguments before sending them to the language model.
}
