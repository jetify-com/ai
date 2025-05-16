package api

// FinishReason indicates why a language model finished generating a response.
type FinishReason string

const (
	// FinishReasonStop indicates the model generated a stop sequence
	FinishReasonStop FinishReason = "stop"

	// FinishReasonLength indicates the model reached the maximum number of tokens
	FinishReasonLength FinishReason = "length"

	// FinishReasonContentFilter indicates a content filter violation stopped the model
	FinishReasonContentFilter FinishReason = "content-filter"

	// FinishReasonToolCalls indicates the model triggered tool calls
	FinishReasonToolCalls FinishReason = "tool-calls"

	// FinishReasonError indicates the model stopped because of an error
	FinishReasonError FinishReason = "error"

	// FinishReasonOther indicates the model stopped for other reasons
	FinishReasonOther FinishReason = "other"

	// FinishReasonUnknown indicates the model has not transmitted a finish reason
	FinishReasonUnknown FinishReason = "unknown"
)
