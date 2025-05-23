package api

import jsonschema "github.com/sashabaranov/go-openai/jsonschema"

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

	// StopSequences specifies sequences that will stop generation when produced.
	// Providers may have limits on the number of stop sequences.
	StopSequences []string `json:"stop_sequences,omitzero"`

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

	// ResponseFormat specifies whether the output should be text or JSON.
	// For JSON output, a schema can optionally guide the model.
	ResponseFormat *ResponseFormat `json:"response_format,omitzero"`

	// Seed provides an integer seed for random sampling.
	// If supported by the model, calls will generate deterministic results.
	Seed int `json:"seed,omitzero"`

	// Tools that are available for the model to use.
	Tools []ToolDefinition `json:"tools,omitzero"`

	// ToolChoice specifies how the model should select which tool to use.
	// Defaults to 'auto'.
	ToolChoice *ToolChoice `json:"tool_choice,omitzero"`

	// Headers specifies additional HTTP headers to send with the request.
	// Only applicable for HTTP-based providers.
	Headers map[string]string `json:"headers,omitzero"`

	// ProviderMetadata contains additional provider-specific metadata.
	// The metadata is passed through to the provider from the AI SDK and enables
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

func (o CallOptions) GetProviderMetadata() *ProviderMetadata { return o.ProviderMetadata }

// ResponseFormat specifies the format of the model's response.
type ResponseFormat struct {
	// Type indicates the response format type ("text" or "json")
	Type string `json:"type"`

	// Schema optionally provides a JSON schema to guide the model's output
	Schema *jsonschema.Definition `json:"schema,omitzero"`

	// Name optionally provides a name for the output to guide the model
	Name string `json:"name,omitzero"`

	// Description optionally provides additional context to guide the model
	Description string `json:"description,omitzero"`
}
