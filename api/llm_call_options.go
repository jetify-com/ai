package api

import jsonschema "github.com/sashabaranov/go-openai/jsonschema"

// TODO: should we call it Config? Settings?
// We should think about the field name if it was being sent as JSON in a request.
// "Request" might also be a better name over "Call": "RequestOptions" or "RequestSettings"

// CallOptions represents the options for language model calls.
type CallOptions struct {
	// === Generative Settings ===
	// Settings are the numerical or generative knobs that tune the model's
	// behavior such as Temperature and MaxTokens.

	// MaxTokens specifies the maximum number of tokens to generate
	MaxTokens int `json:"max_tokens,omitempty"`

	// Temperature controls randomness in the model's output.
	// It is recommended to set either Temperature or TopP, but not both.
	Temperature *float64 `json:"temperature,omitempty"`

	// StopSequences specifies sequences that will stop generation when produced.
	// Providers may have limits on the number of stop sequences.
	StopSequences []string `json:"stop_sequences,omitempty"`

	// TopP controls nucleus sampling.
	// It is recommended to set either Temperature or TopP, but not both.
	TopP float64 `json:"top_p,omitempty"`

	// TopK limits sampling to the top K options for each token.
	// Used to remove "long tail" low probability responses.
	// Recommended for advanced use cases only.
	TopK int `json:"top_k,omitempty"`

	// PresencePenalty affects the likelihood of the model repeating
	// information that is already in the prompt
	PresencePenalty float64 `json:"presence_penalty,omitempty"`

	// FrequencyPenalty affects the likelihood of the model
	// repeatedly using the same words or phrases
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`

	// ResponseFormat specifies whether the output should be text or JSON.
	// For JSON output, a schema can optionally guide the model.
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`

	// Seed provides an integer seed for random sampling.
	// If supported by the model, calls will generate deterministic results.
	Seed int `json:"seed,omitempty"`

	// Headers specifies additional HTTP headers to send with the request.
	// Only applicable for HTTP-based providers.
	Headers map[string]string `json:"headers,omitempty"`

	// InputFormat specifies whether the user provided the input as messages or as a prompt.
	// This can help guide non-chat models in the expansion, as different expansions
	// may be needed for chat vs non-chat use cases.
	InputFormat InputFormat `json:"input_format"` // "messages" or "prompt"

	// === Other Options ===

	// Mode affects the behavior of the language model. It is required to
	// support provider-independent streaming and generation of structured objects.
	// The model can take this information and e.g. configure json mode, the correct
	// low level grammar, etc. It can also be used to optimize the efficiency of the
	// streaming, e.g. tool-delta stream parts are only needed in the
	// object-tool mode.
	//
	// Mode will be removed in v2, and at that point it will be deprecated.
	// All necessary settings will be directly supported through the call settings,
	// in particular responseFormat, toolChoice, and tools.
	Mode ModeConfig `json:"mode"`

	// ProviderMetadata contains additional provider-specific metadata.
	// The metadata is passed through to the provider from the AI SDK and enables
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitempty"`
}

func (o CallOptions) GetProviderMetadata() *ProviderMetadata { return o.ProviderMetadata }

// ResponseFormat specifies the format of the model's response.
type ResponseFormat struct {
	// Type indicates the response format type ("text" or "json")
	Type string `json:"type"`

	// Schema optionally provides a JSON schema to guide the model's output
	Schema *jsonschema.Definition `json:"schema,omitempty"`

	// Name optionally provides a name for the output to guide the model
	Name string `json:"name,omitempty"`

	// Description optionally provides additional context to guide the model
	Description string `json:"description,omitempty"`
}

// ModeConfig represents the different mode configurations available for language model calls
type ModeConfig interface {
	Type() string
	modeConfig()
}

// RegularMode represents the regular mode configuration for streaming text & complete tool calls
type RegularMode struct {
	// Tools that are available for the model
	// TODO Spec V2: move to call settings
	Tools []ToolDefinition `json:"tools,omitempty"`

	// ToolChoice specifies how the tool should be selected. Defaults to 'auto'.
	// TODO Spec V2: move to call settings
	ToolChoice *ToolChoice `json:"tool_choice,omitempty"`
}

var _ ModeConfig = RegularMode{}

func (r RegularMode) Type() string {
	return "regular"
}

func (r RegularMode) modeConfig() {}

// ObjectJSONMode represents object generation with json mode enabled (streaming: text delta)
type ObjectJSONMode struct {
	// Schema is the JSON schema that the generated output should conform to
	Schema *jsonschema.Definition `json:"schema,omitempty"`

	// Name of output that should be generated. Used by some providers for additional LLM guidance.
	Name string `json:"name,omitempty"`

	// Description of the output that should be generated. Used by some providers for additional LLM guidance.
	Description string `json:"description,omitempty"`
}

var _ ModeConfig = ObjectJSONMode{}

func (o ObjectJSONMode) Type() string {
	return "object-json"
}

func (o ObjectJSONMode) modeConfig() {}

// ObjectToolMode represents object generation with tool mode enabled (streaming: tool call deltas)
type ObjectToolMode struct {
	// Tool configuration for object-tool mode
	Tool FunctionTool `json:"tool"`
}

var _ ModeConfig = ObjectToolMode{}

func (o ObjectToolMode) Type() string {
	return "object-tool"
}

func (o ObjectToolMode) modeConfig() {}

// InputFormat specifies whether the input is provided as messages or as a prompt
type InputFormat string

const (
	// InputFormatMessages indicates the input is a sequence of chat messages
	InputFormatMessages InputFormat = "messages"

	// InputFormatPrompt indicates the input is a single text prompt
	InputFormatPrompt InputFormat = "prompt"
)
