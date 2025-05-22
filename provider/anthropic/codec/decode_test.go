package codec

import (
	"encoding/json"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/api"
)

func TestDecodeFinishReason(t *testing.T) {
	tests := []struct {
		name         string
		finishReason anthropic.BetaMessageStopReason
		want         api.FinishReason
	}{
		{
			name:         "end_turn maps to stop",
			finishReason: anthropic.BetaMessageStopReasonEndTurn,
			want:         api.FinishReasonStop,
		},
		{
			name:         "stop_sequence maps to stop",
			finishReason: anthropic.BetaMessageStopReasonStopSequence,
			want:         api.FinishReasonStop,
		},
		{
			name:         "tool_use maps to tool-calls",
			finishReason: anthropic.BetaMessageStopReasonToolUse,
			want:         api.FinishReasonToolCalls,
		},
		{
			name:         "max_tokens maps to length",
			finishReason: anthropic.BetaMessageStopReasonMaxTokens,
			want:         api.FinishReasonLength,
		},
		{
			name:         "empty string maps to unknown",
			finishReason: "",
			want:         api.FinishReasonUnknown,
		},
		{
			name:         "unknown value maps to unknown",
			finishReason: "some_unknown_reason",
			want:         api.FinishReasonUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeFinishReason(tt.finishReason)
			assert.Equal(t, tt.want, got, "decodeFinishReason(%v)", tt.finishReason)
		})
	}
}

func TestDecodeUsage(t *testing.T) {
	tests := []struct {
		name  string
		usage anthropic.BetaUsage
		want  api.Usage
	}{
		{
			name: "normal usage",
			usage: anthropic.BetaUsage{
				InputTokens:  100,
				OutputTokens: 200,
			},
			want: api.Usage{
				InputTokens:  100,
				OutputTokens: 200,
				TotalTokens:  300,
			},
		},
		{
			name:  "zero usage",
			usage: anthropic.BetaUsage{},
			want: api.Usage{
				InputTokens:  0,
				OutputTokens: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeUsage(tt.usage)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDecodeResponseInfo(t *testing.T) {
	tests := []struct {
		name string
		msg  *anthropic.BetaMessage
		want *api.ResponseInfo
	}{
		{
			name: "message with id and model",
			msg: &anthropic.BetaMessage{
				ID:    "msg_123",
				Model: "claude-3-sonnet-20240229",
			},
			want: &api.ResponseInfo{
				ID:      "msg_123",
				ModelID: "claude-3-sonnet-20240229",
			},
		},
		{
			name: "message with empty fields",
			msg:  &anthropic.BetaMessage{},
			want: &api.ResponseInfo{
				ID:      "",
				ModelID: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeResponseInfo(tt.msg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDecodeProviderMetadata(t *testing.T) {
	tests := []struct {
		name string
		msg  *anthropic.BetaMessage
		want *api.ProviderMetadata
	}{
		{
			name: "message with cache tokens",
			msg: &anthropic.BetaMessage{
				Usage: anthropic.BetaUsage{
					CacheCreationInputTokens: 50,
					CacheReadInputTokens:     25,
				},
			},
			want: api.NewProviderMetadata(map[string]any{
				"anthropic": &Metadata{
					Usage: Usage{
						CacheCreationInputTokens: 50,
						CacheReadInputTokens:     25,
					},
				},
			}),
		},
		{
			name: "message with zero cache tokens",
			msg:  &anthropic.BetaMessage{},
			want: api.NewProviderMetadata(map[string]any{
				"anthropic": &Metadata{
					Usage: Usage{},
				},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeProviderMetadata(tt.msg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDecodeReasoning(t *testing.T) {
	tests := []struct {
		name  string
		block anthropic.BetaContentBlock
		want  api.Reasoning
	}{
		{
			name: "thinking block",
			block: anthropic.BetaContentBlock{
				Type:      anthropic.BetaContentBlockTypeThinking,
				Thinking:  "This is my reasoning",
				Signature: "sig123",
			},
			want: &api.ReasoningBlock{
				Text:      "This is my reasoning",
				Signature: "sig123",
			},
		},
		{
			name: "redacted thinking block",
			block: anthropic.BetaContentBlock{
				Type: anthropic.BetaContentBlockTypeRedactedThinking,
				Data: "redacted-data",
			},
			want: &api.RedactedReasoningBlock{
				Data: "redacted-data",
			},
		},
		{
			name: "empty thinking block",
			block: anthropic.BetaContentBlock{
				Type: anthropic.BetaContentBlockTypeThinking,
			},
			want: nil,
		},
		{
			name: "empty redacted thinking block",
			block: anthropic.BetaContentBlock{
				Type: anthropic.BetaContentBlockTypeRedactedThinking,
			},
			want: nil,
		},
		{
			name: "non-reasoning block",
			block: anthropic.BetaContentBlock{
				Type: anthropic.BetaContentBlockTypeText,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeReasoning(tt.block)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDecodeToolUse(t *testing.T) {
	tests := []struct {
		name  string
		block anthropic.BetaContentBlock
		want  api.ToolCallBlock
	}{
		{
			name: "block with input",
			block: anthropic.BetaContentBlock{
				ID:    "call_123",
				Name:  "search",
				Type:  anthropic.BetaContentBlockTypeToolUse,
				Input: json.RawMessage(`{"query":"test"}`),
			},
			want: api.ToolCallBlock{
				ToolCallID: "call_123",
				ToolName:   "search",
				Args:       json.RawMessage(`{"query":"test"}`),
			},
		},
		{
			name: "block without input",
			block: anthropic.BetaContentBlock{
				ID:   "call_456",
				Name: "get_time",
				Type: anthropic.BetaContentBlockTypeToolUse,
			},
			want: api.ToolCallBlock{
				ToolCallID: "call_456",
				ToolName:   "get_time",
				Args:       json.RawMessage(`{}`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeToolUse(tt.block)
			assert.Equal(t, tt.want.ToolCallID, got.ToolCallID)
			assert.Equal(t, tt.want.ToolName, got.ToolName)
			assert.JSONEq(t, string(tt.want.Args), string(got.Args))
		})
	}
}

// TestDecodeToolUseWithMarshalError tests the decodeToolUse function when JSON marshaling fails
func TestDecodeToolUseWithMarshalError(t *testing.T) {
	// Test with malformed JSON to trigger the marshal error path
	block := anthropic.BetaContentBlock{
		ID:    "call_789",
		Name:  "error_call",
		Type:  anthropic.BetaContentBlockTypeToolUse,
		Input: json.RawMessage(`{malformed json`), // Invalid JSON
	}

	expected := api.ToolCallBlock{
		ToolCallID: "call_789",
		ToolName:   "error_call",
		Args:       json.RawMessage(`{}`),
	}

	got := decodeToolUse(block)
	assert.Equal(t, expected.ToolCallID, got.ToolCallID)
	assert.Equal(t, expected.ToolName, got.ToolName)
	assert.JSONEq(t, string(expected.Args), string(got.Args))
}

func TestDecodeContent(t *testing.T) {
	tests := []struct {
		name   string
		blocks []anthropic.BetaContentBlock
		want   responseContent
	}{
		{
			name: "multiple block types",
			blocks: []anthropic.BetaContentBlock{
				{
					Type: anthropic.BetaContentBlockTypeText,
					Text: "Hello world",
				},
				{
					Type:     anthropic.BetaContentBlockTypeThinking,
					Thinking: "Thinking process",
				},
				{
					Type:  anthropic.BetaContentBlockTypeToolUse,
					ID:    "call_789",
					Name:  "get_weather",
					Input: json.RawMessage(`{"location":"New York"}`),
				},
			},
			want: responseContent{
				Text: "Hello world",
				ToolCalls: []api.ToolCallBlock{
					{
						ToolCallID: "call_789",
						ToolName:   "get_weather",
						Args:       json.RawMessage(`{"location":"New York"}`),
					},
				},
				Reasoning: []api.Reasoning{
					&api.ReasoningBlock{
						Text: "Thinking process",
					},
				},
			},
		},
		{
			name:   "nil blocks",
			blocks: nil,
			want: responseContent{
				Text:      "",
				ToolCalls: []api.ToolCallBlock{},
				Reasoning: []api.Reasoning{},
			},
		},
		{
			name:   "empty blocks",
			blocks: []anthropic.BetaContentBlock{},
			want: responseContent{
				Text:      "",
				ToolCalls: []api.ToolCallBlock{},
				Reasoning: []api.Reasoning{},
			},
		},
		{
			name: "empty block type",
			blocks: []anthropic.BetaContentBlock{
				{
					Type: "",
					Text: "Should be skipped",
				},
			},
			want: responseContent{
				Text:      "",
				ToolCalls: []api.ToolCallBlock{},
				Reasoning: []api.Reasoning{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeContent(tt.blocks)
			assert.Equal(t, tt.want.Text, got.Text)
			assert.Equal(t, len(tt.want.ToolCalls), len(got.ToolCalls))

			// Compare tool calls individually to use JSONEq for Args
			for i, wantToolCall := range tt.want.ToolCalls {
				if i < len(got.ToolCalls) {
					gotToolCall := got.ToolCalls[i]
					assert.Equal(t, wantToolCall.ToolCallID, gotToolCall.ToolCallID)
					assert.Equal(t, wantToolCall.ToolName, gotToolCall.ToolName)
					assert.JSONEq(t, string(wantToolCall.Args), string(gotToolCall.Args))
				}
			}

			assert.Equal(t, tt.want.Reasoning, got.Reasoning)
		})
	}
}

func TestDecodeResponse(t *testing.T) {
	tests := []struct {
		name    string
		msg     *anthropic.BetaMessage
		want    api.Response
		wantErr bool
	}{
		{
			name: "full message",
			msg: &anthropic.BetaMessage{
				ID:         "msg_123",
				Model:      "claude-3",
				StopReason: anthropic.BetaMessageStopReasonEndTurn,
				Usage: anthropic.BetaUsage{
					InputTokens:  150,
					OutputTokens: 250,
				},
				Content: []anthropic.BetaContentBlock{
					{
						Type: anthropic.BetaContentBlockTypeText,
						Text: "Hello, I am Claude",
					},
				},
			},
			want: api.Response{
				Text:         "Hello, I am Claude",
				FinishReason: api.FinishReasonStop,
				Usage: api.Usage{
					InputTokens:  150,
					OutputTokens: 250,
					TotalTokens:  400,
				},
				ResponseInfo: &api.ResponseInfo{
					ID:      "msg_123",
					ModelID: "claude-3",
				},
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"anthropic": &Metadata{
						Usage: Usage{
							InputTokens:  150,
							OutputTokens: 250,
						},
					},
				}),
				ToolCalls: []api.ToolCallBlock{},
				Reasoning: []api.Reasoning{},
			},
			wantErr: false,
		},
		{
			name: "empty message",
			msg:  &anthropic.BetaMessage{},
			want: api.Response{
				FinishReason: api.FinishReasonUnknown,
				ResponseInfo: &api.ResponseInfo{},
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"anthropic": &Metadata{
						Usage: Usage{},
					},
				}),
				ToolCalls: []api.ToolCallBlock{},
				Reasoning: []api.Reasoning{},
			},
			wantErr: false,
		},
		{
			name:    "nil message",
			msg:     nil,
			want:    api.Response{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeResponse(tt.msg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Text, got.Text)
			assert.Equal(t, tt.want.FinishReason, got.FinishReason)
			assert.Equal(t, tt.want.Usage, got.Usage)
			assert.Equal(t, tt.want.ResponseInfo, got.ResponseInfo)
			assert.Equal(t, tt.want.ProviderMetadata, got.ProviderMetadata)
			assert.Equal(t, tt.want.ToolCalls, got.ToolCalls)
			assert.Equal(t, tt.want.Reasoning, got.Reasoning)
		})
	}
}
