package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/openrouter/internal/client"
)

func TestDecodeLogProbs(t *testing.T) {
	tests := []struct {
		name     string
		input    *client.LogProbs
		expected []api.LogProb
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: []api.LogProb{},
		},
		{
			name:     "nil content",
			input:    &client.LogProbs{Content: nil},
			expected: []api.LogProb{},
		},
		{
			name: "single token with no top logprobs",
			input: &client.LogProbs{
				Content: []client.LogProb{
					{
						Token:       "hello",
						LogProb:     -0.5,
						TopLogProbs: nil,
					},
				},
			},
			expected: []api.LogProb{
				{
					Token:       "hello",
					LogProb:     -0.5,
					TopLogProbs: []api.TokenLogProb{},
				},
			},
		},
		{
			name: "single token with top logprobs",
			input: &client.LogProbs{
				Content: []client.LogProb{
					{
						Token:   "hello",
						LogProb: -0.5,
						TopLogProbs: []client.TopLogProb{
							{Token: "hello", LogProb: -0.5},
							{Token: "hi", LogProb: -1.0},
						},
					},
				},
			},
			expected: []api.LogProb{
				{
					Token:   "hello",
					LogProb: -0.5,
					TopLogProbs: []api.TokenLogProb{
						{Token: "hello", LogProb: -0.5},
						{Token: "hi", LogProb: -1.0},
					},
				},
			},
		},
		{
			name: "multiple tokens with mixed top logprobs",
			input: &client.LogProbs{
				Content: []client.LogProb{
					{
						Token:   "hello",
						LogProb: -0.5,
						TopLogProbs: []client.TopLogProb{
							{Token: "hello", LogProb: -0.5},
						},
					},
					{
						Token:       "world",
						LogProb:     -0.3,
						TopLogProbs: nil,
					},
				},
			},
			expected: []api.LogProb{
				{
					Token:   "hello",
					LogProb: -0.5,
					TopLogProbs: []api.TokenLogProb{
						{Token: "hello", LogProb: -0.5},
					},
				},
				{
					Token:       "world",
					LogProb:     -0.3,
					TopLogProbs: []api.TokenLogProb{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecodeLogProbs(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDecodeCompletionLogProbs(t *testing.T) {
	tests := []struct {
		name     string
		input    *client.CompletionLogProbs
		expected []api.LogProb
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: []api.LogProb{},
		},
		{
			name: "empty tokens",
			input: &client.CompletionLogProbs{
				Tokens:        []string{},
				TokenLogProbs: []float64{},
				TopLogProbs:   nil,
			},
			expected: []api.LogProb{},
		},
		{
			name: "single token without top logprobs",
			input: &client.CompletionLogProbs{
				Tokens:        []string{"hello"},
				TokenLogProbs: []float64{-0.5},
				TopLogProbs:   nil,
			},
			expected: []api.LogProb{
				{
					Token:       "hello",
					LogProb:     -0.5,
					TopLogProbs: []api.TokenLogProb{},
				},
			},
		},
		{
			name: "token with missing logprob defaults to 0",
			input: &client.CompletionLogProbs{
				Tokens:        []string{"hello"},
				TokenLogProbs: []float64{},
				TopLogProbs:   nil,
			},
			expected: []api.LogProb{
				{
					Token:       "hello",
					LogProb:     0,
					TopLogProbs: []api.TokenLogProb{},
				},
			},
		},
		{
			name: "token with top logprobs",
			input: &client.CompletionLogProbs{
				Tokens:        []string{"hello"},
				TokenLogProbs: []float64{-0.5},
				TopLogProbs: []map[string]float64{
					{
						"hello": -0.5,
						"hi":    -1.0,
					},
				},
			},
			expected: []api.LogProb{
				{
					Token:   "hello",
					LogProb: -0.5,
					TopLogProbs: []api.TokenLogProb{
						{Token: "hello", LogProb: -0.5},
						{Token: "hi", LogProb: -1.0},
					},
				},
			},
		},
		{
			name: "multiple tokens with mixed top logprobs",
			input: &client.CompletionLogProbs{
				Tokens:        []string{"hello", "world"},
				TokenLogProbs: []float64{-0.5, -0.3},
				TopLogProbs: []map[string]float64{
					{
						"hello": -0.5,
						"hi":    -1.0,
					},
					nil,
				},
			},
			expected: []api.LogProb{
				{
					Token:   "hello",
					LogProb: -0.5,
					TopLogProbs: []api.TokenLogProb{
						{Token: "hello", LogProb: -0.5},
						{Token: "hi", LogProb: -1.0},
					},
				},
				{
					Token:       "world",
					LogProb:     -0.3,
					TopLogProbs: []api.TokenLogProb{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecodeCompletionLogProbs(tt.input)
			assert.Equal(t, len(tt.expected), len(result))
			for i := range result {
				assert.Equal(t, tt.expected[i].Token, result[i].Token)
				assert.Equal(t, tt.expected[i].LogProb, result[i].LogProb)
				assert.ElementsMatch(t, tt.expected[i].TopLogProbs, result[i].TopLogProbs)
			}
		})
	}
}
