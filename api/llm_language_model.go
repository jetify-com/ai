package api

import (
	"context"
	"iter"
	"net/http"
	"time"
)

// LanguageModel represents a language model.
type LanguageModel interface {
	// ProviderName returns the name of the provider for logging purposes.
	ProviderName() string

	// ModelID returns the provider-specific model ID for logging purposes.
	ModelID() string

	// SupportedUrls returns URL patterns supported by the model, grouped by media type.
	//
	// The MediaType field contains media type patterns or full media types (e.g. "*/*" for everything,
	// "audio/*", "video/*", or "application/pdf"). The URLPatterns field contains arrays of regular
	// expression patterns that match URL paths.
	//
	// The matching is performed against lowercase URLs.
	//
	// URLs that match these patterns are supported natively by the model and will not
	// be downloaded by the SDK. For non-matching URLs, the SDK will download the content
	// and pass it directly to the model.
	//
	// If nil or an empty slice is returned, the SDK will download all files.
	SupportedUrls() []SupportedURL

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

// SupportedURL defines URL patterns supported for a specific media type
type SupportedURL struct {
	// MediaType is the IANA media type (mime type) of the URL.
	// A simple '*' wildcard is supported for the mime type or subtype
	// (e.g., "application/pdf", "audio/*", "*/*").
	MediaType string

	// TODO: change reasoning to support reasoning blocks as well
	// URLPatterns contains regex patterns for URL paths that match this media type
	URLPatterns []string
}

// Response represents the result of a non-streaming language model generation.
type Response struct {
	// Content contains the ordered list of content blocks that the model generated.
	Content []ContentBlock `json:"content"`

	// FinishReason contains an explanation for why the model finished generating.
	FinishReason FinishReason `json:"finish_reason"`

	// Usage contains information about the number of tokens used by the model.
	Usage Usage `json:"usage"`

	// Additional provider-specific metadata. They are passed through from the
	// provider to enable provider-specific functionality.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`

	// RequestInfo is optional request information for telemetry and debugging purposes.
	RequestInfo *RequestInfo `json:"request,omitzero"`

	// ResponseInfo is optional response information for telemetry and debugging purposes.
	ResponseInfo *ResponseInfo `json:"response,omitzero"`

	// Warnings is a list of warnings that occurred during the call,
	// e.g. unsupported settings.
	Warnings []CallWarning `json:"warnings,omitzero"`
}

func (r Response) GetProviderMetadata() *ProviderMetadata { return r.ProviderMetadata }

// StreamResponse represents the result of a streaming language model call.
type StreamResponse struct {
	// Stream is the sequence of events received from the model.
	// Iterating over events might block if we're waiting for the LLM to respond.
	Stream iter.Seq[StreamEvent]
	// TODO: For now we're always encoding errors as ErrorEvent. Is that the right
	// behavior or should we consider iter.Seq2[StreamEvent, error]?

	// RequestInfo is optional request information for telemetry and debugging purposes.
	RequestInfo *RequestInfo `json:"request,omitzero"`

	// ResponseInfo is optional response information for telemetry and debugging purposes.
	ResponseInfo *ResponseInfo `json:"response,omitzero"`
}

// Usage represents token usage statistics for a model call.
//
// If a provider returns additional usage information besides the ones below,
// that information is added to the provider metadata field.
type Usage struct {
	// InputTokens is the number of tokens used by the input (prompt).
	InputTokens int `json:"input_tokens"`

	// OutputTokens is the number of tokens in the generated output (completion or tool call).
	OutputTokens int `json:"output_tokens"`

	// TotalTokens is the total number of tokens used as reported by the provider.
	// Note that this might be different from the sum of input tokens and output tokens
	// because it might include reasoning tokens or other overhead.
	TotalTokens int `json:"total_tokens"`

	// ReasoningTokens is the number of tokens used by model as part of the reasoning process.
	ReasoningTokens int `json:"reasoning_tokens,omitzero"`

	// CachedInputTokens is the number of input tokens that were cached from a previous call.
	CachedInputTokens int `json:"cached_input_tokens,omitzero"`
}

// IsZero returns true if all fields of the Usage struct are zero.
func (u Usage) IsZero() bool {
	return u.InputTokens == 0 &&
		u.OutputTokens == 0 &&
		u.TotalTokens == 0 &&
		u.ReasoningTokens == 0 &&
		u.CachedInputTokens == 0
}

// RequestInfo contains optional request information for telemetry.
type RequestInfo struct {
	// Body is the raw HTTP body that was sent to the provider
	Body []byte `json:"body,omitzero"`
}

// ResponseInfo contains optional response information for telemetry.
type ResponseInfo struct {
	// ID for the generated response, if the provider sends one.
	ID string `json:"id,omitzero"`

	// Timestamp for the start of the generated response, if the provider sends one.
	Timestamp time.Time `json:"timestamp,omitzero"`

	// ModelID of the model that was used to generate the response, if the provider sends one.
	ModelID string `json:"model_id,omitzero"`

	// Header contains a map of the HTTP response headers.
	Header http.Header

	// Body is the raw HTTP body that was returned by the provider.
	// Not provided for streaming responses.
	Body []byte `json:"body,omitzero"`

	// Status is a status code and message. e.g. "200 OK"
	Status string `json:"status,omitzero"`

	// StatusCode is a status code as integer. e.g. 200
	StatusCode int `json:"status_code,omitzero"`

	// TODO: consider adding a duration field
}
