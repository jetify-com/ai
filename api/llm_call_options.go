package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/tidwall/gjson"
)

// TODO: should we call it Config? Settings?
// We should think about the field name if it was being sent as JSON in a request.
// "Request" might also be a better name over "Call": "RequestOptions" or "RequestSettings"

// CallOptions represents the options for language model calls.
type CallOptions struct {
	// MaxOutputTokens specifies the maximum number of tokens to generate
	MaxOutputTokens int `json:"max_output_tokens,omitzero"`

	// Temperature controls randomness in the model's output.
	// It is recommended to set either Temperature or TopP, but not both.
	Temperature *float64 `json:"temperature,omitzero"`

	// TopP controls nucleus sampling.
	// It is recommended to set either Temperature or TopP, but not both.
	TopP float64 `json:"top_p,omitzero"`

	// TopK limits sampling to the top K options for each token.
	// Used to remove "long tail" low probability responses.
	// Recommended for advanced use cases only.
	TopK int `json:"top_k,omitzero"`

	// PresencePenalty affects the likelihood of the model repeating
	// information that is already in the prompt
	PresencePenalty float64 `json:"presence_penalty,omitzero"`

	// FrequencyPenalty affects the likelihood of the model
	// repeatedly using the same words or phrases
	FrequencyPenalty float64 `json:"frequency_penalty,omitzero"`

	// StopSequences specifies sequences that will stop generation when produced.
	// Providers may have limits on the number of stop sequences.
	StopSequences []string `json:"stop_sequences,omitempty"`

	// Seed provides an integer seed for random sampling.
	// If supported by the model, calls will generate deterministic results.
	Seed int `json:"seed,omitzero"`

	// Headers specifies additional HTTP headers to send with the request.
	// Only applicable for HTTP-based providers.
	Headers http.Header `json:"headers,omitempty"`

	// TODO: Add maxRetries at the ai package level but not provider level.

	// =====
	// Tool-related, might consider moving to a separate struct.
	// =====

	// Tools that are available for the model to use.
	Tools []ToolDefinition `json:"tools,omitempty"`

	// ToolChoice specifies how the model should select which tool to use.
	// Defaults to 'auto'.
	ToolChoice *ToolChoice `json:"tool_choice,omitzero"`

	// =====
	// Provider-specific fields, should not expose in main ai package.
	// =====

	// ResponseFormat specifies whether the output should be text or JSON.
	// For JSON output, a schema can optionally guide the model.
	ResponseFormat *ResponseFormat `json:"response_format,omitzero"` // TODO: Only for providers, not exposed.

	// ProviderMetadata contains additional provider-specific metadata.
	// The metadata is passed through to the provider from the AI SDK and enables
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`

	// TODO:
	// Include raw chunks in the stream. Only applicable for streaming calls.
	// IncludeRawChunks bool

	// TODO:
	// Do we want to let users specify a model name at this level, to let the
	// provider pick it based on a string? Or does that belong outside CallOptions
}

func (o *CallOptions) GetProviderMetadata() *ProviderMetadata { return o.ProviderMetadata }

// UnmarshalJSON implements custom JSON unmarshaling for CallOptions
// to handle the polymorphic ToolDefinition interface
func (o *CallOptions) UnmarshalJSON(data []byte) error {
	// Use a temporary struct to unmarshal everything except tools
	type CallOptionsAlias CallOptions
	temp := struct {
		*CallOptionsAlias
		Tools []json.RawMessage `json:"tools,omitempty"`
	}{
		CallOptionsAlias: &CallOptionsAlias{}, // Create new zero-value instance
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Copy back all fields except Tools (which we handle separately)
	*o = CallOptions(*temp.CallOptionsAlias)

	// Now handle tools with type discrimination using gjson for better performance
	if temp.Tools != nil {
		tools := make([]ToolDefinition, len(temp.Tools))
		for i, toolData := range temp.Tools {
			// Use gjson to extract type without full unmarshaling
			typeResult := gjson.GetBytes(toolData, "type")
			if !typeResult.Exists() {
				return fmt.Errorf("tool at index %d missing required 'type' field", i)
			}

			toolType := typeResult.String()

			// Based on type, unmarshal into appropriate concrete type
			switch toolType {
			case "function":
				var functionTool FunctionTool
				if err := json.Unmarshal(toolData, &functionTool); err != nil {
					return fmt.Errorf("failed to unmarshal function tool at index %d: %w", i, err)
				}
				tools[i] = &functionTool
			case "provider-defined":
				var providerTool ProviderDefinedTool
				if err := json.Unmarshal(toolData, &providerTool); err != nil {
					return fmt.Errorf("failed to unmarshal provider-defined tool at index %d: %w", i, err)
				}
				tools[i] = &providerTool
			default:
				return fmt.Errorf("unknown tool type '%s' at index %d", toolType, i)
			}
		}
		o.Tools = tools
	}

	return nil
}

// ResponseFormat specifies the format of the model's response.
type ResponseFormat struct {
	// Type indicates the response format type ("text" or "json")
	Type string `json:"type"`

	// Schema optionally provides a JSON schema to guide the model's output
	Schema *jsonschema.Schema `json:"schema,omitzero"`

	// Name optionally provides a name for the output to guide the model
	Name string `json:"name,omitzero"`

	// Description optionally provides additional context to guide the model
	Description string `json:"description,omitzero"`
}
