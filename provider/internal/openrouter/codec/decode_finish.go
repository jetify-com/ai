package codec

import (
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/internal/openrouter/client"
)

// DecodeFinishReason converts an OpenRouter finish reason to an AI SDK FinishReason type.
// It handles nil/empty values by returning FinishReasonUnknown.
func DecodeFinishReason(finishReason string) api.FinishReason {
	switch finishReason {
	case client.FinishReasonStop:
		return api.FinishReasonStop
	case client.FinishReasonLength:
		return api.FinishReasonLength
	case client.FinishReasonContentFilter:
		return api.FinishReasonContentFilter
	case client.FinishReasonFunctionCall, client.FinishReasonToolCalls:
		return api.FinishReasonToolCalls
	default:
		return api.FinishReasonUnknown
	}
}
