package codec

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/aitesting"
	"go.jetify.com/ai/api"
)

func TestMergeMessageGroup(t *testing.T) {
	tests := []struct {
		name     string
		messages []api.Message
		want     api.Message
		wantErr  bool
	}{
		{
			name:     "empty_messages",
			messages: nil,
			want:     nil,
			wantErr:  false,
		},
		{
			name: "unsupported_message_type",
			messages: []api.Message{
				&aitesting.MockUnsupportedMessage{},
				&aitesting.MockUnsupportedMessage{},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeMessageGroup(tt.messages)
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeMessageGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMergeMessages(t *testing.T) {
	tests := []struct {
		name     string
		prompt   []api.Message
		expected []api.Message
	}{
		{
			name:     "empty_prompt",
			prompt:   []api.Message{},
			expected: []api.Message{},
		},
		{
			name: "single_message",
			prompt: []api.Message{
				&api.SystemMessage{Content: "Hello"},
			},
			expected: []api.Message{
				&api.SystemMessage{Content: "Hello"},
			},
		},
		{
			name: "consecutive_system_messages",
			prompt: []api.Message{
				&api.SystemMessage{Content: "First"},
				&api.SystemMessage{Content: "Second"},
			},
			expected: []api.Message{
				&api.SystemMessage{Content: "First\nSecond"},
			},
		},
		{
			name: "mixed_message_types",
			prompt: []api.Message{
				&api.SystemMessage{Content: "System"},
				&api.UserMessage{Content: api.ContentFromText("User")},
				&api.AssistantMessage{Content: api.ContentFromText("Assistant")},
			},
			expected: []api.Message{
				&api.SystemMessage{Content: "System"},
				&api.UserMessage{Content: api.ContentFromText("User")},
				&api.AssistantMessage{Content: api.ContentFromText("Assistant")},
			},
		},
		{
			name: "preserve_metadata_in_system_messages",
			prompt: []api.Message{
				&api.SystemMessage{
					Content:          "First",
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"key1": map[string]any{"value": "1"}}),
				},
				&api.SystemMessage{
					Content:          "Second",
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"key2": map[string]any{"value": "2"}}),
				},
			},
			expected: []api.Message{
				&api.SystemMessage{
					Content:          "First\nSecond",
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"key2": map[string]any{"value": "2"}}),
				},
			},
		},
		{
			name: "preserve_block_metadata_in_user_messages",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{
							Text:             "First",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"block1": map[string]any{"value": "1"}}),
						},
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"msg1": map[string]any{"value": "1"}}),
				},
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{
							Text:             "Second",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"block2": map[string]any{"value": "2"}}),
						},
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"msg2": map[string]any{"value": "2"}}),
				},
			},
			expected: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{
							Text:             "First",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"block1": map[string]any{"value": "1"}}),
						},
						&api.TextBlock{
							Text:             "Second",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"block2": map[string]any{"value": "2"}}),
						},
					},
				},
			},
		},
		{
			name: "preserve_message_metadata_in_last_block_if_no_block_metadata",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "First"},
					},
				},
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Second"},
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"msg": map[string]any{"value": "metadata"}}),
				},
			},
			expected: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "First"},
						&api.TextBlock{
							Text:             "Second",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"msg": map[string]any{"value": "metadata"}}),
						},
					},
				},
			},
		},
		{
			name: "empty_message_group",
			prompt: []api.Message{
				nil,
			},
			expected: []api.Message{
				nil,
			},
		},
		{
			name: "unsupported_message_type",
			prompt: []api.Message{
				&aitesting.MockUnsupportedMessage{},
				&aitesting.MockUnsupportedMessage{},
			},
			expected: []api.Message{
				&aitesting.MockUnsupportedMessage{},
				&aitesting.MockUnsupportedMessage{},
			},
		},
		{
			name: "combine_with_image_and_file_blocks",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.ImageBlock{
							URL:              "image1.jpg",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"img1": map[string]any{"value": "1"}}),
						},
						&api.FileBlock{
							URL:              "file1.txt",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"file1": map[string]any{"value": "1"}}),
						},
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"msg1": map[string]any{"value": "1"}}),
				},
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.ImageBlock{URL: "image2.jpg"},
						&api.FileBlock{URL: "file2.txt"},
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"msg2": map[string]any{"value": "2"}}),
				},
			},
			expected: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.ImageBlock{
							URL:              "image1.jpg",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"img1": map[string]any{"value": "1"}}),
						},
						&api.FileBlock{
							URL:              "file1.txt",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"file1": map[string]any{"value": "1"}}),
						},
						&api.ImageBlock{URL: "image2.jpg"},
						&api.FileBlock{
							URL:              "file2.txt",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"msg2": map[string]any{"value": "2"}}),
						},
					},
				},
			},
		},
		{
			name: "combine_with_reasoning_blocks",
			prompt: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.ReasoningBlock{
							Text:      "reasoning1",
							Signature: "sig1",
						},
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"msg1": map[string]any{"value": "1"}}),
				},
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.RedactedReasoningBlock{
							Data: "redacted1",
						},
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"msg2": map[string]any{"value": "2"}}),
				},
			},
			expected: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.ReasoningBlock{
							Text:      "reasoning1",
							Signature: "sig1",
						},
						&api.RedactedReasoningBlock{
							Data:             "redacted1",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"msg2": map[string]any{"value": "2"}}),
						},
					},
				},
			},
		},
		{
			name: "nil_messages_in_combine",
			prompt: []api.Message{
				&api.SystemMessage{Content: "Hello"},
				&api.SystemMessage{Content: "World"},
			},
			expected: []api.Message{
				&api.SystemMessage{Content: "Hello\nWorld"},
			},
		},
		{
			name: "empty_content_blocks",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{},
				},
				&api.UserMessage{
					Content:          []api.ContentBlock{},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"msg": map[string]any{"value": "metadata"}}),
				},
			},
			expected: []api.Message{
				&api.UserMessage{
					Content: nil,
				},
			},
		},
		{
			name: "nil_content_blocks",
			prompt: []api.Message{
				&api.AssistantMessage{
					Content: nil,
				},
				&api.AssistantMessage{
					Content:          nil,
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"msg": map[string]any{"value": "metadata"}}),
				},
			},
			expected: []api.Message{
				&api.AssistantMessage{
					Content: nil,
				},
			},
		},
		{
			name: "mixed_content_blocks_no_combine",
			prompt: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Text"},
						&api.ImageBlock{URL: "image.jpg"},
					},
				},
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Response"},
						&api.ReasoningBlock{Text: "Reasoning", Signature: "sig"},
					},
				},
				&api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID: "tool1",
							ToolName:   "test-tool",
							Result:     json.RawMessage(`"Tool result"`),
						},
					},
				},
			},
			expected: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Text"},
						&api.ImageBlock{URL: "image.jpg"},
					},
				},
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Response"},
						&api.ReasoningBlock{Text: "Reasoning", Signature: "sig"},
					},
				},
				&api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID: "tool1",
							ToolName:   "test-tool",
							Result:     json.RawMessage(`"Tool result"`),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeMessages(tt.prompt)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMergeSystemMessages(t *testing.T) {
	tests := []struct {
		name     string
		messages []api.Message
		want     api.Message
		wantErr  bool
	}{
		{
			name: "combine two messages",
			messages: []api.Message{
				&api.SystemMessage{Content: "First"},
				&api.SystemMessage{Content: "Second"},
			},
			want: &api.SystemMessage{Content: "First\nSecond"},
		},
		{
			name: "preserve last message metadata",
			messages: []api.Message{
				&api.SystemMessage{
					Content:          "First",
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"key1": "value1"}}),
				},
				&api.SystemMessage{
					Content:          "Second",
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"key2": "value2"}}),
				},
			},
			want: &api.SystemMessage{
				Content:          "First\nSecond",
				ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"key2": "value2"}}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeSystemMessages(tt.messages)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMergeUserMessages(t *testing.T) {
	tests := []struct {
		name     string
		messages []api.Message
		want     api.Message
		wantErr  bool
	}{
		{
			name: "combine text blocks",
			messages: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "First"},
					},
				},
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Second"},
					},
				},
			},
			want: &api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "First"},
					&api.TextBlock{Text: "Second"},
				},
			},
		},
		{
			name: "preserve block metadata",
			messages: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{
							Text:             "First",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"block1": "value1"}}),
						},
					},
				},
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{
							Text:             "Second",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"block2": "value2"}}),
						},
					},
				},
			},
			want: &api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{
						Text:             "First",
						ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"block1": "value1"}}),
					},
					&api.TextBlock{
						Text:             "Second",
						ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"block2": "value2"}}),
					},
				},
			},
		},
		{
			name: "preserve message metadata in last block",
			messages: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "First"},
					},
				},
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Second"},
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"msg": "value"}}),
				},
			},
			want: &api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "First"},
					&api.TextBlock{
						Text:             "Second",
						ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"msg": "value"}}),
					},
				},
			},
		},
		{
			name: "mixed content types",
			messages: []api.Message{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Text"},
						&api.ImageBlock{URL: "image.jpg"},
					},
				},
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.FileBlock{URL: "file.txt"},
					},
				},
			},
			want: &api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Text"},
					&api.ImageBlock{URL: "image.jpg"},
					&api.FileBlock{URL: "file.txt"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeUserMessages(tt.messages)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMergeAssistantMessages(t *testing.T) {
	tests := []struct {
		name     string
		messages []api.Message
		want     api.Message
		wantErr  bool
	}{
		{
			name: "combine text blocks",
			messages: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "First"},
					},
				},
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Second"},
					},
				},
			},
			want: &api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "First"},
					&api.TextBlock{Text: "Second"},
				},
			},
		},
		{
			name: "preserve block metadata",
			messages: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{
							Text:             "First",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"block1": "value1"}}),
						},
					},
				},
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{
							Text:             "Second",
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"block2": "value2"}}),
						},
					},
				},
			},
			want: &api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{
						Text:             "First",
						ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"block1": "value1"}}),
					},
					&api.TextBlock{
						Text:             "Second",
						ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"block2": "value2"}}),
					},
				},
			},
		},
		{
			name: "mixed content types",
			messages: []api.Message{
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "Text"},
						&api.ToolCallBlock{
							ToolCallID: "tool-1",
							ToolName:   "test_tool",
						},
					},
				},
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.ReasoningBlock{
							Text:      "Reasoning",
							Signature: "sig-1",
						},
					},
				},
			},
			want: &api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Text"},
					&api.ToolCallBlock{
						ToolCallID: "tool-1",
						ToolName:   "test_tool",
					},
					&api.ReasoningBlock{
						Text:      "Reasoning",
						Signature: "sig-1",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeAssistantMessages(tt.messages)
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeAssistantMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMergeToolMessages(t *testing.T) {
	tests := []struct {
		name     string
		messages []api.Message
		want     api.Message
		wantErr  bool
	}{
		{
			name: "combine tool results",
			messages: []api.Message{
				&api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID: "tool-1",
							ToolName:   "test_tool",
							Result:     json.RawMessage(`"result1"`),
						},
					},
				},
				&api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID: "tool-2",
							ToolName:   "test_tool",
							Result:     json.RawMessage(`"result2"`),
						},
					},
				},
			},
			want: &api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID: "tool-1",
						ToolName:   "test_tool",
						Result:     json.RawMessage(`"result1"`),
					},
					{
						ToolCallID: "tool-2",
						ToolName:   "test_tool",
						Result:     json.RawMessage(`"result2"`),
					},
				},
			},
		},
		{
			name: "preserve metadata",
			messages: []api.Message{
				&api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID:       "tool-1",
							ToolName:         "test_tool",
							Result:           json.RawMessage(`"result1"`),
							ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"block1": "value1"}}),
						},
					},
				},
				&api.ToolMessage{
					Content: []api.ToolResultBlock{
						{
							ToolCallID: "tool-2",
							ToolName:   "test_tool",
							Result:     json.RawMessage(`"result2"`),
						},
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"msg": "value"}}),
				},
			},
			want: &api.ToolMessage{
				Content: []api.ToolResultBlock{
					{
						ToolCallID:       "tool-1",
						ToolName:         "test_tool",
						Result:           json.RawMessage(`"result1"`),
						ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"block1": "value1"}}),
					},
					{
						ToolCallID:       "tool-2",
						ToolName:         "test_tool",
						Result:           json.RawMessage(`"result2"`),
						ProviderMetadata: api.NewProviderMetadata(map[string]any{"metadata": map[string]any{"msg": "value"}}),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeToolMessages(tt.messages)
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeToolMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
