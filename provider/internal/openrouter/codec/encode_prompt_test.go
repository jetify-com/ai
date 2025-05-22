package codec

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/aitesting"
	"go.jetify.com/ai/api"
)

func TestEncodePrompt(t *testing.T) {
	tests := []struct {
		name      string
		prompt    []api.Message
		expected  string // JSON string of expected output
		wantError bool
	}{
		{
			name: "system message",
			prompt: []api.Message{
				&api.SystemMessage{Content: "test system message"},
			},
			expected: `[{"role":"system","content":"test system message"}]`,
		},
		{
			name: "user message with single text block",
			prompt: []api.Message{
				&api.UserMessage{
					Content: api.ContentFromText("hello"),
				},
			},
			expected: `[{"role":"user","content":"hello"}]`,
		},
		{
			name: "user message with multiple blocks",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "hello"},
						&api.ImageBlock{Data: []byte{0, 1, 2, 3}, MediaType: "image/png"},
						&api.FileBlock{URL: "http://example.com/file.txt"},
					},
				},
			},
			expected: `[{"role":"user","content":[
				{"type":"text","text":"hello"},
				{"type":"image_url","image_url":{"url":"data:image/png;base64,AAECAw=="}},
				{"type":"text","text":"http://example.com/file.txt"}
			]}]`,
		},
		{
			name: "assistant message with text",
			prompt: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "hello"},
					},
				},
			},
			expected: `[{"role":"assistant","content":"hello"}]`,
		},
		{
			name: "assistant message with tool calls",
			prompt: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Using calculator"},
						&api.ToolCallBlock{
							ToolCallID: "call_123",
							ToolName:   "calculator",
							Args:       json.RawMessage(`{"x": 1, "y": 2}`),
						},
					},
				},
			},
			expected: `[{"role":"assistant","content":"Using calculator","tool_calls":[
				{"type":"function","id":"call_123","function":{"name":"calculator","arguments":"{\"x\":1,\"y\":2}"}}
			]}]`,
		},
		{
			name: "tool message with multiple results",
			prompt: []api.Message{
				&api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID: "call_123",
							ToolName:   "calculator",
							Result:     json.RawMessage(`{"result": 3}`),
						},
						{
							ToolCallID: "call_456",
							ToolName:   "calculator",
							Result:     json.RawMessage(`{"result": 4}`),
						},
					},
				},
			},
			expected: `[
				{"role":"tool","content":"{\"result\":3}","tool_call_id":"call_123"},
				{"role":"tool","content":"{\"result\":4}","tool_call_id":"call_456"}
			]`,
		},
		{
			name: "user message with binary image data",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "hello"},
						api.ImageBlockFromData([]byte{0, 1, 2, 3}, "image/png"),
					},
				},
			},
			expected: `[{"role":"user","content":[
				{"type":"text","text":"hello"},
				{"type":"image_url","image_url":{"url":"data:image/png;base64,AAECAw=="}}
			]}]`,
		},
		{
			name: "user message with image URL",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "hello"},
						api.ImageBlockFromURL("https://example.com/image.jpg"),
					},
				},
			},
			expected: `[{"role":"user","content":[
				{"type":"text","text":"hello"},
				{"type":"image_url","image_url":{"url":"https://example.com/image.jpg"}}
			]}]`,
		},
		{
			name: "user message with file URL",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "hello"},
						api.FileBlockFromURL("http://example.com/file.txt"),
					},
				},
			},
			expected: `[{"role":"user","content":[
				{"type":"text","text":"hello"},
				{"type":"text","text":"http://example.com/file.txt"}
			]}]`,
		},
		{
			name: "user message with audio file data",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						api.FileBlockFromData([]byte{0, 1, 2, 3}, "audio/wav"),
					},
				},
			},
			expected: `[{"role":"user","content":[
				{"type":"text","text":"data:audio/wav;base64,AAECAw=="}
			]}]`,
		},
		{
			name: "user message with image data and missing mime type",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.ImageBlock{
							Data: []byte{0, 1, 2, 3},
							// MimeType intentionally omitted
						},
					},
				},
			},
			expected: `[{"role":"user","content":[
				{"type":"image_url","image_url":{"url":"data:image/jpeg;base64,AAECAw=="}}
			]}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages, err := EncodePrompt(tt.prompt)
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			data, err := json.Marshal(messages)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))
		})
	}
}

func TestEncodePrompt_Failures(t *testing.T) {
	tests := []struct {
		name          string
		prompt        []api.Message
		expectedError string
	}{
		{
			name: "user message with unsupported block",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&aitesting.MockUnsupportedBlock{},
					},
				},
			},
			expectedError: "unsupported content block type",
		},
		{
			name: "user message with tool call block",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.ToolCallBlock{
							ToolCallID: "call_123",
							ToolName:   "calculator",
							Args:       json.RawMessage(`{"x": 1}`),
						},
					},
				},
			},
			expectedError: "unsupported content block type",
		},
		{
			name: "assistant message with unsupported block",
			prompt: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&aitesting.MockUnsupportedBlock{},
					},
				},
			},
			expectedError: "unsupported assistant content block type",
		},
		{
			name: "assistant message with file block",
			prompt: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						api.FileBlockFromURL("http://example.com/file.txt"),
					},
				},
			},
			expectedError: "unsupported assistant content block type",
		},
		{
			name: "unsupported message type",
			prompt: []api.Message{
				&aitesting.MockUnsupportedMessage{},
			},
			expectedError: "unsupported message type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EncodePrompt(tt.prompt)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}
