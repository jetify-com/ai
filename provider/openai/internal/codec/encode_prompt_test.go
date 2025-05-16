package codec

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/responses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

// testCase represents a test case for the EncodePrompt function
type testCase struct {
	name             string
	input            []api.Message
	modelConfig      modelConfig // Optional, default provided if nil
	expectedMessages []string
	expectedWarnings []api.CallWarning // Optional, defaults to empty slice
	expectedError    string            // Optional, if set, error is expected
}

// Helper function to run a slice of test cases
func runTestCases(t *testing.T, tests []testCase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function
			result, err := EncodePrompt(tt.input, tt.modelConfig)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			assertMessages(t, tt.expectedMessages, result.Messages)
			assertWarnings(t, tt.expectedWarnings, result.Warnings)
		})
	}
}

// Helper function to assert messages match expected JSON
func assertMessages(t *testing.T, expected []string, actual []responses.ResponseInputItemUnionParam) {
	assert.Equal(t, len(expected), len(actual))
	for i, expectedMessageJSON := range expected {
		if i < len(actual) {
			actualMessageJSON, err := json.Marshal(actual[i])
			require.NoError(t, err)
			assert.JSONEq(t, expectedMessageJSON, string(actualMessageJSON))
		}
	}
}

// Helper function to assert warnings match
func assertWarnings(t *testing.T, expected, actual []api.CallWarning) {
	if expected == nil {
		expected = []api.CallWarning{}
	}
	assert.ElementsMatch(t, expected, actual)
}

// System message test cases
var systemMessageTests = []testCase{
	{
		name: "system message with system role",
		input: []api.Message{
			&api.SystemMessage{
				Content: "Hello",
			},
		},
		modelConfig: modelConfig{
			SystemMessageMode: "system",
		},
		expectedMessages: []string{
			`{
				"role": "system",
				"content": "Hello"
			}`,
		},
	},
	{
		name: "system message with developer role",
		input: []api.Message{
			&api.SystemMessage{
				Content: "Hello",
			},
		},
		modelConfig: modelConfig{
			SystemMessageMode: "developer",
		},
		expectedMessages: []string{
			`{
				"role": "developer",
				"content": "Hello"
			}`,
		},
	},
	{
		name: "system message removed",
		input: []api.Message{
			&api.SystemMessage{
				Content: "Hello",
			},
		},
		modelConfig: modelConfig{
			SystemMessageMode: "remove",
		},
		expectedMessages: []string{},
		expectedWarnings: []api.CallWarning{
			{
				Type:    "other",
				Message: "system messages are removed for this model",
			},
		},
	},
	{
		name: "invalid system message mode",
		input: []api.Message{
			&api.SystemMessage{
				Content: "Hello",
			},
		},
		modelConfig: modelConfig{
			SystemMessageMode: "invalid_mode",
		},
		expectedError: "unsupported system message mode: invalid_mode",
	},
	{
		name: "system message removed with warning verification",
		input: []api.Message{
			&api.SystemMessage{
				Content: "Hello",
			},
		},
		modelConfig: modelConfig{
			SystemMessageMode: "remove",
		},
		expectedMessages: []string{},
		expectedWarnings: []api.CallWarning{
			{
				Type:    "other",
				Message: "system messages are removed for this model",
			},
		},
	},
}

// User message test cases
var userMessageTests = []testCase{
	{
		name: "user message with text only",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{
						Text: "Hello",
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": [
					{
						"type": "input_text",
						"text": "Hello"
					}
				]
			}`,
		},
	},
	{
		name: "user message with text and image URL",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{
						Text: "Hello",
					},
					&api.ImageBlock{
						URL: "https://example.com/image.jpg",
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": [
					{
						"type": "input_text",
						"text": "Hello"
					},
					{
						"type": "input_image",
						"image_url": "https://example.com/image.jpg"
					}
				]
			}`,
		},
	},
	{
		name: "user message with image binary data",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.ImageBlock{
						Data:     []byte{0, 1, 2, 3},
						MimeType: "image/png",
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": [
					{
						"type": "input_image",
						"image_url": "data:image/png;base64,AAECAw=="
					}
				]
			}`,
		},
	},
	{
		name: "user message with image binary data and default mime type",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.ImageBlock{
						Data: []byte{0, 1, 2, 3},
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": [
					{
						"type": "input_image",
						"image_url": "data:image/jpeg;base64,AAECAw=="
					}
				]
			}`,
		},
	},
	{
		name: "user message with image and detail metadata",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.ImageBlock{
						Data:     []byte{0, 1, 2, 3},
						MimeType: "image/png",
						ProviderMetadata: api.NewProviderMetadata(map[string]any{
							"openai": Metadata{
								ImageDetail: "low",
							},
						}),
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": [
					{
						"type": "input_image",
						"image_url": "data:image/png;base64,AAECAw==",
						"detail": "low"
					}
				]
			}`,
		},
	},
	{
		name: "user message with PDF file",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.FileBlock{
						Data:     []byte{1, 2, 3, 4, 5},
						MimeType: "application/pdf",
						ProviderMetadata: api.NewProviderMetadata(map[string]any{
							"openai": Metadata{
								Filename: "document.pdf",
							},
						}),
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": [
					{
						"type": "input_file",
						"filename": "document.pdf",
						"file_data": "data:application/pdf;base64,AQIDBAU="
					}
				]
			}`,
		},
	},
	{
		name: "user message with PDF file and default filename",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.FileBlock{
						Data:     []byte{1, 2, 3, 4, 5},
						MimeType: "application/pdf",
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": [
					{
						"type": "input_file",
						"filename": "file.pdf",
						"file_data": "data:application/pdf;base64,AQIDBAU="
					}
				]
			}`,
		},
	},
	{
		name: "unsupported file type",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.FileBlock{
						Data:     []byte{1, 2, 3, 4, 5},
						MimeType: "text/plain",
					},
				},
			},
		},
		expectedError: "only PDF files are supported in user messages",
	},
	{
		name: "file URL instead of data",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.FileBlock{
						URL:      "https://example.com/document.pdf",
						MimeType: "application/pdf",
					},
				},
			},
		},
		expectedError: "file URLs in user messages",
	},
	{
		name: "user message with invalid image detail",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.ImageBlock{
						URL: "https://example.com/image.jpg",
						ProviderMetadata: api.NewProviderMetadata(map[string]any{
							"openai": Metadata{
								ImageDetail: "invalid",
							},
						}),
					},
				},
			},
		},
		expectedError: "encoding user message: failed to encode content block: invalid image detail level: invalid (must be one of 'high', 'low', or 'auto')",
	},
	{
		name: "user message with empty image block",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.ImageBlock{},
				},
			},
		},
		expectedError: "encoding user message: failed to encode content block: image block must have either URL or Data",
	},
	{
		name: "user message with image and high detail",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.ImageBlock{
						URL: "https://example.com/image.jpg",
						ProviderMetadata: api.NewProviderMetadata(map[string]any{
							"openai": Metadata{
								ImageDetail: "high",
							},
						}),
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": [
					{
						"type": "input_image",
						"image_url": "https://example.com/image.jpg",
						"detail": "high"
					}
				]
			}`,
		},
	},
	{
		name: "user message with image and low detail",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.ImageBlock{
						URL: "https://example.com/image.jpg",
						ProviderMetadata: api.NewProviderMetadata(map[string]any{
							"openai": Metadata{
								ImageDetail: "low",
							},
						}),
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": [
					{
						"type": "input_image",
						"image_url": "https://example.com/image.jpg",
						"detail": "low"
					}
				]
			}`,
		},
	},
	{
		name: "user message with image and auto detail",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.ImageBlock{
						URL: "https://example.com/image.jpg",
						ProviderMetadata: api.NewProviderMetadata(map[string]any{
							"openai": Metadata{
								ImageDetail: "auto",
							},
						}),
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": [
					{
						"type": "input_image",
						"image_url": "https://example.com/image.jpg",
						"detail": "auto"
					}
				]
			}`,
		},
	},
	{
		name: "user message with image data and custom mime type",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.ImageBlock{
						Data:     []byte{1, 2, 3, 4},
						MimeType: "image/png",
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": [
					{
						"type": "input_image",
						"image_url": "data:image/png;base64,AQIDBA=="
					}
				]
			}`,
		},
	},
	{
		name: "user message with empty content array",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": []
			}`,
		},
	},
	{
		name: "user message with mixed valid and invalid blocks",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{
						Text: "Valid text",
					},
					nil,
					&api.TextBlock{
						Text: "More valid text",
					},
				},
			},
		},
		expectedError: "failed to encode content block: unsupported content block type: <nil>",
	},
	{
		name: "user message with nil text block",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					(*api.TextBlock)(nil),
				},
			},
		},
		expectedError: "failed to encode content block: text block cannot be nil",
	},
}

// Assistant message test cases
var assistantMessageTests = []testCase{
	{
		name: "assistant message with text only",
		input: []api.Message{
			&api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{
						Text: "Hello",
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "assistant",
				"content": [
					{
						"type": "output_text",
						"text": "Hello"
					}
				],
				"type": "message"
			}`,
		},
	},
	{
		name: "assistant message with text and tool call",
		input: []api.Message{
			&api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{
						Text: "I will search for that information.",
					},
					&api.ToolCallBlock{
						ToolCallID: "call_123",
						ToolName:   "search",
						Args:       json.RawMessage(`{"query":"weather in San Francisco"}`),
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "assistant",
				"content": [
					{
						"type": "output_text",
						"text": "I will search for that information."
					}
				],
				"type": "message"
			}`,
			`{
				"type": "function_call",
				"call_id": "call_123",
				"name": "search",
				"arguments": "{\"query\":\"weather in San Francisco\"}"
			}`,
		},
	},
	{
		name: "assistant message with multiple tool calls",
		input: []api.Message{
			&api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "call_123",
						ToolName:   "search",
						Args:       json.RawMessage(`{"query":"weather in San Francisco"}`),
					},
					&api.ToolCallBlock{
						ToolCallID: "call_456",
						ToolName:   "calculator",
						Args:       json.RawMessage(`{"expression":"2 + 2"}`),
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"type": "function_call",
				"call_id": "call_123",
				"name": "search",
				"arguments": "{\"query\":\"weather in San Francisco\"}"
			}`,
			`{
				"type": "function_call",
				"call_id": "call_456",
				"name": "calculator",
				"arguments": "{\"expression\":\"2 + 2\"}"
			}`,
		},
	},
	{
		name: "assistant message with empty tool name",
		input: []api.Message{
			&api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "call_123",
						ToolName:   "",
						Args:       json.RawMessage(`{"query":"test"}`),
					},
				},
			},
		},
		expectedError: "tool call is missing tool name",
	},
	{
		name: "assistant message with invalid tool call args",
		input: []api.Message{
			&api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "call_123",
						ToolName:   "search",
						Args:       json.RawMessage(`{invalid json`),
					},
				},
			},
		},
		expectedError: "failed to marshal tool call arguments",
	},
	{
		name: "user message with nil content block",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					nil,
				},
			},
		},
		expectedError: "encoding user message: failed to encode content block: unsupported content block type: <nil>",
	},
	{
		name: "user message with empty text",
		input: []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{
						Text: "",
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"role": "user",
				"content": [
					{
						"type": "input_text",
						"text": ""
					}
				]
			}`,
		},
	},
	{
		name: "assistant message with nil tool call block",
		input: []api.Message{
			&api.AssistantMessage{
				Content: []api.ContentBlock{
					(*api.ToolCallBlock)(nil),
				},
			},
		},
		expectedError: "encoding tool call block: tool call block cannot be nil",
	},
	{
		name: "assistant message with empty content",
		input: []api.Message{
			&api.AssistantMessage{
				Content: []api.ContentBlock{},
			},
		},
		expectedMessages: []string{},
	},
	{
		name: "assistant message with nil text block",
		input: []api.Message{
			&api.AssistantMessage{
				Content: []api.ContentBlock{
					(*api.TextBlock)(nil),
				},
			},
		},
		expectedError: "encoding text block: text block cannot be nil",
	},
	{
		name: "assistant message with unsupported block type",
		input: []api.Message{
			&api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.ImageBlock{}, // Image blocks aren't supported in assistant messages
				},
			},
		},
		expectedError: "unsupported content block type in assistant message: *api.ImageBlock",
	},
}

// Tool message test cases
var toolMessageTests = []testCase{
	{
		name: "tool message with single result",
		input: []api.Message{
			&api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "call_123",
						ToolName:   "search",
						Result:     json.RawMessage(`{"temperature":"72°F","condition":"Sunny"}`),
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"type": "function_call_output",
				"call_id": "call_123",
				"output": "{\"temperature\":\"72°F\",\"condition\":\"Sunny\"}"
			}`,
		},
	},
	{
		name: "tool message with multiple results",
		input: []api.Message{
			&api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "call_123",
						ToolName:   "search",
						Result:     json.RawMessage(`{"temperature":"72°F","condition":"Sunny"}`),
					},
					{
						ToolCallID: "call_456",
						ToolName:   "calculator",
						Result:     json.RawMessage(`4`),
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"type": "function_call_output",
				"call_id": "call_123",
				"output": "{\"temperature\":\"72°F\",\"condition\":\"Sunny\"}"
			}`,
			`{
				"type": "function_call_output",
				"call_id": "call_456",
				"output": "4"
			}`,
		},
	},
	{
		name: "tool message with invalid JSON result",
		input: []api.Message{
			&api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "call_123",
						ToolName:   "search",
						Result:     json.RawMessage(`{invalid json`),
					},
				},
			},
		},
		expectedError: "failed to marshal tool result",
	},
	{
		name: "tool message with empty tool call ID",
		input: []api.Message{
			&api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "",
						ToolName:   "search",
						Result:     json.RawMessage(`{"data":"test"}`),
					},
				},
			},
		},
		expectedError: "tool result is missing tool call ID",
	},
	{
		name: "tool message with empty content",
		input: []api.Message{
			&api.ToolMessage{
				Content: []api.ToolResultBlock{},
			},
		},
		expectedMessages: []string{},
	},
	{
		name: "tool message with nil result",
		input: []api.Message{
			&api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "call_123",
						ToolName:   "search",
						Result:     nil,
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"type": "function_call_output",
				"call_id": "call_123",
				"output": "null"
			}`,
		},
	},
	{
		name: "computer tool result with safety checks",
		input: []api.Message{
			&api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "openai.computer_use_preview",
						Content: []api.ContentBlock{
							api.ImageBlock{
								Data:     []byte("test-image-data"),
								MimeType: "image/png",
							},
						},
						ProviderMetadata: api.NewProviderMetadata(map[string]any{
							"openai": &Metadata{
								ComputerSafetyChecks: []ComputerSafetyCheck{
									{
										ID:      "check1",
										Code:    "file_access",
										Message: "File access requested",
									},
									{
										ID:      "check2",
										Code:    "network_access",
										Message: "Network access requested",
									},
								},
							},
						}),
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"type": "computer_call_output",
				"call_id": "openai.computer_use_preview",
				"output": {
					"type": "computer_screenshot",
					"image_url": "data:image/png;base64,dGVzdC1pbWFnZS1kYXRh"
				},
				"acknowledged_safety_checks": [
					{
						"id": "check1",
						"code": "file_access",
						"message": "File access requested"
					},
					{
						"id": "check2",
						"code": "network_access",
						"message": "Network access requested"
					}
				]
			}`,
		},
	},
	{
		name: "computer tool result without safety checks",
		input: []api.Message{
			&api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "openai.computer_use_preview",
						Content: []api.ContentBlock{
							api.ImageBlock{
								Data:     []byte("test-image-data"),
								MimeType: "image/png",
							},
						},
					},
				},
			},
		},
		expectedMessages: []string{
			`{
				"type": "computer_call_output",
				"call_id": "openai.computer_use_preview",
				"output": {
					"type": "computer_screenshot",
					"image_url": "data:image/png;base64,dGVzdC1pbWFnZS1kYXRh"
				}
			}`,
		},
	},
	{
		name: "computer tool result with no content blocks",
		input: []api.Message{
			&api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "openai.computer_use_preview",
						Content:    []api.ContentBlock{},
					},
				},
			},
		},
		expectedError: "expected 1 content block for computer use tool result, got 0",
	},
	{
		name: "computer tool result with multiple content blocks",
		input: []api.Message{
			&api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "openai.computer_use_preview",
						Content: []api.ContentBlock{
							api.ImageBlock{
								Data:     []byte("test-image-data"),
								MimeType: "image/png",
							},
							api.ImageBlock{
								Data:     []byte("test-image-data-2"),
								MimeType: "image/png",
							},
						},
					},
				},
			},
		},
		expectedError: "expected 1 content block for computer use tool result, got 2",
	},
	{
		name: "computer tool result with wrong content type",
		input: []api.Message{
			&api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "openai.computer_use_preview",
						Content: []api.ContentBlock{
							&api.TextBlock{
								Text: "not an image",
							},
						},
					},
				},
			},
		},
		expectedError: "expected image block for computer use tool result",
	},
}

func TestEncodePrompt(t *testing.T) {
	t.Run("SystemMessages", func(t *testing.T) {
		runTestCases(t, systemMessageTests)
	})

	t.Run("UserMessages", func(t *testing.T) {
		runTestCases(t, userMessageTests)
	})

	t.Run("AssistantMessages", func(t *testing.T) {
		runTestCases(t, assistantMessageTests)
	})

	t.Run("ToolMessages", func(t *testing.T) {
		runTestCases(t, toolMessageTests)
	})
}
