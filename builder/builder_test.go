package builder

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/aitesting"
	"go.jetify.com/ai/api"
	"go.jetify.com/pkg/pointer"
)

func TestResponseBuilder(t *testing.T) {
	tests := []struct {
		name     string
		events   []api.StreamEvent
		expected api.Response
		metadata *api.StreamResponse
	}{
		{
			name: "single text delta",
			events: []api.StreamEvent{
				&api.TextDeltaEvent{TextDelta: "Hello"},
			},
			expected: api.Response{
				Text: "Hello",
			},
		},
		{
			name: "multiple text deltas",
			events: []api.StreamEvent{
				&api.TextDeltaEvent{TextDelta: "Hello "},
				&api.TextDeltaEvent{TextDelta: "World"},
			},
			expected: api.Response{
				Text: "Hello World",
			},
		},
		{
			name: "tool call",
			events: []api.StreamEvent{
				&api.ToolCallEvent{
					ToolCallID: "call_1",
					ToolName:   "test_tool",
					Args:       json.RawMessage(`["arg1", "arg2"]`),
				},
			},
			expected: api.Response{
				ToolCalls: []api.ToolCallBlock{
					{
						ToolCallID: "call_1",
						ToolName:   "test_tool",
						Args:       json.RawMessage(`["arg1", "arg2"]`),
					},
				},
			},
		},
		{
			name: "mixed events",
			events: []api.StreamEvent{
				&api.TextDeltaEvent{TextDelta: "Hello "},
				&api.ToolCallEvent{
					ToolCallID: "call_1",
					ToolName:   "test_tool",
					Args:       json.RawMessage(`["arg1", "arg2"]`),
				},
				&api.TextDeltaEvent{TextDelta: "World"},
			},
			expected: api.Response{
				Text: "Hello World",
				ToolCalls: []api.ToolCallBlock{
					{
						ToolCallID: "call_1",
						ToolName:   "test_tool",
						Args:       json.RawMessage(`["arg1", "arg2"]`),
					},
				},
			},
		},
		{
			name: "reasoning events",
			events: []api.StreamEvent{
				&api.ReasoningEvent{TextDelta: "Let's think about this"},
				&api.ReasoningSignatureEvent{Signature: "sig123"},
				&api.RedactedReasoningEvent{Data: "redacted_data"},
			},
			expected: api.Response{
				Reasoning: []api.Reasoning{
					&api.ReasoningBlock{
						Text:      "Let's think about this",
						Signature: "sig123",
					},
					&api.RedactedReasoningBlock{
						Data: "redacted_data",
					},
				},
			},
		},
		{
			name: "tool call delta",
			events: []api.StreamEvent{
				&api.ToolCallEvent{
					ToolCallID: "call_1",
					ToolName:   "test_tool",
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_1",
					ToolName:   "test_tool",
					ArgsDelta:  json.RawMessage(`["arg1",`),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_1",
					ToolName:   "test_tool",
					ArgsDelta:  json.RawMessage(`"arg2"]`),
				},
			},
			expected: api.Response{
				ToolCalls: []api.ToolCallBlock{
					{
						ToolCallID: "call_1",
						ToolName:   "test_tool",
						Args:       json.RawMessage(`["arg1","arg2"]`),
					},
				},
			},
		},
		{
			name: "source and file events",
			events: []api.StreamEvent{
				&api.SourceEvent{
					Source: api.Source{
						SourceType: "url",
						ID:         "test_source",
						URL:        "test.txt",
					},
				},
				&api.FileEvent{
					MimeType: "text/plain",
					Data:     []byte("test data"),
				},
			},
			expected: api.Response{
				Sources: []api.Source{
					{
						SourceType: "url",
						ID:         "test_source",
						URL:        "test.txt",
					},
				},
				Files: []api.FileBlock{
					{
						MimeType: "text/plain",
						Data:     []byte("test data"),
					},
				},
			},
		},
		{
			name: "metadata and finish events",
			events: []api.StreamEvent{
				&api.ResponseMetadataEvent{
					ID:        "resp_123",
					Timestamp: pointer.Ptr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ModelID:   "model_123",
				},
				&api.FinishEvent{
					FinishReason: api.FinishReason("stop"),
					Usage: &api.Usage{
						PromptTokens:     10,
						CompletionTokens: 20,
					},
				},
			},
			expected: api.Response{
				ResponseInfo: &api.ResponseInfo{
					ID:        "resp_123",
					Timestamp: pointer.Ptr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ModelID:   "model_123",
				},
				FinishReason: api.FinishReason("stop"),
				Usage: api.Usage{
					PromptTokens:     10,
					CompletionTokens: 20,
				},
			},
		},
		{
			name: "error event",
			events: []api.StreamEvent{
				&api.ErrorEvent{
					Err: "test error",
				},
			},
			expected: api.Response{},
		},
		{
			name: "add metadata",
			events: []api.StreamEvent{
				&api.TextDeltaEvent{TextDelta: "Hello"},
			},
			expected: api.Response{
				Text: "Hello",
				Warnings: []api.CallWarning{
					{
						Message: "test warning",
					},
				},
				RawCall: api.RawCall{
					RawPrompt: "test prompt",
					RawSettings: map[string]interface{}{
						"test": "setting",
					},
				},
			},
			metadata: &api.StreamResponse{
				Warnings: []api.CallWarning{
					{
						Message: "test warning",
					},
				},
				RawCall: api.RawCall{
					RawPrompt: "test prompt",
					RawSettings: map[string]interface{}{
						"test": "setting",
					},
				},
			},
		},
		{
			name: "create tool call through deltas",
			events: []api.StreamEvent{
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_1",
					ToolName:   "test_tool",
					ArgsDelta:  json.RawMessage(`{"arg1":`),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_1",
					ToolName:   "test_tool",
					ArgsDelta:  json.RawMessage(`"value1"}`),
				},
			},
			expected: api.Response{
				ToolCalls: []api.ToolCallBlock{
					{
						ToolCallID: "call_1",
						ToolName:   "test_tool",
						Args:       json.RawMessage(`{"arg1":"value1"}`),
					},
				},
			},
		},
		{
			name: "multiple tool calls through deltas",
			events: []api.StreamEvent{
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_1",
					ToolName:   "tool_1",
					ArgsDelta:  json.RawMessage(`{"arg1":`),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_2",
					ToolName:   "tool_2",
					ArgsDelta:  json.RawMessage(`{"arg2":`),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_1",
					ToolName:   "tool_1",
					ArgsDelta:  json.RawMessage(`"value1"}`),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_2",
					ToolName:   "tool_2",
					ArgsDelta:  json.RawMessage(`"value2"}`),
				},
			},
			expected: api.Response{
				ToolCalls: []api.ToolCallBlock{
					{
						ToolCallID: "call_1",
						ToolName:   "tool_1",
						Args:       json.RawMessage(`{"arg1":"value1"}`),
					},
					{
						ToolCallID: "call_2",
						ToolName:   "tool_2",
						Args:       json.RawMessage(`{"arg2":"value2"}`),
					},
				},
			},
		},
		{
			name: "interleaved tool calls and text",
			events: []api.StreamEvent{
				&api.TextDeltaEvent{TextDelta: "Starting "},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_1",
					ToolName:   "tool_1",
					ArgsDelta:  json.RawMessage(`{"arg1":`),
				},
				&api.TextDeltaEvent{TextDelta: "middle "},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_1",
					ToolName:   "tool_1",
					ArgsDelta:  json.RawMessage(`"value1"}`),
				},
				&api.TextDeltaEvent{TextDelta: "end"},
			},
			expected: api.Response{
				Text: "Starting middle end",
				ToolCalls: []api.ToolCallBlock{
					{
						ToolCallID: "call_1",
						ToolName:   "tool_1",
						Args:       json.RawMessage(`{"arg1":"value1"}`),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewResponseBuilder()

			for _, event := range tt.events {
				err := builder.AddEvent(event)
				assert.NoError(t, err)
			}

			if tt.metadata != nil {
				err := builder.AddMetadata(tt.metadata)
				assert.NoError(t, err)
			}

			resp, err := builder.Build()
			if tt.name == "error event" {
				assert.Error(t, err)
				assert.Equal(t, "test error", err.Error())
				return
			}
			assert.NoError(t, err)
			aitesting.ResponseContains(t, tt.expected, resp)
		})
	}
}
