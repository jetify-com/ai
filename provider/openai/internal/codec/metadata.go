package codec

import (
	"go.jetify.com/ai/api"
)

// For now we are using a single type for all metadata.
// TODO: Decide if we will need different types for different metadata.
type Metadata struct {
	// --- Used in requests ---

	// ParallelToolCalls determines whether to allow the model to run tool calls in parallel.
	// When not specified (nil), OpenAI defaults this to true.
	ParallelToolCalls *bool `json:"parallel_tool_calls,omitempty"`

	// PreviousResponseID is the unique ID of the previous response to the model. Use this to create
	// multi-turn conversations. Learn more about
	// [conversation state](https://platform.openai.com/docs/guides/conversation-state).
	PreviousResponseID string `json:"previous_response_id,omitempty"`

	// Store determines whether to store the generated model response for later retrieval via API.
	// When not specified (nil), OpenAI defaults this to true.
	Store *bool `json:"store,omitempty"`

	// User is a unique identifier representing your end-user, which can help OpenAI to monitor
	// and detect abuse.
	// [Learn more](https://platform.openai.com/docs/guides/safety-best-practices#end-user-ids)
	User string `json:"user,omitempty"`

	// Instructions is a system (or developer) message that is inserted as the first item
	// in the model's context.
	//
	// When using along with `previous_response_id`, the instructions from a previous
	// response will not be carried over to the next response. This makes it simple to
	// swap out system (or developer) messages in new responses.
	Instructions string `json:"instructions,omitempty"`

	// StrictSchemas determines whether JSON schema validation should be strict.
	// When true (default), the model will strictly follow the provided JSON schema.
	//
	// Whether to enable strict schema adherence when generating the output. If set to
	// true, the model will always follow the exact schema defined in the `schema`
	// field. Only a subset of JSON Schema is supported when `strict` is `true`. To
	// learn more, read the
	// [Structured Outputs guide](https://platform.openai.com/docs/guides/structured-outputs).
	StrictSchemas *bool `json:"strict_schemas,omitempty"`

	// ReasoningEffort controls the amount of reasoning the model puts into its response.
	// Reducing reasoning effort can result in faster responses and fewer tokens used on
	// reasoning in a response.
	//
	// Only supported by reasoning models (O series).
	//
	// Read more about [reasoning models](https://platform.openai.com/docs/guides/reasoning)
	// for more information.
	//
	// Supported values are `low`, `medium`, and `high`.
	ReasoningEffort string `json:"reasoning_effort,omitempty"`

	// ReasoningSummary indicates the level of detail that should be used when
	// summarizing the reasoning performed by the model. This can be useful for
	// debugging and understanding the model's reasoning process.
	//
	// Only supported by computer_use_preview.
	//
	// Supported values are `concise` and `detailed`.
	ReasoningSummary string `json:"reasoning_summary,omitempty"`

	// --- Used in blocks ---

	// ImageDetail indicates the level of detail that should be used when processing
	// and understanding the image that is being sent to the model.
	//
	// One of `high`, `low`, or `auto`. Defaults to `auto`.
	ImageDetail string `json:"image_detail,omitempty"`

	// Filename is the custom filename to use when sending a file to the model.
	Filename string `json:"filename,omitempty"`

	// --- Used in responses ---

	// ResponseID is the unique ID of the response.
	ResponseID string `json:"response_id,omitempty"`
	// TODO: Decide whether to promote ID to a top-level field.

	// Usage stores token usage details including input tokens, output tokens, a
	// breakdown of output tokens, and the total tokens used.
	Usage Usage `json:"usage,omitempty"`

	// ComputerSafetyChecks is a list of pending safety checks for the computer call.
	ComputerSafetyChecks []ComputerSafetyCheck `json:"computer_safety_checks,omitempty"`
}

func GetMetadata(source api.MetadataSource) *Metadata {
	return api.GetMetadata[Metadata]("openai", source)
}

// Usage stores token usage details including input tokens, output tokens, a
// breakdown of output tokens, and the total tokens used.
type Usage struct {
	// The number of input tokens.
	InputTokens int `json:"input_tokens,omitempty"`

	// The number of tokens that were retrieved from the cache.
	// [More on prompt caching](https://platform.openai.com/docs/guides/prompt-caching).
	InputCachedTokens int `json:"cached_tokens,omitempty"`

	// The number of output tokens.
	OutputTokens int `json:"output_tokens,omitempty"`

	// The number of reasoning tokens.
	OutputReasoningTokens int `json:"reasoning_tokens,omitempty"`
}
