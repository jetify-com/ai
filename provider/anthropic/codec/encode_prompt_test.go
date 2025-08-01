package codec

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/aitesting"
	"go.jetify.com/ai/api"
)

func TestEncodePrompt(t *testing.T) {
	tests := []struct {
		name    string
		prompt  []api.Message
		want    *AnthropicPrompt
		wantErr bool
	}{
		{
			name: "simple conversation with system message",
			prompt: []api.Message{
				&api.SystemMessage{Content: "You are a helpful assistant"},
				&api.UserMessage{Content: api.ContentFromText("Hello")},
				&api.AssistantMessage{Content: api.ContentFromText("Hi there!")},
			},
			want: &AnthropicPrompt{
				System: []anthropic.BetaTextBlockParam{
					NewTextBlock("You are a helpful assistant"),
				},
				Messages: []anthropic.BetaMessageParam{
					anthropic.NewBetaUserMessage(anthropic.BetaContentBlockParamUnion{
						OfText: &anthropic.BetaTextBlockParam{
							Type: "text",
							Text: "Hello",
						},
					}),
					NewAssistantMessage(anthropic.BetaContentBlockParamUnion{
						OfText: &anthropic.BetaTextBlockParam{
							Type: "text",
							Text: "Hi there!",
						},
					}),
				},
			},
		},
		{
			name: "multiple consecutive system messages",
			prompt: []api.Message{
				&api.SystemMessage{Content: "First system message"},
				&api.SystemMessage{Content: "Second system message"},
				&api.UserMessage{Content: api.ContentFromText("Hello")},
			},
			want: &AnthropicPrompt{
				System: []anthropic.BetaTextBlockParam{
					NewTextBlock("First system message\nSecond system message"),
				},
				Messages: []anthropic.BetaMessageParam{
					anthropic.NewBetaUserMessage(anthropic.BetaContentBlockParamUnion{
						OfText: &anthropic.BetaTextBlockParam{
							Type: "text",
							Text: "Hello",
						},
					}),
				},
			},
		},
		{
			name: "system message after non-system not allowed",
			prompt: []api.Message{
				&api.SystemMessage{Content: "First system message"},
				&api.UserMessage{Content: api.ContentFromText("Hello")},
				&api.SystemMessage{Content: "Second system message"},
			},
			wantErr: true,
		},
		{
			name: "conversation without system message",
			prompt: []api.Message{
				&api.UserMessage{Content: api.ContentFromText("Hello")},
				&api.AssistantMessage{Content: api.ContentFromText("Hi there!")},
			},
			want: &AnthropicPrompt{
				Messages: []anthropic.BetaMessageParam{
					anthropic.NewBetaUserMessage(anthropic.BetaContentBlockParamUnion{
						OfText: &anthropic.BetaTextBlockParam{
							Type: "text",
							Text: "Hello",
						},
					}),
					NewAssistantMessage(anthropic.BetaContentBlockParamUnion{
						OfText: &anthropic.BetaTextBlockParam{
							Type: "text",
							Text: "Hi there!",
						},
					}),
				},
			},
		},
		{
			name: "user message with multiple parts",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "hello"},
						api.ImageBlockFromData([]byte{0, 1, 2, 3}, "image/png"),
						api.FileBlockFromURL("http://example.com/file.txt"),
					},
				},
			},
			want: &AnthropicPrompt{
				Messages: []anthropic.BetaMessageParam{
					anthropic.NewBetaUserMessage(
						anthropic.BetaContentBlockParamUnion{
							OfText: &anthropic.BetaTextBlockParam{
								Type: "text",
								Text: "hello",
							},
						},
						anthropic.BetaContentBlockParamUnion{
							OfImage: &anthropic.BetaImageBlockParam{
								Type: "image",
								Source: anthropic.BetaImageBlockParamSourceUnion{
									OfBase64: &anthropic.BetaBase64ImageSourceParam{
										Type:      "base64",
										Data:      "AAECAw==",
										MediaType: "image/png",
									},
								},
							},
						},
						anthropic.BetaContentBlockParamUnion{
							OfDocument: &anthropic.BetaRequestDocumentBlockParam{
								Type: "document",
								Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
									OfText: &anthropic.BetaPlainTextSourceParam{
										Type:      "text",
										Data:      "http://example.com/file.txt",
										MediaType: "text/plain",
									},
								},
							},
						},
					),
				},
			},
		},
		{
			name: "with tool calls and results",
			prompt: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Let me check the weather"},
						&api.ToolCallBlock{
							ToolCallID: "weather-1",
							ToolName:   "get_weather",
							Args:       json.RawMessage(`{"location":"London"}`),
						},
					},
				},
				&api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID: "weather-1",
							ToolName:   "get_weather",
							Result:     json.RawMessage(`{"temperature": 20, "condition": "sunny"}`),
							IsError:    false,
						},
					},
				},
			},
			want: &AnthropicPrompt{
				Messages: []anthropic.BetaMessageParam{
					NewAssistantMessage(
						anthropic.BetaContentBlockParamUnion{
							OfText: &anthropic.BetaTextBlockParam{
								Type: "text",
								Text: "Let me check the weather",
							},
						},
						anthropic.BetaContentBlockParamUnion{
							OfToolUse: &anthropic.BetaToolUseBlockParam{
								Type:  "tool_use",
								ID:    "weather-1",
								Name:  "get_weather",
								Input: json.RawMessage(`{"location":"London"}`),
							},
						},
					),
					anthropic.NewBetaUserMessage(
						anthropic.BetaContentBlockParamUnion{
							OfToolResult: &anthropic.BetaToolResultBlockParam{
								Type:      "tool_result",
								ToolUseID: "weather-1",
								Content: []anthropic.BetaToolResultBlockParamContentUnion{
									{
										OfText: &anthropic.BetaTextBlockParam{
											Type: "text",
											Text: `{"temperature":20,"condition":"sunny"}`,
										},
									},
								},
								IsError: anthropic.Bool(false),
							},
						},
					),
				},
			},
		},
		{
			name: "user message with PDF file",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Here's a PDF"},
						&api.FileBlock{
							URL:       "http://example.com/doc.pdf",
							MediaType: "application/pdf",
						},
					},
				},
			},
			want: &AnthropicPrompt{
				Messages: []anthropic.BetaMessageParam{
					anthropic.NewBetaUserMessage(
						anthropic.BetaContentBlockParamUnion{
							OfText: &anthropic.BetaTextBlockParam{
								Type: "text",
								Text: "Here's a PDF",
							},
						},
						anthropic.BetaContentBlockParamUnion{
							OfDocument: &anthropic.BetaRequestDocumentBlockParam{
								Type: "document",
								Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
									OfURL: &anthropic.BetaURLPDFSourceParam{
										Type: "url",
										URL:  "http://example.com/doc.pdf",
									},
								},
							},
						},
					),
				},
				Betas: []anthropic.AnthropicBeta{"pdfs-2024-09-25"},
			},
		},
		{
			name: "user message with non-PDF file",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Here's a text file"},
						&api.FileBlock{
							URL:       "http://example.com/file.txt",
							MediaType: "text/plain",
						},
					},
				},
			},
			want: &AnthropicPrompt{
				Messages: []anthropic.BetaMessageParam{
					anthropic.NewBetaUserMessage(
						anthropic.BetaContentBlockParamUnion{
							OfText: &anthropic.BetaTextBlockParam{
								Type: "text",
								Text: "Here's a text file",
							},
						},
						anthropic.BetaContentBlockParamUnion{
							OfDocument: &anthropic.BetaRequestDocumentBlockParam{
								Type: "document",
								Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
									OfText: &anthropic.BetaPlainTextSourceParam{
										Type:      "text",
										Data:      "http://example.com/file.txt",
										MediaType: "text/plain",
									},
								},
							},
						},
					),
				},
			},
		},
		{
			name:   "empty prompt",
			prompt: []api.Message{},
			want: &AnthropicPrompt{
				Messages: []anthropic.BetaMessageParam{},
			},
		},
		{
			name:    "nil message in prompt",
			prompt:  []api.Message{nil},
			wantErr: true,
		},
		{
			name: "user message with encoding error",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						(*api.TextBlock)(nil), // Will cause encoding error
					},
				},
			},
			wantErr: true,
		},
		{
			name: "assistant message with encoding error",
			prompt: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						(*api.TextBlock)(nil), // Will cause encoding error
					},
				},
			},
			wantErr: true,
		},
		{
			name: "tool message with encoding error",
			prompt: []api.Message{
				&api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID: "", // Will cause encoding error
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "unsupported message type",
			prompt: []api.Message{
				&aitesting.MockUnsupportedMessage{},
			},
			wantErr: true,
		},
		{
			name: "tool message before system message",
			prompt: []api.Message{
				&api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID: "tool-1",
							ToolName:   "test_tool",
							Result:     json.RawMessage(`"simple string result"`),
							IsError:    false,
						},
					},
				},
				&api.SystemMessage{Content: "System message"},
			},
			want: &AnthropicPrompt{
				System: []anthropic.BetaTextBlockParam{
					NewTextBlock("System message"),
				},
				Messages: []anthropic.BetaMessageParam{
					anthropic.NewBetaUserMessage(
						anthropic.BetaContentBlockParamUnion{
							OfToolResult: &anthropic.BetaToolResultBlockParam{
								Type:      "tool_result",
								ToolUseID: "tool-1",
								Content: []anthropic.BetaToolResultBlockParamContentUnion{
									{
										OfText: &anthropic.BetaTextBlockParam{
											Type: "text",
											Text: `"simple string result"`,
										},
									},
								},
								IsError: anthropic.Bool(false),
							},
						},
					),
				},
			},
		},
		{
			name: "user message with value type (not pointer)",
			prompt: []api.Message{
				api.UserMessage{Content: api.ContentFromText("Hello from value type")},
			},
			want: &AnthropicPrompt{
				Messages: []anthropic.BetaMessageParam{
					anthropic.NewBetaUserMessage(anthropic.BetaContentBlockParamUnion{
						OfText: &anthropic.BetaTextBlockParam{
							Type: "text",
							Text: "Hello from value type",
						},
					}),
				},
			},
		},
		{
			name: "system message with value type (not pointer)",
			prompt: []api.Message{
				api.SystemMessage{Content: "System from value type"},
				api.UserMessage{Content: api.ContentFromText("Hello")},
			},
			want: &AnthropicPrompt{
				System: []anthropic.BetaTextBlockParam{
					NewTextBlock("System from value type"),
				},
				Messages: []anthropic.BetaMessageParam{
					anthropic.NewBetaUserMessage(anthropic.BetaContentBlockParamUnion{
						OfText: &anthropic.BetaTextBlockParam{
							Type: "text",
							Text: "Hello",
						},
					}),
				},
			},
		},
		{
			name: "assistant message with value type (not pointer)",
			prompt: []api.Message{
				api.AssistantMessage{Content: api.ContentFromText("Assistant from value type")},
			},
			want: &AnthropicPrompt{
				Messages: []anthropic.BetaMessageParam{
					NewAssistantMessage(anthropic.BetaContentBlockParamUnion{
						OfText: &anthropic.BetaTextBlockParam{
							Type: "text",
							Text: "Assistant from value type",
						},
					}),
				},
			},
		},
		{
			name: "tool message with value type (not pointer)",
			prompt: []api.Message{
				api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID: "tool-123",
							ToolName:   "test_tool",
							Result:     json.RawMessage(`"result from value type"`),
							IsError:    false,
						},
					},
				},
			},
			want: &AnthropicPrompt{
				Messages: []anthropic.BetaMessageParam{
					anthropic.NewBetaUserMessage(
						anthropic.BetaContentBlockParamUnion{
							OfToolResult: &anthropic.BetaToolResultBlockParam{
								Type:      "tool_result",
								ToolUseID: "tool-123",
								Content: []anthropic.BetaToolResultBlockParamContentUnion{
									{
										OfText: &anthropic.BetaTextBlockParam{
											Type: "text",
											Text: `"result from value type"`,
										},
									},
								},
								IsError: anthropic.Bool(false),
							},
						},
					),
				},
			},
		},
		{
			name: "user message with value type content blocks",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						api.TextBlock{
							Text: "Text as value",
						},
						api.ImageBlock{
							URL: "https://example.com/img.jpg",
						},
						api.FileBlock{
							Data:      []byte("test pdf data"),
							MediaType: "application/pdf",
						},
					},
				},
			},
			want: &AnthropicPrompt{
				Messages: []anthropic.BetaMessageParam{
					anthropic.NewBetaUserMessage(
						anthropic.BetaContentBlockParamUnion{
							OfText: &anthropic.BetaTextBlockParam{
								Type: "text",
								Text: "Text as value",
							},
						},
						anthropic.BetaContentBlockParamUnion{
							OfImage: &anthropic.BetaImageBlockParam{
								Type: "image",
								Source: anthropic.BetaImageBlockParamSourceUnion{
									OfURL: &anthropic.BetaURLImageSourceParam{
										Type: "url",
										URL:  "https://example.com/img.jpg",
									},
								},
							},
						},
						anthropic.BetaContentBlockParamUnion{
							OfDocument: &anthropic.BetaRequestDocumentBlockParam{
								Type: "document",
								Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
									OfBase64: &anthropic.BetaBase64PDFSourceParam{
										Type:      "base64",
										Data:      base64.StdEncoding.EncodeToString([]byte("test pdf data")),
										MediaType: "application/pdf",
									},
								},
							},
						},
					),
				},
				Betas: []anthropic.AnthropicBeta{"pdfs-2024-09-25"},
			},
		},
		{
			name: "assistant message with value type content blocks",
			prompt: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						api.TextBlock{
							Text: "Assistant text as value",
						},
						api.ToolCallBlock{
							ToolCallID: "tool-456",
							ToolName:   "search_tool",
							Args:       json.RawMessage(`{"query":"test"}`),
						},
					},
				},
			},
			want: &AnthropicPrompt{
				Messages: []anthropic.BetaMessageParam{
					NewAssistantMessage(
						anthropic.BetaContentBlockParamUnion{
							OfText: &anthropic.BetaTextBlockParam{
								Type: "text",
								Text: "Assistant text as value",
							},
						},
						anthropic.BetaContentBlockParamUnion{
							OfToolUse: &anthropic.BetaToolUseBlockParam{
								Type:  "tool_use",
								ID:    "tool-456",
								Name:  "search_tool",
								Input: json.RawMessage(`{"query":"test"}`),
							},
						},
					),
				},
			},
		},
		{
			name: "tool message with value type content blocks in result",
			prompt: []api.Message{
				&api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID: "tool-789",
							ToolName:   "image_tool",
							Content: []api.ContentBlock{
								api.TextBlock{
									Text: "Tool result text",
								},
								api.ImageBlock{
									Data:      []byte{1, 2, 3, 4},
									MediaType: "image/png",
								},
							},
						},
					},
				},
			},
			want: &AnthropicPrompt{
				Messages: []anthropic.BetaMessageParam{
					anthropic.NewBetaUserMessage(
						anthropic.BetaContentBlockParamUnion{
							OfToolResult: &anthropic.BetaToolResultBlockParam{
								Type:      "tool_result",
								ToolUseID: "tool-789",
								Content: []anthropic.BetaToolResultBlockParamContentUnion{
									{
										OfText: &anthropic.BetaTextBlockParam{
											Type: "text",
											Text: "Tool result text",
										},
									},
									{
										OfImage: &anthropic.BetaImageBlockParam{
											Type: "image",
											Source: anthropic.BetaImageBlockParamSourceUnion{
												OfBase64: &anthropic.BetaBase64ImageSourceParam{
													Type:      "base64",
													Data:      base64.StdEncoding.EncodeToString([]byte{1, 2, 3, 4}),
													MediaType: "image/png",
												},
											},
										},
									},
								},
								IsError: anthropic.Bool(false),
							},
						},
					),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodePrompt(tt.prompt)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare JSON representations
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), string(gotJSON))
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
			name: "assistant message with image block",
			prompt: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						api.ImageBlockFromURL("http://example.com/image.jpg"),
					},
				},
			},
			expectedError: "unsupported assistant content block type",
		},
		{
			name: "image block without URL or data",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.ImageBlock{},
					},
				},
			},
			expectedError: "image block must have either URL or Data",
		},
		{
			name: "tool message with error",
			prompt: []api.Message{
				&api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID: "", // Invalid empty ID
							ToolName:   "test_tool",
							Result:     json.RawMessage(`"result"`),
							IsError:    false,
						},
					},
				},
			},
			expectedError: "tool call ID cannot be empty",
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

func TestEncodeUserContentPart(t *testing.T) {
	tests := []struct {
		name     string
		block    api.ContentBlock
		want     anthropic.BetaContentBlockParamUnion
		wantBeta []anthropic.AnthropicBeta
		wantErr  bool
	}{
		{
			name:  "text block",
			block: &api.TextBlock{Text: "Hello, world!"},
			want: anthropic.BetaContentBlockParamUnion{
				OfText: &anthropic.BetaTextBlockParam{
					Type: "text",
					Text: "Hello, world!",
				},
			},
		},
		{
			name:  "image block with URL",
			block: &api.ImageBlock{URL: "https://example.com/image.jpg"},
			want: anthropic.BetaContentBlockParamUnion{
				OfImage: &anthropic.BetaImageBlockParam{
					Type: "image",
					Source: anthropic.BetaImageBlockParamSourceUnion{
						OfURL: &anthropic.BetaURLImageSourceParam{
							Type: "url",
							URL:  "https://example.com/image.jpg",
						},
					},
				},
			},
		},
		{
			name: "image block with data",
			block: &api.ImageBlock{
				Data:      []byte("fake-image-data"),
				MediaType: "image/jpeg",
			},
			want: anthropic.BetaContentBlockParamUnion{
				OfImage: &anthropic.BetaImageBlockParam{
					Type: "image",
					Source: anthropic.BetaImageBlockParamSourceUnion{
						OfBase64: &anthropic.BetaBase64ImageSourceParam{
							Type:      "base64",
							Data:      "ZmFrZS1pbWFnZS1kYXRh",
							MediaType: "image/jpeg",
						},
					},
				},
			},
		},
		{
			name: "image block with data and missing mime type",
			block: &api.ImageBlock{
				Data: []byte("fake-image-data"),
			},
			want: anthropic.BetaContentBlockParamUnion{
				OfImage: &anthropic.BetaImageBlockParam{
					Type: "image",
					Source: anthropic.BetaImageBlockParamSourceUnion{
						OfBase64: &anthropic.BetaBase64ImageSourceParam{
							Type:      "base64",
							Data:      "ZmFrZS1pbWFnZS1kYXRh",
							MediaType: "image/jpeg",
						},
					},
				},
			},
		},
		{
			name:  "file block with URL",
			block: &api.FileBlock{URL: "https://example.com/file.txt"},
			want: anthropic.BetaContentBlockParamUnion{
				OfDocument: &anthropic.BetaRequestDocumentBlockParam{
					Type: "document",
					Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
						OfText: &anthropic.BetaPlainTextSourceParam{
							Type:      "text",
							Data:      "https://example.com/file.txt",
							MediaType: "text/plain",
						},
					},
				},
			},
		},
		{
			name:  "file block with text data",
			block: &api.FileBlock{Data: []byte("Hello from file")},
			want: anthropic.BetaContentBlockParamUnion{
				OfDocument: &anthropic.BetaRequestDocumentBlockParam{
					Type: "document",
					Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
						OfText: &anthropic.BetaPlainTextSourceParam{
							Type:      "text",
							Data:      "Hello from file",
							MediaType: "text/plain",
						},
					},
				},
			},
		},
		{
			name: "file block with binary data and mime type",
			block: &api.FileBlock{
				Data:      []byte{0, 1, 2, 3},
				MediaType: "audio/wav",
			},
			wantErr: true,
		},
		{
			name:    "empty file block",
			block:   &api.FileBlock{},
			wantErr: true,
		},
		{
			name:    "invalid block type",
			block:   &aitesting.MockUnsupportedBlock{},
			wantErr: true,
		},
		{
			name:    "nil block",
			block:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotBetas, err := EncodeUserContentBlock(tt.block)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare JSON representations
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), string(gotJSON))

			assert.Equal(t, tt.wantBeta, gotBetas)
		})
	}
}

func TestEncodeAssistantMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     *api.AssistantMessage
		want    anthropic.BetaMessageParam
		wantErr bool
	}{
		{
			name: "text only",
			msg: &api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Hello!"},
				},
			},
			want: NewAssistantMessage(
				anthropic.BetaContentBlockParamUnion{
					OfText: &anthropic.BetaTextBlockParam{
						Type: "text",
						Text: "Hello!",
					},
				},
			),
		},
		{
			name: "with tool call",
			msg: &api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Let me check"},
					&api.ToolCallBlock{
						ToolCallID: "tool-1",
						ToolName:   "test_tool",
						Args:       json.RawMessage(`{"key":"value"}`),
					},
				},
			},
			want: NewAssistantMessage(
				anthropic.BetaContentBlockParamUnion{
					OfText: &anthropic.BetaTextBlockParam{
						Type: "text",
						Text: "Let me check",
					},
				},
				anthropic.BetaContentBlockParamUnion{
					OfToolUse: &anthropic.BetaToolUseBlockParam{
						Type:  "tool_use",
						ID:    "tool-1",
						Name:  "test_tool",
						Input: json.RawMessage(`{"key":"value"}`),
					},
				},
			),
		},
		{
			name: "with reasoning",
			msg: &api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.ReasoningBlock{
						Text:      "Let me think about this...",
						Signature: "sig123",
					},
				},
			},
			want: NewAssistantMessage(
				anthropic.BetaContentBlockParamUnion{
					OfThinking: &anthropic.BetaThinkingBlockParam{
						Type:      "thinking",
						Thinking:  "Let me think about this...",
						Signature: "sig123",
					},
				},
			),
		},
		{
			name: "with redacted reasoning",
			msg: &api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.RedactedReasoningBlock{
						Data: "redacted-data",
					},
				},
			},
			want: NewAssistantMessage(
				anthropic.BetaContentBlockParamUnion{
					OfRedactedThinking: &anthropic.BetaRedactedThinkingBlockParam{
						Type: "redacted_thinking",
						Data: "redacted-data",
					},
				},
			),
		},
		{
			name: "invalid content type",
			msg: &api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.ImageBlock{}, // Images not supported in assistant messages
				},
			},
			wantErr: true,
		},
		{
			name: "with unsupported part",
			msg: &api.AssistantMessage{
				Content: []api.ContentBlock{
					&aitesting.MockUnsupportedBlock{},
				},
			},
			wantErr: true,
		},
		{
			name: "nil text block",
			msg: &api.AssistantMessage{
				Content: []api.ContentBlock{
					(*api.TextBlock)(nil),
				},
			},
			wantErr: true,
		},
		{
			name: "nil tool call block",
			msg: &api.AssistantMessage{
				Content: []api.ContentBlock{
					(*api.ToolCallBlock)(nil),
				},
			},
			wantErr: true,
		},
		{
			name: "nil reasoning block",
			msg: &api.AssistantMessage{
				Content: []api.ContentBlock{
					(*api.ReasoningBlock)(nil),
				},
			},
			wantErr: true,
		},
		{
			name: "nil redacted reasoning block",
			msg: &api.AssistantMessage{
				Content: []api.ContentBlock{
					(*api.RedactedReasoningBlock)(nil),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeAssistantMessage(tt.msg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare JSON representations
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}

func TestEncodeToolMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     *api.ToolMessage
		want    anthropic.BetaMessageParam
		wantErr bool
	}{
		{
			name: "simple result",
			msg: &api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "tool-1",
						ToolName:   "test_tool",
						Result:     json.RawMessage(`"success"`),
						IsError:    false,
					},
				},
			},
			want: anthropic.NewBetaUserMessage(
				anthropic.BetaContentBlockParamUnion{
					OfToolResult: &anthropic.BetaToolResultBlockParam{
						Type:      "tool_result",
						ToolUseID: "tool-1",
						Content: []anthropic.BetaToolResultBlockParamContentUnion{
							{
								OfText: &anthropic.BetaTextBlockParam{
									Type: "text",
									Text: `"success"`,
								},
							},
						},
						IsError: anthropic.Bool(false),
					},
				},
			),
		},
		{
			name: "error result",
			msg: &api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "tool-1",
						ToolName:   "test_tool",
						Result:     json.RawMessage(`"failed"`),
						IsError:    true,
					},
				},
			},
			want: anthropic.NewBetaUserMessage(
				anthropic.BetaContentBlockParamUnion{
					OfToolResult: &anthropic.BetaToolResultBlockParam{
						Type:      "tool_result",
						ToolUseID: "tool-1",
						Content: []anthropic.BetaToolResultBlockParamContentUnion{
							{
								OfText: &anthropic.BetaTextBlockParam{
									Type: "text",
									Text: `"failed"`,
								},
							},
						},
						IsError: anthropic.Bool(true),
					},
				},
			),
		},
		{
			name: "structured result",
			msg: &api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "tool-1",
						ToolName:   "test_tool",
						Result:     json.RawMessage(`{"status":"ok","data":123}`),
						IsError:    false,
					},
				},
			},
			want: anthropic.NewBetaUserMessage(
				anthropic.BetaContentBlockParamUnion{
					OfToolResult: &anthropic.BetaToolResultBlockParam{
						Type:      "tool_result",
						ToolUseID: "tool-1",
						Content: []anthropic.BetaToolResultBlockParamContentUnion{
							{
								OfText: &anthropic.BetaTextBlockParam{
									Type: "text",
									Text: `{"status":"ok","data":123}`,
								},
							},
						},
						IsError: anthropic.Bool(false),
					},
				},
			),
		},
		{
			name: "tool message with rich content",
			msg: &api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "tool-1",
						ToolName:   "image_generator",
						Content: []api.ContentBlock{
							&api.TextBlock{
								Text: "Generated image:",
							},
							&api.ImageBlock{
								Data:      []byte("base64data"),
								MediaType: "image/png",
							},
						},
					},
				},
			},
			want: anthropic.NewBetaUserMessage(
				anthropic.BetaContentBlockParamUnion{
					OfToolResult: &anthropic.BetaToolResultBlockParam{
						Type:      "tool_result",
						ToolUseID: "tool-1",
						Content: []anthropic.BetaToolResultBlockParamContentUnion{
							{
								OfText: &anthropic.BetaTextBlockParam{
									Type: "text",
									Text: "Generated image:",
								},
							},
							{
								OfImage: &anthropic.BetaImageBlockParam{
									Type: "image",
									Source: anthropic.BetaImageBlockParamSourceUnion{
										OfBase64: &anthropic.BetaBase64ImageSourceParam{
											Type:      "base64",
											Data:      "YmFzZTY0ZGF0YQ==",
											MediaType: "image/png",
										},
									},
								},
							},
						},
						IsError: anthropic.Bool(false),
					},
				},
			),
		},
		{
			name: "tool message with invalid content type",
			msg: &api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "tool-1",
						Content: []api.ContentBlock{
							&aitesting.MockUnsupportedBlock{},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "with cache control",
			msg: &api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "tool-1",
						ToolName:   "test_tool",
						Result:     json.RawMessage(`"success"`),
						ProviderMetadata: api.NewProviderMetadata(map[string]any{
							"anthropic": Metadata{
								CacheControl: "ephemeral",
							},
						}),
					},
				},
			},
			want: anthropic.NewBetaUserMessage(
				anthropic.BetaContentBlockParamUnion{
					OfToolResult: &anthropic.BetaToolResultBlockParam{
						Type:      "tool_result",
						ToolUseID: "tool-1",
						Content: []anthropic.BetaToolResultBlockParamContentUnion{
							{
								OfText: &anthropic.BetaTextBlockParam{
									Type: "text",
									Text: `"success"`,
								},
							},
						},
						IsError: anthropic.Bool(false),
						CacheControl: anthropic.BetaCacheControlEphemeralParam{
							Type: "ephemeral",
						},
					},
				},
			),
		},
		{
			name: "empty tool call ID",
			msg: &api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "", // Empty ID should cause error
						ToolName:   "test_tool",
						Result:     json.RawMessage(`"success"`),
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeToolMessage(tt.msg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare JSON representations
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}

func TestEncodeSystemMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     *api.SystemMessage
		want    anthropic.BetaTextBlockParam
		wantErr bool
	}{
		{
			name: "simple system message",
			msg:  &api.SystemMessage{Content: "You are a helpful assistant"},
			want: NewTextBlock("You are a helpful assistant"),
		},
		{
			name: "system message with cache control",
			msg: &api.SystemMessage{
				Content: "You are a helpful assistant",
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"anthropic": Metadata{
						CacheControl: "ephemeral",
					},
				}),
			},
			want: anthropic.BetaTextBlockParam{
				Type: "text",
				Text: "You are a helpful assistant",
				CacheControl: anthropic.BetaCacheControlEphemeralParam{
					Type: "ephemeral",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeSystemMessage(tt.msg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare JSON representations
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}

func TestEncodeTextPart(t *testing.T) {
	tests := []struct {
		name    string
		block   *api.TextBlock
		want    anthropic.BetaTextBlockParam
		wantErr bool
	}{
		{
			name:  "simple text",
			block: &api.TextBlock{Text: "Hello, world!"},
			want: anthropic.BetaTextBlockParam{
				Type: "text",
				Text: "Hello, world!",
			},
		},
		{
			name:  "empty text",
			block: &api.TextBlock{Text: ""},
			want: anthropic.BetaTextBlockParam{
				Type: "text",
				Text: "",
			},
		},
		{
			name: "text with cache control",
			block: &api.TextBlock{
				Text: "Hello",
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"anthropic": Metadata{
						CacheControl: "ephemeral",
					},
				}),
			},
			want: anthropic.BetaTextBlockParam{
				Type: "text",
				Text: "Hello",
				CacheControl: anthropic.BetaCacheControlEphemeralParam{
					Type: "ephemeral",
				},
			},
		},
		{
			name:    "nil block",
			block:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeTextBlock(tt.block)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare JSON representations
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}

func TestEncodeToolCallPart(t *testing.T) {
	tests := []struct {
		name    string
		block   *api.ToolCallBlock
		want    anthropic.BetaToolUseBlockParam
		wantErr bool
	}{
		{
			name: "simple tool call",
			block: &api.ToolCallBlock{
				ToolCallID: "tool-1",
				ToolName:   "test_tool",
				Args:       json.RawMessage(`{"key":"value"}`),
			},
			want: anthropic.BetaToolUseBlockParam{
				Type:  "tool_use",
				ID:    "tool-1",
				Name:  "test_tool",
				Input: json.RawMessage(`{"key":"value"}`),
			},
		},
		{
			name: "empty tool name",
			block: &api.ToolCallBlock{
				ToolCallID: "tool-1",
				ToolName:   "",
				Args:       json.RawMessage(`{}`),
			},
			want: anthropic.BetaToolUseBlockParam{
				Type:  "tool_use",
				ID:    "tool-1",
				Name:  "",
				Input: json.RawMessage(`{}`),
			},
		},
		{
			name: "with cache control",
			block: &api.ToolCallBlock{
				ToolCallID: "tool-1",
				ToolName:   "test_tool",
				Args:       json.RawMessage(`{"key":"value"}`),
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"anthropic": Metadata{
						CacheControl: "ephemeral",
					},
				}),
			},
			want: anthropic.BetaToolUseBlockParam{
				Type:  "tool_use",
				ID:    "tool-1",
				Name:  "test_tool",
				Input: json.RawMessage(`{"key":"value"}`),
				CacheControl: anthropic.BetaCacheControlEphemeralParam{
					Type: "ephemeral",
				},
			},
		},
		{
			name:    "nil block",
			block:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeToolCallBlock(tt.block)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare JSON representations
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}

func TestEncodeFilePart(t *testing.T) {
	tests := []struct {
		name     string
		block    *api.FileBlock
		want     anthropic.BetaRequestDocumentBlockParam
		wantBeta []anthropic.AnthropicBeta
		wantErr  bool
	}{
		{
			name:  "file with URL",
			block: &api.FileBlock{URL: "http://example.com/file.txt"},
			want: anthropic.BetaRequestDocumentBlockParam{
				Type: "document",
				Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
					OfText: &anthropic.BetaPlainTextSourceParam{
						Type:      "text",
						Data:      "http://example.com/file.txt",
						MediaType: "text/plain",
					},
				},
			},
		},
		{
			name:  "PDF file with URL",
			block: &api.FileBlock{URL: "http://example.com/doc.pdf"},
			want: anthropic.BetaRequestDocumentBlockParam{
				Type: "document",
				Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
					OfURL: &anthropic.BetaURLPDFSourceParam{
						Type: "url",
						URL:  "http://example.com/doc.pdf",
					},
				},
			},
			wantBeta: []anthropic.AnthropicBeta{anthropic.AnthropicBetaPDFs2024_09_25},
		},
		{
			name:  "PDF file with URL and query parameters",
			block: &api.FileBlock{URL: "http://example.com/doc.pdf?version=2&user=test"},
			want: anthropic.BetaRequestDocumentBlockParam{
				Type: "document",
				Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
					OfURL: &anthropic.BetaURLPDFSourceParam{
						Type: "url",
						URL:  "http://example.com/doc.pdf?version=2&user=test",
					},
				},
			},
			wantBeta: []anthropic.AnthropicBeta{anthropic.AnthropicBetaPDFs2024_09_25},
		},
		{
			name:  "PDF file with URL and fragment",
			block: &api.FileBlock{URL: "http://example.com/doc.pdf#page=5"},
			want: anthropic.BetaRequestDocumentBlockParam{
				Type: "document",
				Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
					OfURL: &anthropic.BetaURLPDFSourceParam{
						Type: "url",
						URL:  "http://example.com/doc.pdf#page=5",
					},
				},
			},
			wantBeta: []anthropic.AnthropicBeta{anthropic.AnthropicBetaPDFs2024_09_25},
		},
		{
			name:  "PDF file with complex path",
			block: &api.FileBlock{URL: "http://example.com/files/documents/2023/report.pdf"},
			want: anthropic.BetaRequestDocumentBlockParam{
				Type: "document",
				Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
					OfURL: &anthropic.BetaURLPDFSourceParam{
						Type: "url",
						URL:  "http://example.com/files/documents/2023/report.pdf",
					},
				},
			},
			wantBeta: []anthropic.AnthropicBeta{anthropic.AnthropicBetaPDFs2024_09_25},
		},
		{
			name: "PDF file with data",
			block: &api.FileBlock{
				Data:      []byte("PDF data"),
				MediaType: "application/pdf",
			},
			want: anthropic.BetaRequestDocumentBlockParam{
				Type: "document",
				Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
					OfBase64: &anthropic.BetaBase64PDFSourceParam{
						Type:      "base64",
						Data:      base64.StdEncoding.EncodeToString([]byte("PDF data")),
						MediaType: "application/pdf",
					},
				},
			},
			wantBeta: []anthropic.AnthropicBeta{anthropic.AnthropicBetaPDFs2024_09_25},
		},
		{
			name:  "file with text data",
			block: &api.FileBlock{Data: []byte("Hello from file")},
			want: anthropic.BetaRequestDocumentBlockParam{
				Type: "document",
				Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
					OfText: &anthropic.BetaPlainTextSourceParam{
						Type:      "text",
						Data:      "Hello from file",
						MediaType: "text/plain",
					},
				},
			},
		},
		{
			name: "file with binary data and mime type",
			block: &api.FileBlock{
				Data:      []byte{0, 1, 2, 3},
				MediaType: "audio/wav",
			},
			wantErr: true,
		},
		{
			name:    "empty file part",
			block:   &api.FileBlock{},
			wantErr: true,
		},
		{
			name: "with cache control",
			block: &api.FileBlock{
				URL: "http://example.com/file.txt",
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"anthropic": Metadata{
						CacheControl: "ephemeral",
					},
				}),
			},
			want: anthropic.BetaRequestDocumentBlockParam{
				Type: "document",
				Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
					OfText: &anthropic.BetaPlainTextSourceParam{
						Type:      "text",
						Data:      "http://example.com/file.txt",
						MediaType: "text/plain",
					},
				},
				CacheControl: anthropic.BetaCacheControlEphemeralParam{
					Type: "ephemeral",
				},
			},
		},
		{
			name:    "nil block",
			block:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotBetas, err := EncodeFileBlock(tt.block)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare JSON representations
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), string(gotJSON))

			assert.Equal(t, tt.wantBeta, gotBetas)
		})
	}
}

func TestGetCacheControl(t *testing.T) {
	tests := []struct {
		name   string
		source api.MetadataSource
		want   *anthropic.BetaCacheControlEphemeralParam
	}{
		{
			name:   "nil metadata",
			source: nil,
			want:   nil,
		},
		{
			name:   "empty metadata",
			source: &aitesting.MockMetadataSource{ProviderMetadata: api.NewProviderMetadata(nil)},
			want:   nil,
		},
		{
			name: "no anthropic metadata",
			source: &aitesting.MockMetadataSource{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"other": map[string]any{},
				}),
			},
			want: nil,
		},
		{
			name: "no cache control in anthropic metadata",
			source: &aitesting.MockMetadataSource{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"anthropic": Metadata{},
				}),
			},
			want: nil,
		},
		{
			name: "with cache control",
			source: &aitesting.MockMetadataSource{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"anthropic": Metadata{
						CacheControl: "ephemeral",
					},
				}),
			},
			want: &anthropic.BetaCacheControlEphemeralParam{
				Type: "ephemeral",
			},
		},
		{
			name: "with invalid cache control type",
			source: &aitesting.MockMetadataSource{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"anthropic": Metadata{
						CacheControl: "something-else",
					},
				}),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getCacheControl(tt.source)
			if tt.want == nil {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)

			// Compare JSON representations
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}

func TestEncodeImagePart(t *testing.T) {
	tests := []struct {
		name    string
		block   *api.ImageBlock
		want    anthropic.BetaImageBlockParam
		wantErr bool
	}{
		{
			name:  "image with URL",
			block: &api.ImageBlock{URL: "http://example.com/image.jpg"},
			want: anthropic.BetaImageBlockParam{
				Type: "image",
				Source: anthropic.BetaImageBlockParamSourceUnion{
					OfURL: &anthropic.BetaURLImageSourceParam{
						Type: "url",
						URL:  "http://example.com/image.jpg",
					},
				},
			},
		},
		{
			name: "image with data and mime type",
			block: &api.ImageBlock{
				Data:      []byte{0, 1, 2, 3},
				MediaType: "image/png",
			},
			want: NewImageBlockBase64("image/png", "AAECAw=="),
		},
		{
			name: "image with data and no mime type",
			block: &api.ImageBlock{
				Data: []byte{0, 1, 2, 3},
			},
			want: NewImageBlockBase64("image/jpeg", "AAECAw=="),
		},
		{
			name: "with cache control",
			block: &api.ImageBlock{
				URL: "http://example.com/image.jpg",
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"anthropic": Metadata{
						CacheControl: "ephemeral",
					},
				}),
			},
			want: anthropic.BetaImageBlockParam{
				Type: "image",
				Source: anthropic.BetaImageBlockParamSourceUnion{
					OfURL: &anthropic.BetaURLImageSourceParam{
						Type: "url",
						URL:  "http://example.com/image.jpg",
					},
				},
				CacheControl: anthropic.BetaCacheControlEphemeralParam{
					Type: "ephemeral",
				},
			},
		},
		{
			name:    "empty image block",
			block:   &api.ImageBlock{},
			wantErr: true,
		},
		{
			name:    "nil block",
			block:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeImageBlock(tt.block)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare JSON representations
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}

func TestEncodeReasoningBlock(t *testing.T) {
	tests := []struct {
		name    string
		block   *api.ReasoningBlock
		want    anthropic.BetaThinkingBlockParam
		wantErr bool
	}{
		{
			name: "valid reasoning block",
			block: &api.ReasoningBlock{
				Text:      "Let me think about this...",
				Signature: "sig123",
			},
			want: anthropic.BetaThinkingBlockParam{
				Type:      "thinking",
				Thinking:  "Let me think about this...",
				Signature: "sig123",
			},
		},
		{
			name:    "nil block",
			block:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeReasoningBlock(tt.block)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare JSON representations
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}

func TestEncodeRedactedReasoningBlock(t *testing.T) {
	tests := []struct {
		name    string
		block   *api.RedactedReasoningBlock
		want    anthropic.BetaRedactedThinkingBlockParam
		wantErr bool
	}{
		{
			name: "valid redacted reasoning block",
			block: &api.RedactedReasoningBlock{
				Data: "redacted-data",
			},
			want: anthropic.BetaRedactedThinkingBlockParam{
				Type: "redacted_thinking",
				Data: "redacted-data",
			},
		},
		{
			name:    "nil block",
			block:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeRedactedReasoningBlock(tt.block)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare JSON representations
			wantJSON, err := json.Marshal(tt.want)
			require.NoError(t, err)
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)
			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}
