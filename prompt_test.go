package aisdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

func TestUserMessage(t *testing.T) {
	tests := []struct {
		name    string
		args    []any
		want    api.UserMessage
		wantErr bool
		errMsg  string
	}{
		{
			name: "simple string content",
			args: []any{"hello world"},
			want: api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello world"},
				},
			},
		},
		{
			name: "multiple content blocks",
			args: []any{
				"first block",
				&api.TextBlock{Text: "second block"},
			},
			want: api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "first block"},
					&api.TextBlock{Text: "second block"},
				},
			},
		},
		{
			name: "with metadata",
			args: []any{
				"content",
				api.NewProviderMetadata(map[string]any{"test": map[string]any{"key": "value"}}),
			},
			want: api.UserMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "content"},
				},
				ProviderMetadata: api.NewProviderMetadata(map[string]any{"test": map[string]any{"key": "value"}}),
			},
		},
		{
			name:    "invalid argument type",
			args:    []any{42},
			wantErr: true,
			errMsg:  "unsupported argument type: int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UserMessage(tt.args...)
			if tt.wantErr {
				require.True(t, result.IsErr())
				assert.Contains(t, result.Err().Error(), tt.errMsg)
				return
			}
			require.True(t, result.IsOk())
			assert.Equal(t, tt.want, result.MustGet())
		})
	}
}

func TestAssistantMessage(t *testing.T) {
	tests := []struct {
		name    string
		args    []any
		want    api.AssistantMessage
		wantErr bool
		errMsg  string
	}{
		{
			name: "simple string content",
			args: []any{"hello world"},
			want: api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello world"},
				},
			},
		},
		{
			name: "multiple content blocks",
			args: []any{
				"first block",
				&api.TextBlock{Text: "second block"},
			},
			want: api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "first block"},
					&api.TextBlock{Text: "second block"},
				},
			},
		},
		{
			name: "with metadata",
			args: []any{
				"content",
				api.NewProviderMetadata(map[string]any{"test": map[string]any{"key": "value"}}),
			},
			want: api.AssistantMessage{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "content"},
				},
				ProviderMetadata: api.NewProviderMetadata(map[string]any{"test": map[string]any{"key": "value"}}),
			},
		},
		{
			name:    "invalid argument type",
			args:    []any{42},
			wantErr: true,
			errMsg:  "unsupported argument type: int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AssistantMessage(tt.args...)
			if tt.wantErr {
				require.True(t, result.IsErr())
				assert.Contains(t, result.Err().Error(), tt.errMsg)
				return
			}
			require.True(t, result.IsOk())
			assert.Equal(t, tt.want, result.MustGet())
		})
	}
}

func TestSystemMessage(t *testing.T) {
	tests := []struct {
		name    string
		args    []any
		want    api.SystemMessage
		wantErr bool
		errMsg  string
	}{
		{
			name: "simple string content",
			args: []any{"system instruction"},
			want: api.SystemMessage{
				Content: "system instruction",
			},
		},
		{
			name: "with metadata",
			args: []any{
				"system instruction",
				api.NewProviderMetadata(map[string]any{"test": map[string]any{"key": "value"}}),
			},
			want: api.SystemMessage{
				Content:          "system instruction",
				ProviderMetadata: api.NewProviderMetadata(map[string]any{"test": map[string]any{"key": "value"}}),
			},
		},
		{
			name:    "multiple string contents",
			args:    []any{"first", "second"},
			wantErr: true,
			errMsg:  "multiple string contents provided for SystemMessage",
		},
		{
			name:    "no content",
			args:    []any{},
			wantErr: true,
			errMsg:  "no content provided for SystemMessage",
		},
		{
			name:    "invalid argument type",
			args:    []any{42},
			wantErr: true,
			errMsg:  "unsupported argument type for SystemMessage: int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SystemMessage(tt.args...)
			if tt.wantErr {
				require.True(t, result.IsErr())
				assert.Contains(t, result.Err().Error(), tt.errMsg)
				return
			}
			require.True(t, result.IsOk())
			assert.Equal(t, tt.want, result.MustGet())
		})
	}
}
