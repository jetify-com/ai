package codec

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/openai/openai-go/responses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

func TestDecodeStreamEvents(t *testing.T) {
	tests := []struct {
		name       string
		eventJSONs []string
		want       []api.StreamEvent
	}{
		{
			name: "simple text stream",
			eventJSONs: []string{
				`{"type": "response.created", "response": {"id": "resp_123", "created_at": 1741269019, "model": "gpt-4"}}`,
				`{"type": "response.output_text.delta", "delta": "Hello world"}`,
				`{"type": "response.completed", "response": {"usage": {"input_tokens": 10, "output_tokens": 5}}}`,
			},
			want: []api.StreamEvent{
				&api.ResponseMetadataEvent{
					ID:        "resp_123",
					Timestamp: time.Date(2025, 3, 6, 13, 50, 19, 0, time.UTC),
					ModelID:   "gpt-4",
				},
				&api.TextDeltaEvent{
					TextDelta: "Hello world",
				},
				&api.FinishEvent{
					FinishReason: api.FinishReasonStop,
					Usage: api.Usage{
						InputTokens:  10,
						OutputTokens: 5,
						TotalTokens:  15,
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{
						"openai": &Metadata{
							ResponseID: "resp_123",
							Usage: Usage{
								InputTokens:  10,
								OutputTokens: 5,
							},
						},
					}),
				},
			},
		},
		{
			name: "tool call stream",
			eventJSONs: []string{
				`{"type": "response.created", "response": {"id": "resp_456", "created_at": 1741269019, "model": "gpt-4"}}`,
				`{"type": "response.output_item.added", "output_index": 0, "item": {"type": "function_call", "call_id": "call_123", "name": "get_weather", "arguments": "{\"location\":\"New York\"}"}}`,
				`{"type": "response.completed", "response": {"usage": {"input_tokens": 15, "output_tokens": 8}}}`,
			},
			want: []api.StreamEvent{
				&api.ResponseMetadataEvent{
					ID:        "resp_456",
					Timestamp: time.Date(2025, 3, 6, 13, 50, 19, 0, time.UTC),
					ModelID:   "gpt-4",
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_123",
					ToolName:   "get_weather",
					ArgsDelta:  []byte(`{"location":"New York"}`),
				},
				&api.FinishEvent{
					FinishReason: api.FinishReasonToolCalls,
					Usage: api.Usage{
						InputTokens:  15,
						OutputTokens: 8,
						TotalTokens:  23,
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{
						"openai": &Metadata{
							ResponseID: "resp_456",
							Usage: Usage{
								InputTokens:  15,
								OutputTokens: 8,
							},
						},
					}),
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Parse the JSON events
			var events []responses.ResponseStreamEventUnion
			for _, jsonStr := range testCase.eventJSONs {
				var event responses.ResponseStreamEventUnion
				err := json.Unmarshal([]byte(jsonStr), &event)
				require.NoError(t, err)
				events = append(events, event)
			}

			// Create a mock stream
			stream := newMockStreamReader(events)

			// Decode the stream
			result, err := DecodeStream(stream)
			require.NoError(t, err)

			// Collect all events from the stream
			var got []api.StreamEvent
			for event := range result.Stream {
				got = append(got, event)
			}

			// Compare events using deep equality
			assert.Equal(t, testCase.want, got)
		})
	}
}

// mockStreamReader implements the StreamReader interface for testing
type mockStreamReader struct {
	events []responses.ResponseStreamEventUnion
	index  int
	err    error
}

// newMockStreamReader creates a new mock stream reader with the given events
func newMockStreamReader(events []responses.ResponseStreamEventUnion) *mockStreamReader {
	return &mockStreamReader{
		events: events,
		index:  -1,
	}
}

// Next advances to the next event, returning true if there is one, false otherwise
func (m *mockStreamReader) Next() bool {
	m.index++
	return m.index < len(m.events)
}

// Current returns the current event
func (m *mockStreamReader) Current() responses.ResponseStreamEventUnion {
	return m.events[m.index]
}

// Err returns any error that occurred while reading the stream
func (m *mockStreamReader) Err() error {
	return m.err
}
