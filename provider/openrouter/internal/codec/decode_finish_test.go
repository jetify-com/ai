package codec

import (
	"testing"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/openrouter/internal/client"
)

func TestDecodeFinishReason(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected api.FinishReason
	}{
		{
			name:     "stop reason",
			input:    client.FinishReasonStop,
			expected: api.FinishReasonStop,
		},
		{
			name:     "length reason",
			input:    client.FinishReasonLength,
			expected: api.FinishReasonLength,
		},
		{
			name:     "content filter reason",
			input:    client.FinishReasonContentFilter,
			expected: api.FinishReasonContentFilter,
		},
		{
			name:     "function call reason",
			input:    client.FinishReasonFunctionCall,
			expected: api.FinishReasonToolCalls,
		},
		{
			name:     "tool calls reason",
			input:    client.FinishReasonToolCalls,
			expected: api.FinishReasonToolCalls,
		},
		{
			name:     "empty string",
			input:    "",
			expected: api.FinishReasonUnknown,
		},
		{
			name:     "unknown reason",
			input:    "something_else",
			expected: api.FinishReasonUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecodeFinishReason(tt.input)
			if result != tt.expected {
				t.Errorf("DecodeFinishReason(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
