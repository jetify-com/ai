package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		message  Message
		expected string
	}{
		{
			name: "system message",
			message: &SystemMessage{
				Content: "test content",
			},
			expected: `{"role":"system","content":"test content"}`,
		},
		{
			name: "user message with string content",
			message: &UserMessage{
				Content: UserMessageContent{
					Text: "test content",
				},
			},
			expected: `{"role":"user","content":"test content"}`,
		},
		{
			name: "user message with content parts",
			message: &UserMessage{
				Content: UserMessageContent{
					Parts: []ContentPart{
						&TextPart{
							Text: "test text",
						},
						&ImagePart{
							ImageURL: struct {
								URL string `json:"url"`
							}{
								URL: "http://example.com/image.jpg",
							},
						},
					},
				},
			},
			expected: `{"role":"user","content":[{"type":"text","text":"test text"},{"type":"image_url","image_url":{"url":"http://example.com/image.jpg"}}]}`,
		},
		{
			name: "assistant message with tool calls",
			message: &AssistantMessage{
				Content: "test content",
				ToolCalls: []ToolCall{
					{
						Type: "function",
						ID:   "call_123",
						Function: struct {
							Name      string `json:"name"`
							Arguments string `json:"arguments"`
						}{
							Name:      "test_function",
							Arguments: `{"arg":"value"}`,
						},
					},
				},
			},
			expected: `{"role":"assistant","content":"test content","tool_calls":[{"type":"function","id":"call_123","function":{"name":"test_function","arguments":"{\"arg\":\"value\"}"}}]}`,
		},
		{
			name: "tool message",
			message: &ToolMessage{
				Content:    "test content",
				ToolCallID: "call_123",
			},
			expected: `{"role":"tool","content":"test content","tool_call_id":"call_123"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test direct marshaling matches expected JSON
			data, err := json.Marshal(tt.message)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))

			// Unmarshaling should match the original struct
			unmarshaled, err := UnmarshalMessage(data)
			assert.NoError(t, err)
			assert.Equal(t, tt.message, unmarshaled)

			// Remarshaling should match the original JSON
			remarshaled, err := json.Marshal(unmarshaled)
			assert.NoError(t, err)
			assert.JSONEq(t, string(data), string(remarshaled))
		})
	}
}

func TestMessageUnmarshalingWithInvalidRoles(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		targetMsg   Message
		expectedErr string
	}{
		{
			name:        "system message with invalid role",
			json:        `{"role":"user","content":"test content"}`,
			targetMsg:   &SystemMessage{},
			expectedErr: "invalid role for SystemMessage: user",
		},
		{
			name:        "user message with invalid role",
			json:        `{"role":"system","content":"test content"}`,
			targetMsg:   &UserMessage{},
			expectedErr: "invalid role for UserMessage: system",
		},
		{
			name:        "assistant message with invalid role",
			json:        `{"role":"user","content":"test content"}`,
			targetMsg:   &AssistantMessage{},
			expectedErr: "invalid role for AssistantMessage: user",
		},
		{
			name:        "tool message with invalid role",
			json:        `{"role":"assistant","content":"test content","tool_call_id":"123"}`,
			targetMsg:   &ToolMessage{},
			expectedErr: "invalid role for ToolMessage: assistant",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(tt.json), tt.targetMsg)
			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}

func TestContentPartUnmarshalingWithInvalidTypes(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		targetPart  ContentPart
		expectedErr string
	}{
		{
			name:        "text part with invalid type",
			json:        `{"type":"image_url","text":"test content"}`,
			targetPart:  &TextPart{},
			expectedErr: "invalid type for TextPart: image_url",
		},
		{
			name:        "image part with invalid type",
			json:        `{"type":"text","image_url":{"url":"http://example.com/image.jpg"}}`,
			targetPart:  &ImagePart{},
			expectedErr: "invalid type for ImagePart: text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(tt.json), tt.targetPart)
			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}
