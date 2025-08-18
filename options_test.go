package ai

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/api"
	"go.jetify.com/pkg/pointer"
)

func TestCallOptionBuilders(t *testing.T) {
	tests := []struct {
		name     string
		option   GenerateOption
		expected GenerateOptions
	}{
		{
			name:   "WithMaxOutputTokens",
			option: WithMaxOutputTokens(100),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{MaxOutputTokens: 100},
			},
		},
		{
			name:   "WithTemperature",
			option: WithTemperature(0.7),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{Temperature: pointer.Float64(0.7)},
			},
		},
		{
			name:   "WithStopSequences",
			option: WithStopSequences("stop1", "stop2"),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{StopSequences: []string{"stop1", "stop2"}},
			},
		},
		{
			name:   "WithStopSequences_Empty",
			option: WithStopSequences(),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{StopSequences: nil},
			},
		},
		{
			name:   "WithTopP",
			option: WithTopP(0.9),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{TopP: 0.9},
			},
		},
		{
			name:   "WithTopK",
			option: WithTopK(40),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{TopK: 40},
			},
		},
		{
			name:   "WithPresencePenalty",
			option: WithPresencePenalty(1.0),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{PresencePenalty: 1.0},
			},
		},
		{
			name:   "WithFrequencyPenalty",
			option: WithFrequencyPenalty(1.5),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{FrequencyPenalty: 1.5},
			},
		},
		{
			name: "WithResponseFormat",
			option: WithResponseFormat(&api.ResponseFormat{
				Type:        "json",
				Schema:      &jsonschema.Schema{},
				Name:        "test",
				Description: "test desc",
			}),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{
					ResponseFormat: &api.ResponseFormat{
						Type:        "json",
						Schema:      &jsonschema.Schema{},
						Name:        "test",
						Description: "test desc",
					},
				},
			},
		},
		{
			name:   "WithSeed",
			option: WithSeed(42),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{Seed: 42},
			},
		},
		{
			name:   "WithHeaders",
			option: WithHeaders(http.Header{"key": []string{"value"}}),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{Headers: http.Header{"key": []string{"value"}}},
			},
		},
		{
			name:   "WithTools",
			option: WithTools(&api.FunctionTool{Name: "test-tool"}),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{
					Tools: []api.ToolDefinition{&api.FunctionTool{Name: "test-tool"}},
				},
			},
		},
		{
			name: "WithProviderMetadata_SingleProvider",
			option: WithProviderMetadata("test-provider", map[string]any{
				"key": "value",
			}),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{
					ProviderMetadata: api.NewProviderMetadata(map[string]any{
						"test-provider": map[string]any{
							"key": "value",
						},
					}),
				},
			},
		},
		{
			name: "WithProviderMetadata_MultipleProviders",
			option: func() GenerateOption {
				return func(o *GenerateOptions) {
					WithProviderMetadata("provider1", map[string]any{"key1": "value1"})(o)
					WithProviderMetadata("provider2", map[string]any{"key2": "value2"})(o)
				}
			}(),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{
					ProviderMetadata: api.NewProviderMetadata(map[string]any{
						"provider1": map[string]any{"key1": "value1"},
						"provider2": map[string]any{"key2": "value2"},
					}),
				},
			},
		},
		{
			name:   "WithModel",
			option: WithModel(&mockLanguageModel{name: "test-model"}),
			expected: GenerateOptions{
				Model: &mockLanguageModel{name: "test-model"},
			},
		},
		{
			name: "WithCallOptions",
			option: WithCallOptions(api.CallOptions{
				MaxOutputTokens:  500,
				Temperature:      pointer.Float64(0.5),
				TopP:             0.8,
				StopSequences:    []string{"END"},
				Seed:             123,
				PresencePenalty:  0.1,
				FrequencyPenalty: 0.2,
			}),
			expected: GenerateOptions{
				CallOptions: api.CallOptions{
					MaxOutputTokens:  500,
					Temperature:      pointer.Float64(0.5),
					TopP:             0.8,
					StopSequences:    []string{"END"},
					Seed:             123,
					PresencePenalty:  0.1,
					FrequencyPenalty: 0.2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &GenerateOptions{}
			tt.option(opts)
			assert.Equal(t, tt.expected, *opts)
		})
	}
}

func TestBuildCallOptions(t *testing.T) {
	tests := []struct {
		name     string
		opts     []GenerateOption
		expected GenerateOptions
	}{
		{
			name: "Default options",
			opts: []GenerateOption{},
			expected: GenerateOptions{
				CallOptions: api.CallOptions{
					ProviderMetadata: api.NewProviderMetadata(map[string]any{}),
				},
				Model: DefaultLanguageModel(),
			},
		},
		{
			name: "Multiple options",
			opts: []GenerateOption{
				WithMaxOutputTokens(100),
				WithTemperature(0.7),
			},
			expected: GenerateOptions{
				CallOptions: api.CallOptions{
					MaxOutputTokens:  100,
					Temperature:      pointer.Float64(0.7),
					ProviderMetadata: api.NewProviderMetadata(map[string]any{}),
				},
				Model: DefaultLanguageModel(),
			},
		},
		{
			name: "With tools",
			opts: []GenerateOption{
				WithTools(&api.FunctionTool{Name: "test-tool"}),
			},
			expected: GenerateOptions{
				CallOptions: api.CallOptions{
					Tools:            []api.ToolDefinition{&api.FunctionTool{Name: "test-tool"}},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{}),
				},
				Model: DefaultLanguageModel(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := buildGenerateConfig(tt.opts)
			assert.Equal(t, tt.expected, opts)
		})
	}
}

// mockLanguageModel implements api.LanguageModel for testing
type mockLanguageModel struct {
	name string
}

func (m *mockLanguageModel) Generate(ctx context.Context, prompt []api.Message, opts api.CallOptions) (*api.Response, error) {
	return &api.Response{}, nil
}

func (m *mockLanguageModel) Stream(ctx context.Context, prompt []api.Message, opts api.CallOptions) (*api.StreamResponse, error) {
	return &api.StreamResponse{}, nil
}

func (m *mockLanguageModel) ModelID() string {
	return m.name
}

func (m *mockLanguageModel) ProviderName() string {
	return "mock-provider"
}

func (m *mockLanguageModel) SupportedUrls() []api.SupportedURL {
	return nil
}
