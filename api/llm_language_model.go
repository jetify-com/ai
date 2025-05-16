package api

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"net/url"
	"time"
)

// LanguageModel represents a language model.
type LanguageModel interface {
	// SpecificationVersion returns which language model interface version is implemented.
	// This allows evolving the interface while maintaining backwards compatibility.
	SpecificationVersion() string

	// ProviderName returns the name of the provider for logging purposes.
	ProviderName() string

	// ModelID returns the provider-specific model ID for logging purposes.
	ModelID() string

	// DefaultObjectGenerationMode returns the default object generation mode that should be used
	// when no mode is specified. Should return the mode with the best results for this model.
	// Returns empty string if object generation is not supported.
	DefaultObjectGenerationMode() ObjectGenerationMode

	// SupportsImageURLs indicates whether the model supports image URLs.
	// Defaults to true. If false, the AI SDK will download the image and
	// pass the image data to the model.
	SupportsImageURLs() bool

	// SupportsStructuredOutputs indicates whether the model supports
	// grammar-guided generation, i.e., follows JSON schemas for
	// object generation when the response format is set to "json"
	// or when the `object-json` mode is used.
	//
	// This means that the model guarantees that the generated JSON
	// will be a valid JSON object AND that the object will match the
	// JSON schema.
	//
	// Default is false.
	SupportsStructuredOutputs() bool

	// SupportsURL checks if the model supports the given URL for file parts
	// natively. If the model does not support the URL, the AI SDK will
	// download the file and pass the file data to the model.
	//
	// If not implemented or if returning false, the AI SDK will download the file.
	SupportsURL(u *url.URL) bool

	// Generate generates a language model output (non-streaming).
	//
	// The prompt parameter is a standardized prompt type, not the user-facing prompt.
	// The AI SDK methods will map the user-facing prompt types such as chat or
	// instruction prompts to this format.
	Generate(ctx context.Context, prompt []Message, opts CallOptions) (Response, error)

	// Stream generates a language model output (streaming).
	// Returns a stream of events from the model.
	//
	// The prompt parameter is a standardized prompt type, not the user-facing prompt.
	// The AI SDK methods will map the user-facing prompt types such as chat or
	// instruction prompts to this format.
	Stream(ctx context.Context, prompt []Message, opts CallOptions) (StreamResponse, error)
}

// Response represents the result of a non-streaming language model generation.
type Response struct {
	// Text that the model has generated. Can be empty if the model has only
	// generated tool calls.
	Text string `json:"text,omitempty"`

	// Reasoning text that the model has generated. Can be empty if the model
	// has only generated text.
	Reasoning []Reasoning `json:"reasoning,omitempty"`

	// Files that the model has generated as binary data.
	// Can be empty if the model didn't generate any files.
	Files []FileBlock `json:"files,omitempty"`

	// TODO: change reasoning to support reasoning blocks as well

	// Tool calls that the model has generated. Can be empty if the model has
	// only generated text.
	ToolCalls []ToolCallBlock `json:"tool_calls,omitempty"`

	// Reason why the model finished generating
	FinishReason FinishReason `json:"finish_reason"`

	// Usage contains usage information.
	Usage Usage `json:"usage"`

	// Warnings is a list of warnings that occurred during the call,
	// e.g. unsupported settings.
	Warnings []CallWarning `json:"warnings,omitempty"`

	// Additional provider-specific metadata. They are passed through from the
	// provider to enable provider-specific functionality.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitempty"`

	// Sources are the sources that have been used as input to
	// generate the response.
	Sources []Source `json:"sources,omitempty"`
	// TODO: decide if sources and citations are the same thing and if so how to
	// merge the concepts. We currently have not implemented citations but they're
	// available in anthropic for example.

	// LogProbs are the log probabilities for the completion.
	// nil if the mode does not support logprobs or if it was not enabled.
	// @deprecated will be changed into a provider-specific extension in v2
	LogProbs LogProbs `json:"logprobs,omitempty"`

	// TODO: revisit the types for the four fields below:
	// - RawCall, RawResponse, RequestInfo, ResponseInfo
	// Decide if we can use go's http request and response objects.

	// RawCall is the raw prompt and settings information for observability
	// provider integration.
	// TODO: remove in v2
	RawCall RawCall `json:"raw_call"`

	// RawResponse is optional response information for telemetry and debugging purposes.
	// TODO: rename to response in v2 or remove
	RawResponse *RawResponse `json:"raw_response,omitempty"`

	// RequestInfo is optional request information for telemetry and debugging purposes.
	RequestInfo *RequestInfo `json:"request,omitempty"`

	// ResponseInfo is optional response information for telemetry and debugging purposes.
	ResponseInfo *ResponseInfo `json:"response,omitempty"`
}

func (r Response) GetProviderMetadata() *ProviderMetadata { return r.ProviderMetadata }

// StreamResponse represents the result of a streaming language model generation.
//
// With the exception of the Events field, all other fields are for data known at
// the start of the stream.
//
// Anything that results in an update, including additional metadata, will be sent
// as a StreamEvent.
type StreamResponse struct {
	// Sequence of events received from the model.
	// Iterating over events might block if we're waiting for the LLM to respond.
	Events iter.Seq[StreamEvent]
	// TODO: For now we're always encoding errors as ErrorEvent. Is that the right
	// behavior or should we consider iter.Seq2[StreamEvent, error]?

	// RawCall is the raw prompt and settings information for observability
	// provider integration.
	// TODO: remove in v2 (there is now request info)
	RawCall RawCall `json:"raw_call"`

	// RawResponse is optional raw response data.
	// TODO: rename to response in v2
	RawResponse *RawResponse `json:"raw_response,omitempty"`

	// RequestInfo is optional request information for telemetry and debugging purposes.
	RequestInfo *RequestInfo `json:"request,omitempty"`

	// Warnings is a list of warnings that occurred during the call,
	// e.g. unsupported settings.
	Warnings []CallWarning `json:"warnings,omitempty"`

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitempty"`
}

func (r StreamResponse) GetProviderMetadata() *ProviderMetadata { return r.ProviderMetadata }

// EventType represents the different types of stream events.
type EventType string

const (
	// EventTextDelta represents an incremental text response from the model.
	EventTextDelta EventType = "text-delta"

	// EventReasoning is an optional reasoning or intermediate explanation generated by the model.
	EventReasoning EventType = "reasoning"

	// EventReasoningSignature represents a signature that verifies reasoning content.
	EventReasoningSignature EventType = "reasoning-signature"

	// EventRedactedReasoning represents redacted reasoning data.
	EventRedactedReasoning EventType = "redacted-reasoning"

	// EventSource represents a citation that was used to generate the response.
	EventSource EventType = "source"

	// EventFile represents a file generated by the model.
	EventFile EventType = "file"

	// EventToolCall represents a completed tool call with all arguments provided.
	EventToolCall EventType = "tool-call"

	// EventToolCallDelta is an incremental update for tool call arguments.
	EventToolCallDelta EventType = "tool-call-delta"

	// EventResponseMetadata contains additional response metadata, such as timestamps or provider details.
	EventResponseMetadata EventType = "response-metadata"

	// EventFinish is the final part of the stream, providing the finish reason and usage statistics.
	EventFinish EventType = "finish"

	// EventError indicates that an error occurred during the stream.
	EventError EventType = "error"

	// TODO: How should we handle refusal events? Do we need an additional event type?
)

// StreamEvent represents a streamed incremental update of the language model output.
type StreamEvent interface {
	// Type returns the type of event being received.
	Type() EventType
}

// TextDeltaEvent represents an incremental text response from the model
//
// Used to update a TextBlock incrementally.
type TextDeltaEvent struct {
	// TextDelta is a partial text response from the model
	TextDelta string `json:"text_delta"`
}

func (b TextDeltaEvent) Type() EventType { return EventTextDelta }

// ReasoningEvent represents an incremental reasoning response from the model.
//
// Used to update the text of a ReasoningBlock.
type ReasoningEvent struct {
	// TextDelta is a partial reasoning text from the model
	TextDelta string `json:"text_delta"`
}

func (b ReasoningEvent) Type() EventType { return EventReasoning }

// ReasoningSignatureEvent represents an incremental signature update for reasoning text.
//
// Used to update the signature field of a ReasoningBlock.
type ReasoningSignatureEvent struct {
	// Signature is the cryptographic signature for verifying reasoning
	Signature string `json:"signature"`
}

func (b ReasoningSignatureEvent) Type() EventType { return EventReasoningSignature }

// RedactedReasoningEvent represents an update to redacted reasoning data.
//
// Used to update the data field of a RedactedReasoningBlock.
type RedactedReasoningEvent struct {
	// Data contains redacted reasoning data
	Data string `json:"data"`
}

func (b RedactedReasoningEvent) Type() EventType { return EventRedactedReasoning }

// SourceEvent represents a source that was used to generate the response.
//
// Used to add a source to the response.
type SourceEvent struct {
	// Source contains information about the source
	Source Source `json:"source"`
}

func (b SourceEvent) Type() EventType { return EventSource }

// FileEvent represents a file generated by the model.
//
// Used to add a file to the response via a FileBlock.
type FileEvent struct {
	// MimeType is the mime type of the file
	MimeType string `json:"mime_type"`
	// Data contains the generated file as a byte array
	Data []byte `json:"data"`
}

func (b FileEvent) Type() EventType { return EventFile }

// ToolCallEvent represents a complete tool call with all arguments.
//
// Used to add a tool call to the response via a ToolCallBlock.
type ToolCallEvent struct {
	// ToolCallID is the ID of the tool call. This ID is used to match the tool call with the tool result.
	ToolCallID string `json:"tool_call_id"`

	// ToolName is the name of the tool being invoked.
	ToolName string `json:"tool_name"`

	// Args contains the arguments of the tool call as a JSON payload matching
	// the tool's input schema.
	// Note that args are often generated by the language model and may be
	// malformed.
	Args json.RawMessage `json:"args"`
}

func (b ToolCallEvent) Type() EventType { return EventToolCall }

// ToolCallDeltaEvent represents a tool call with incremental arguments.
// Tool call deltas are only needed for object generation modes.
// The tool call deltas must be partial JSON.
type ToolCallDeltaEvent struct {
	// ToolCallID is the ID of the tool call
	ToolCallID string `json:"tool_call_id"`
	// ToolCallType specifies the type of tool call (always "function")
	ToolCallType string `json:"tool_call_type"`
	// ToolName is the name of the tool being invoked
	ToolName string `json:"tool_name"`
	// ArgsDelta is a partial JSON byte slice update for the tool call arguments
	ArgsDelta []byte `json:"args_delta"`
}

func (b ToolCallDeltaEvent) Type() EventType { return EventToolCallDelta }

// ResponseMetadataEvent contains additional response metadata.
//
// It will be sent as soon as it is available, without having to wait for
// the FinishEvent.
type ResponseMetadataEvent struct {
	// ID for the generated response, if the provider sends one
	ID string `json:"id,omitempty"`
	// Timestamp represents when the stream part was generated
	Timestamp *time.Time `json:"timestamp,omitempty"`
	// ModelID for the generated response, if the provider sends one
	ModelID string `json:"model_id,omitempty"`
}

func (b ResponseMetadataEvent) Type() EventType {
	return EventResponseMetadata
}

// FinishEvent represents the final part of the stream.
//
// It will be sent once the stream has finished processing.
type FinishEvent struct {
	// FinishReason indicates why the model stopped generating
	FinishReason FinishReason `json:"finish_reason"`
	// ProviderMetadata contains provider-specific metadata
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitempty"`
	// Usage contains token usage statistics
	Usage *Usage `json:"usage,omitempty"`
	// LogProbs are the log probabilities for the completion
	// @deprecated will be changed into a provider-specific extension in v2
	LogProbs *LogProbs `json:"logprobs,omitempty"`
}

func (b FinishEvent) Type() EventType { return EventFinish }

func (b FinishEvent) GetProviderMetadata() *ProviderMetadata { return b.ProviderMetadata }

// ErrorEvent represents an error that occurred during streaming.
type ErrorEvent struct {
	// Err contains any error messages or error objects encountered during the stream
	Err any `json:"error"`

	// TODO: We might want to make sure that the error field is always serializable as JSON,
	// or we might force it to be an error type defined by our AI SDK so that the shape is
	// known if transmitted over the network.
}

func (b ErrorEvent) Type() EventType { return EventError }
func (b ErrorEvent) Error() string   { return fmt.Sprintf("%v", b.Err) }

// Usage represents token usage statistics for a model call.
type Usage struct {
	// PromptTokens is the number of tokens in the prompt
	PromptTokens int `json:"prompt_tokens"`

	// CompletionTokens is the number of tokens in the completion
	CompletionTokens int `json:"completion_tokens"`
}

// RawCall contains raw prompt and settings information for observability
// and debugging purposes.
type RawCall struct {
	// RawPrompt is the raw prompt after expansion/conversion
	RawPrompt interface{} `json:"raw_prompt"`

	// RawSettings are the raw settings used for the API call (provider-specific)
	RawSettings map[string]interface{} `json:"raw_settings"`
}

// RawResponse contains optional raw response data.
type RawResponse struct {
	// Headers are the response headers, if any
	Headers map[string]string `json:"headers,omitempty"`
}

// RequestInfo contains optional request information for telemetry.
type RequestInfo struct {
	// Body is the raw HTTP body that was sent to the provider (JSON stringified)
	Body string `json:"body,omitempty"`
}

// ResponseInfo is a placeholder for optional response information.
type ResponseInfo struct {
	// ID for the generated response, if the provider sends one.
	ID string `json:"id,omitempty"`

	// Timestamp for the start of the generated response, if the provider sends one.
	Timestamp *time.Time `json:"timestamp,omitempty"`

	// ModelID for the generated response, if the provider sends one.
	ModelID string `json:"model_id,omitempty"`
}
