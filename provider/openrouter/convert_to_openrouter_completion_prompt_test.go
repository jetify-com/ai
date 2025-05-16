package openrouter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/api"
)

func TestConvertToOpenRouterCompletionPrompt(t *testing.T) {
	t.Run("direct prompt", func(t *testing.T) {
		prompt := []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello world"},
				},
			},
		}

		text, stop, err := ConvertToOpenRouterCompletionPrompt(CompletionPromptOptions{
			Prompt:      prompt,
			InputFormat: InputFormatPrompt,
		})

		assert.NoError(t, err)
		assert.Equal(t, "hello world", text)
		assert.Nil(t, stop)
	})

	t.Run("system message prefix", func(t *testing.T) {
		prompt := []api.Message{
			&api.SystemMessage{Content: "system instruction"},
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "user message"},
				},
			},
		}

		text, stop, err := ConvertToOpenRouterCompletionPrompt(CompletionPromptOptions{
			Prompt:      prompt,
			InputFormat: InputFormatMessages,
		})

		assert.NoError(t, err)
		assert.Equal(t, "system instruction\n\nuser:\nuser message\n\nassistant:\n", text)
		assert.Equal(t, []string{"\nuser:"}, stop)
	})

	t.Run("user and assistant messages", func(t *testing.T) {
		prompt := []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello"},
				},
			},
			&api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hi there"},
				},
			},
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "how are you?"},
				},
			},
		}

		text, stop, err := ConvertToOpenRouterCompletionPrompt(CompletionPromptOptions{
			Prompt:      prompt,
			InputFormat: InputFormatMessages,
		})

		assert.NoError(t, err)
		assert.Equal(t, "user:\nhello\n\nassistant:\nhi there\n\nuser:\nhow are you?\n\nassistant:\n", text)
		assert.Equal(t, []string{"\nuser:"}, stop)
	})

	t.Run("custom user and assistant labels", func(t *testing.T) {
		prompt := []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello"},
				},
			},
			&api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hi"},
				},
			},
		}

		text, stop, err := ConvertToOpenRouterCompletionPrompt(CompletionPromptOptions{
			Prompt:      prompt,
			InputFormat: InputFormatMessages,
			User:        "Human",
			Assistant:   "AI",
		})

		assert.NoError(t, err)
		assert.Equal(t, "Human:\nhello\n\nAI:\nhi\n\nAI:\n", text)
		assert.Equal(t, []string{"\nHuman:"}, stop)
	})

	t.Run("unsupported image content", func(t *testing.T) {
		prompt := []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					api.ImageBlockFromURL("http://example.com/image.jpg"),
				},
			},
		}

		_, _, err := ConvertToOpenRouterCompletionPrompt(CompletionPromptOptions{
			Prompt:      prompt,
			InputFormat: InputFormatMessages,
		})

		assert.Error(t, err)
		assert.IsType(t, &api.UnsupportedFunctionalityError{}, err)
		assert.Contains(t, err.Error(), "images")
	})

	t.Run("unsupported file content", func(t *testing.T) {
		prompt := []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					api.FileBlockFromURL("http://example.com/doc.pdf"),
				},
			},
		}

		_, _, err := ConvertToOpenRouterCompletionPrompt(CompletionPromptOptions{
			Prompt:      prompt,
			InputFormat: InputFormatMessages,
		})

		assert.Error(t, err)
		assert.IsType(t, &api.UnsupportedFunctionalityError{}, err)
		assert.Contains(t, err.Error(), "file attachments")
	})

	t.Run("unsupported tool call", func(t *testing.T) {
		prompt := []api.Message{
			&api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "123",
						ToolName:   "calculator",
						Args:       json.RawMessage(`{"x": 1, "y": 2}`),
					},
				},
			},
		}

		_, _, err := ConvertToOpenRouterCompletionPrompt(CompletionPromptOptions{
			Prompt:      prompt,
			InputFormat: InputFormatMessages,
		})

		assert.Error(t, err)
		assert.IsType(t, &api.UnsupportedFunctionalityError{}, err)
		assert.Contains(t, err.Error(), "tool-call messages")
	})

	t.Run("unexpected system message", func(t *testing.T) {
		prompt := []api.Message{
			&api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello"},
				},
			},
			&api.SystemMessage{Content: "unexpected system"},
		}

		_, _, err := ConvertToOpenRouterCompletionPrompt(CompletionPromptOptions{
			Prompt:      prompt,
			InputFormat: InputFormatMessages,
		})

		assert.Error(t, err)
		assert.IsType(t, &api.InvalidPromptError{}, err)
		assert.Contains(t, err.Error(), "unexpected system message")
	})
}
