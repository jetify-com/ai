package aisdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
	"go.jetify.com/pkg/try"
)

func TestToLLMArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []any
		want    *llmArgs
		wantErr bool
		errMsg  string
	}{
		{
			name: "simple string content",
			args: []any{"hello world"},
			want: &llmArgs{
				Prompt: []api.Message{
					&api.UserMessage{
						Content: []api.ContentBlock{
							&api.TextBlock{Text: "hello world"},
						},
					},
				},
				Config: GenerateTextConfig{},
			},
		},
		{
			name: "multiple content blocks",
			args: []any{
				"first block",
				&api.TextBlock{Text: "second block"},
			},
			want: &llmArgs{
				Prompt: []api.Message{
					&api.UserMessage{
						Content: []api.ContentBlock{
							&api.TextBlock{Text: "first block"},
							&api.TextBlock{Text: "second block"},
						},
					},
				},
				Config: GenerateTextConfig{},
			},
		},
		{
			name: "with messages",
			args: []any{
				&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "user message"},
					},
				},
				&api.AssistantMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "assistant message"},
					},
				},
			},
			want: &llmArgs{
				Prompt: []api.Message{
					&api.UserMessage{
						Content: []api.ContentBlock{
							&api.TextBlock{Text: "user message"},
						},
					},
					&api.AssistantMessage{
						Content: []api.ContentBlock{
							&api.TextBlock{Text: "assistant message"},
						},
					},
				},
				Config: GenerateTextConfig{},
			},
		},
		{
			name: "with options",
			args: []any{
				"hello world",
				WithMaxTokens(100),
				WithTemperature(0.7),
			},
			want: &llmArgs{
				Prompt: []api.Message{
					&api.UserMessage{
						Content: []api.ContentBlock{
							&api.TextBlock{Text: "hello world"},
						},
					},
				},
				Config: GenerateTextConfig{
					CallOptions: api.CallOptions{
						MaxTokens:   100,
						Temperature: float64Value(0.7),
					},
				},
			},
		},
		{
			name:    "invalid argument type",
			args:    []any{42},
			wantErr: true,
			errMsg:  "unsupported argument type: int",
		},
		{
			name: "pointer to content block",
			args: []any{
				func() *api.ContentBlock {
					block := api.ContentBlock(&api.TextBlock{Text: "pointer block"})
					return &block
				}(),
			},
			want: &llmArgs{
				Prompt: []api.Message{
					&api.UserMessage{
						Content: []api.ContentBlock{
							&api.TextBlock{Text: "pointer block"},
						},
					},
				},
				Config: GenerateTextConfig{},
			},
		},
		{
			name: "try-wrapped content block",
			args: []any{
				try.Ok[api.ContentBlock](&api.TextBlock{Text: "try block"}),
			},
			want: &llmArgs{
				Prompt: []api.Message{
					&api.UserMessage{
						Content: []api.ContentBlock{
							&api.TextBlock{Text: "try block"},
						},
					},
				},
				Config: GenerateTextConfig{},
			},
		},
		{
			name: "try-wrapped pointer to content block",
			args: []any{
				func() try.Try[*api.ContentBlock] {
					block := api.ContentBlock(&api.TextBlock{Text: "try pointer block"})
					return try.Ok(&block)
				}(),
			},
			want: &llmArgs{
				Prompt: []api.Message{
					&api.UserMessage{
						Content: []api.ContentBlock{
							&api.TextBlock{Text: "try pointer block"},
						},
					},
				},
				Config: GenerateTextConfig{},
			},
		},
		{
			name: "try-wrapped message",
			args: []any{
				try.Ok[api.Message](&api.UserMessage{
					Content: []api.ContentBlock{
						&api.TextBlock{Text: "try message"},
					},
				}),
			},
			want: &llmArgs{
				Prompt: []api.Message{
					&api.UserMessage{
						Content: []api.ContentBlock{
							&api.TextBlock{Text: "try message"},
						},
					},
				},
				Config: GenerateTextConfig{},
			},
		},
		{
			name: "error in try-wrapped content block",
			args: []any{
				try.Err[api.ContentBlock](assert.AnError),
			},
			wantErr: true,
			errMsg:  assert.AnError.Error(),
		},
		{
			name: "error in try-wrapped message",
			args: []any{
				try.Err[api.Message](assert.AnError),
			},
			wantErr: true,
			errMsg:  assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toLLMArgs(tt.args...)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}
			require.NoError(t, err)

			// Compare prompts
			assert.Equal(t, tt.want.Prompt, got.Prompt)

			// Compare GenerateConfig
			assert.Equal(t, tt.want.Config, got.Config)
		})
	}
}

// Helper function to create a pointer to a float64
func float64Value(v float64) *float64 {
	return &v
}
