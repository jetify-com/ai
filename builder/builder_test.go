package builder

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/aitesting"
	"go.jetify.com/ai/api"
)

func TestResponseBuilder(t *testing.T) {
	tests := []struct {
		name     string
		events   []api.StreamEvent
		expected *api.Response
		metadata *api.StreamResponse
	}{
		{
			name: "single text delta",
			events: []api.StreamEvent{
				&api.TextDeltaEvent{TextDelta: "Hello"},
			},
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Hello"},
				},
			},
		},
		{
			name: "multiple text deltas",
			events: []api.StreamEvent{
				&api.TextDeltaEvent{TextDelta: "Hello "},
				&api.TextDeltaEvent{TextDelta: "World"},
			},
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Hello World"},
				},
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
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
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
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Hello "},
					&api.ToolCallBlock{
						ToolCallID: "call_1",
						ToolName:   "test_tool",
						Args:       json.RawMessage(`["arg1", "arg2"]`),
					},
					&api.TextBlock{Text: "World"},
				},
			},
		},
		{
			name: "reasoning events",
			events: []api.StreamEvent{
				&api.ReasoningEvent{TextDelta: "Let's think about this"},
				&api.ReasoningSignatureEvent{Signature: "sig123"},
			},
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.ReasoningBlock{
						Text:      "Let's think about this",
						Signature: "sig123",
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
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
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
					MediaType: "text/plain",
					Data:      []byte("test data"),
				},
			},
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.SourceBlock{
						ID:    "test_source",
						URL:   "test.txt",
						Title: "",
					},
					&api.FileBlock{
						MediaType: "text/plain",
						Data:      []byte("test data"),
					},
				},
			},
		},
		{
			name: "metadata and finish events",
			events: []api.StreamEvent{
				&api.ResponseMetadataEvent{
					ID:        "resp_123",
					Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					ModelID:   "model_123",
				},
				&api.FinishEvent{
					FinishReason: api.FinishReason("stop"),
					Usage: api.Usage{
						InputTokens:  10,
						OutputTokens: 20,
						TotalTokens:  30,
					},
				},
			},
			expected: &api.Response{
				ResponseInfo: &api.ResponseInfo{
					ID:        "resp_123",
					Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					ModelID:   "model_123",
				},
				FinishReason: api.FinishReason("stop"),
				Usage: api.Usage{
					InputTokens:  10,
					OutputTokens: 20,
					TotalTokens:  30,
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
			expected: &api.Response{},
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
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
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
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "call_1",
						ToolName:   "tool_1",
						Args:       json.RawMessage(`{"arg1":"value1"}`),
					},
					&api.ToolCallBlock{
						ToolCallID: "call_2",
						ToolName:   "tool_2",
						Args:       json.RawMessage(`{"arg2":"value2"}`),
					},
				},
			},
		},
		{
			name: "multiple function calls",
			events: []api.StreamEvent{
				&api.ToolCallEvent{
					ToolCallID: "call1",
					ToolName:   "get_weather",
					Args:       json.RawMessage(`{"location":"New York"}`),
				},
				&api.ToolCallEvent{
					ToolCallID: "call2",
					ToolName:   "get_time",
					Args:       json.RawMessage(`{"timezone":"EST"}`),
				},
			},
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "call1",
						ToolName:   "get_weather",
						Args:       json.RawMessage(`{"location":"New York"}`),
					},
					&api.ToolCallBlock{
						ToolCallID: "call2",
						ToolName:   "get_time",
						Args:       json.RawMessage(`{"timezone":"EST"}`),
					},
				},
			},
		},
		{
			name: "reasoning with multiple summaries",
			events: []api.StreamEvent{
				&api.ReasoningEvent{TextDelta: "First point"},
				&api.ReasoningEvent{TextDelta: "\nSecond point"},
				&api.ReasoningSignatureEvent{Signature: "sig123"},
			},
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.ReasoningBlock{
						Text:      "First point\nSecond point",
						Signature: "sig123",
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
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Starting "},
					&api.ToolCallBlock{
						ToolCallID: "call_1",
						ToolName:   "tool_1",
						Args:       json.RawMessage(`{"arg1":"value1"}`),
					},
					&api.TextBlock{Text: "middle end"},
				},
			},
		},
		{
			name: "reasoning value type",
			events: []api.StreamEvent{
				&api.ReasoningEvent{TextDelta: "Thinking..."}, // Value type
			},
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.ReasoningBlock{Text: "Thinking..."},
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

// TestResponseBuilder_ValueTypes tests that value types (not pointers) are handled correctly
func TestResponseBuilder_ValueTypes(t *testing.T) {
	tests := []struct {
		name     string
		events   []api.StreamEvent
		expected *api.Response
		wantErr  bool
	}{
		{
			name: "text delta value type",
			events: []api.StreamEvent{
				&api.TextDeltaEvent{TextDelta: "Hello"}, // Value type, not pointer
			},
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Hello"},
				},
			},
		},
		{
			name: "multiple event value types",
			events: []api.StreamEvent{
				&api.TextDeltaEvent{TextDelta: "Hello "}, // Value type
				&api.TextDeltaEvent{TextDelta: "World"},  // Value type
				&api.SourceEvent{Source: api.Source{ // Value type - sources can be mixed with text
					SourceType: "url",
					ID:         "test_source",
					URL:        "example.com",
				}},
			},
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Hello World"},
					&api.SourceBlock{
						ID:    "test_source",
						URL:   "example.com",
						Title: "",
					},
				},
			},
		},
		{
			name: "tool call value type",
			events: []api.StreamEvent{
				&api.ToolCallEvent{ // Value type
					ToolCallID: "call_1",
					ToolName:   "test_tool",
					Args:       json.RawMessage(`["arg1", "arg2"]`),
				},
			},
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "call_1",
						ToolName:   "test_tool",
						Args:       json.RawMessage(`["arg1", "arg2"]`),
					},
				},
			},
		},
		{
			name: "mixed value and pointer types",
			events: []api.StreamEvent{
				&api.TextDeltaEvent{TextDelta: "Hello "},                 // Value type
				&api.TextDeltaEvent{TextDelta: "from pointer"},           // Pointer type
				&api.FinishEvent{FinishReason: api.FinishReason("stop")}, // Value type
			},
			expected: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Hello from pointer"},
				},
				FinishReason: api.FinishReason("stop"),
			},
		},
		{
			name: "error event value type",
			events: []api.StreamEvent{
				&api.ErrorEvent{Err: "test error"}, // Value type
			},
			expected: &api.Response{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewResponseBuilder()

			for _, event := range tt.events {
				err := builder.AddEvent(event)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
			}

			resp, err := builder.Build()
			if tt.name == "error event value type" {
				assert.Error(t, err)
				assert.Equal(t, "test error", err.Error())
				return
			}
			assert.NoError(t, err)
			aitesting.ResponseContains(t, tt.expected, resp)
		})
	}
}
