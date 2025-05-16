package codec

import (
	"encoding/json"
	"errors"

	"github.com/anthropics/anthropic-sdk-go"
	"go.jetify.com/ai/api"
)

// DecodeResponse converts an Anthropic Message to the AI SDK Response type
func DecodeResponse(msg *anthropic.BetaMessage) (api.Response, error) {
	if msg == nil {
		return api.Response{}, errors.New("nil message provided")
	}

	response := api.Response{
		FinishReason:     decodeFinishReason(msg.StopReason),
		Usage:            decodeUsage(msg.Usage),
		ResponseInfo:     decodeResponseInfo(msg),
		ProviderMetadata: decodeProviderMetadata(msg),
	}

	content := decodeContent(msg.Content)
	response.Text = content.Text
	response.ToolCalls = content.ToolCalls
	response.Reasoning = content.Reasoning

	return response, nil
}

// decodeResponseInfo extracts the response info from an Anthropic message
func decodeResponseInfo(msg *anthropic.BetaMessage) *api.ResponseInfo {
	return &api.ResponseInfo{
		ID:      msg.ID,
		ModelID: msg.Model,
	}
}

// decodeProviderMetadata extracts Anthropic-specific metadata
func decodeProviderMetadata(msg *anthropic.BetaMessage) *api.ProviderMetadata {
	return api.NewProviderMetadata(map[string]any{
		"anthropic": &Metadata{
			Usage: Usage{
				InputTokens:              msg.Usage.InputTokens,
				OutputTokens:             msg.Usage.OutputTokens,
				CacheCreationInputTokens: msg.Usage.CacheCreationInputTokens,
				CacheReadInputTokens:     msg.Usage.CacheReadInputTokens,
			},
		},
	})
}

// responseContent holds the parsed content from an Anthropic message
type responseContent struct {
	Text      string
	ToolCalls []api.ToolCallBlock
	Reasoning []api.Reasoning
}

// decodeContent processes the content blocks from an Anthropic message
// and returns a responseContent with text, tool calls, and reasoning blocks
func decodeContent(blocks []anthropic.BetaContentBlock) responseContent {
	// Initialize with empty slices to avoid nil references
	content := responseContent{
		ToolCalls: []api.ToolCallBlock{},
		Reasoning: []api.Reasoning{},
	}

	if blocks == nil {
		return content
	}

	for _, block := range blocks {
		switch block.Type {
		case anthropic.BetaContentBlockTypeText:
			content.Text += block.Text
		case anthropic.BetaContentBlockTypeToolUse:
			content.ToolCalls = append(content.ToolCalls, decodeToolUse(block))
		case anthropic.BetaContentBlockTypeThinking, anthropic.BetaContentBlockTypeRedactedThinking:
			if reasoningBlock := decodeReasoning(block); reasoningBlock != nil {
				content.Reasoning = append(content.Reasoning, reasoningBlock)
			}
		}
	}

	return content
}

// decodeToolUse converts an Anthropic tool use block to an AI SDK FunctionToolCall
func decodeToolUse(block anthropic.BetaContentBlock) api.ToolCallBlock {
	var args string
	if block.Input != nil {
		rawArgs, err := json.Marshal(block.Input)
		if err == nil {
			args = string(rawArgs)
		} else {
			// If marshaling fails, use empty JSON object
			args = "{}"
		}
	} else {
		args = "{}"
	}

	return api.ToolCallBlock{
		ToolCallID: block.ID,
		ToolName:   block.Name,
		Args:       json.RawMessage(args),
	}
}

// decodeReasoning converts an Anthropic thinking block to an AI SDK ReasoningBlock
func decodeReasoning(block anthropic.BetaContentBlock) api.Reasoning {
	if block.Type == anthropic.BetaContentBlockTypeThinking {
		// Check for nil or empty thinking text
		if block.Thinking == "" {
			return nil
		}
		return &api.ReasoningBlock{
			Text:      block.Thinking,
			Signature: block.Signature,
		}
	} else if block.Type == anthropic.BetaContentBlockTypeRedactedThinking {
		// Check for nil or empty data
		if block.Data == "" {
			return nil
		}
		return &api.RedactedReasoningBlock{
			Data: block.Data,
		}
	}
	return nil
}

// decodeUsage converts Anthropic Usage to API SDK Usage
func decodeUsage(usage anthropic.BetaUsage) api.Usage {
	return api.Usage{
		PromptTokens:     int(usage.InputTokens),
		CompletionTokens: int(usage.OutputTokens),
	}
}

// decodeFinishReason converts an Anthropic stop reason to an AI SDK FinishReason type.
// It handles nil/empty values by returning FinishReasonUnknown.
func decodeFinishReason(finishReason anthropic.BetaMessageStopReason) api.FinishReason {
	switch finishReason {
	case anthropic.BetaMessageStopReasonEndTurn, anthropic.BetaMessageStopReasonStopSequence:
		return api.FinishReasonStop
	case anthropic.BetaMessageStopReasonToolUse:
		return api.FinishReasonToolCalls
	case anthropic.BetaMessageStopReasonMaxTokens:
		return api.FinishReasonLength
	default:
		return api.FinishReasonUnknown
	}
}
