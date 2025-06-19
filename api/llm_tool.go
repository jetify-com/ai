package api

// TODO: This schema package is pretty small.
// It might be best to just in line it into our AI SDK.
import (
	"encoding/json"

	"github.com/sashabaranov/go-openai/jsonschema"
)

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
	// Type is the type of the tool. Either "function" or "provider-defined".
	Type() string

	isToolDefinition() bool
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
func (t FunctionTool) Type() string { return "function" }

// isToolDefinition is a marker method to satisfy the ToolDefinition interface
func (t FunctionTool) isToolDefinition() bool { return true }

// FunctionTool JSON marshaling - automatically includes "type" field
func (t FunctionTool) MarshalJSON() ([]byte, error) {
	type Alias FunctionTool
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "function",
		Alias: (*Alias)(&t),
	})
}

// ProviderDefinedTool represents a tool that has built-in support by the provider.
// Provider implementations will usually predefine these.
type ProviderDefinedTool struct {
	// ID is the tool identifier in format "<provider-name>.<tool-name>"
	ID string `json:"id"`

	// Name returns the unique name used to identify this tool within the model's messages.
	// This is the name that will be used by the language model as the value of the ToolName field
	// in ToolCall blocks.
	Name string `json:"name"`

	// Args contains the arguments for configuring the tool. Must match the expected arguments
	// defined by the provider for this tool.
	// Providers should support both a JSON-serializable type and a map[string]interface{} type.
	Args any `json:"args"`
}

var _ ToolDefinition = &ProviderDefinedTool{}

// Type is the type of the tool (always "provider-defined")
func (t ProviderDefinedTool) Type() string { return "provider-defined" }

// isToolDefinition is a marker method to satisfy the ToolDefinition interface
func (t ProviderDefinedTool) isToolDefinition() bool { return true }

// ProviderDefinedTool JSON marshaling - automatically includes "type" field
func (t ProviderDefinedTool) MarshalJSON() ([]byte, error) {
	type Alias ProviderDefinedTool
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "provider-defined",
		Alias: (*Alias)(&t),
	})
}
