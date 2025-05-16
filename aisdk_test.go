package aisdk

import (
	"context"
	"net/url"
	"testing"

	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/api"
)

func TestCallOptionBuilders(t *testing.T) {
	tests := []struct {
		name     string
		option   GenerateOption
		validate func(*testing.T, *GenerateTextConfig)
	}{
		{
			name:   "WithMaxTokens",
			option: WithMaxTokens(100),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				assert.Equal(t, 100, opts.CallOptions.MaxTokens)
			},
		},
		{
			name:   "WithTemperature",
			option: WithTemperature(0.7),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				assert.NotNil(t, opts.CallOptions.Temperature)
				assert.Equal(t, 0.7, *opts.CallOptions.Temperature)
			},
		},
		{
			name:   "WithStopSequences",
			option: WithStopSequences("stop1", "stop2"),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				assert.Equal(t, []string{"stop1", "stop2"}, opts.CallOptions.StopSequences)
			},
		},
		{
			name:   "WithStopSequences_Empty",
			option: WithStopSequences(),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				assert.Empty(t, opts.CallOptions.StopSequences)
			},
		},
		{
			name:   "WithTopP",
			option: WithTopP(0.9),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				assert.Equal(t, 0.9, opts.CallOptions.TopP)
			},
		},
		{
			name:   "WithTopK",
			option: WithTopK(40),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				assert.Equal(t, 40, opts.CallOptions.TopK)
			},
		},
		{
			name:   "WithPresencePenalty",
			option: WithPresencePenalty(1.0),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				assert.Equal(t, 1.0, opts.CallOptions.PresencePenalty)
			},
		},
		{
			name:   "WithFrequencyPenalty",
			option: WithFrequencyPenalty(1.5),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				assert.Equal(t, 1.5, opts.CallOptions.FrequencyPenalty)
			},
		},
		{
			name: "WithResponseFormat",
			option: WithResponseFormat(&api.ResponseFormat{
				Type:        "json",
				Schema:      &jsonschema.Definition{},
				Name:        "test",
				Description: "test desc",
			}),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				assert.Equal(t, "json", opts.CallOptions.ResponseFormat.Type)
				assert.NotNil(t, opts.CallOptions.ResponseFormat.Schema)
				assert.Equal(t, "test", opts.CallOptions.ResponseFormat.Name)
				assert.Equal(t, "test desc", opts.CallOptions.ResponseFormat.Description)
			},
		},
		{
			name:   "WithSeed",
			option: WithSeed(42),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				assert.Equal(t, 42, opts.CallOptions.Seed)
			},
		},
		{
			name:   "WithHeaders",
			option: WithHeaders(map[string]string{"key": "value"}),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				assert.Equal(t, map[string]string{"key": "value"}, opts.CallOptions.Headers)
			},
		},
		{
			name:   "WithInputFormat",
			option: WithInputFormat(api.InputFormatMessages),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				assert.Equal(t, api.InputFormatMessages, opts.CallOptions.InputFormat)
			},
		},
		{
			name: "WithMode",
			option: WithMode(api.ObjectJSONMode{
				Name:        "test",
				Description: "test desc",
			}),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				mode, ok := opts.CallOptions.Mode.(api.ObjectJSONMode)
				assert.True(t, ok)
				assert.Equal(t, "test", mode.Name)
				assert.Equal(t, "test desc", mode.Description)
			},
		},
		{
			name: "WithProviderMetadata_SingleProvider",
			option: WithProviderMetadata("test-provider", map[string]any{
				"key": "value",
			}),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				expected := api.NewProviderMetadata(map[string]any{
					"test-provider": map[string]any{
						"key": "value",
					},
				})
				assert.Equal(t, expected, opts.CallOptions.ProviderMetadata)
			},
		},
		{
			name: "WithProviderMetadata_MultipleProviders",
			option: func() GenerateOption {
				return func(o *GenerateTextConfig) {
					WithProviderMetadata("provider1", map[string]any{"key1": "value1"})(o)
					WithProviderMetadata("provider2", map[string]any{"key2": "value2"})(o)
				}
			}(),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				expected := api.NewProviderMetadata(map[string]any{
					"provider1": map[string]any{"key1": "value1"},
					"provider2": map[string]any{"key2": "value2"},
				})
				assert.Equal(t, expected, opts.CallOptions.ProviderMetadata)
			},
		},
		{
			name:   "WithModel",
			option: WithModel(&mockLanguageModel{name: "test-model"}),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				model, ok := opts.Model.(*mockLanguageModel)
				assert.True(t, ok)
				assert.Equal(t, "test-model", model.name)
			},
		},
		{
			name:   "WithTools",
			option: WithTools(api.FunctionTool{Name: "test-tool"}),
			validate: func(t *testing.T, opts *GenerateTextConfig) {
				mode, ok := opts.CallOptions.Mode.(api.RegularMode)
				assert.True(t, ok)
				assert.Len(t, mode.Tools, 1)
				tool, ok := mode.Tools[0].(api.FunctionTool)
				assert.True(t, ok)
				assert.Equal(t, "test-tool", tool.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &GenerateTextConfig{}
			tt.option(opts)
			tt.validate(t, opts)
		})
	}
}

func TestBuildCallOptions(t *testing.T) {
	tests := []struct {
		name     string
		opts     []GenerateOption
		validate func(*testing.T, GenerateTextConfig)
	}{
		{
			name: "Default options",
			opts: []GenerateOption{},
			validate: func(t *testing.T, opts GenerateTextConfig) {
				assert.Equal(t, api.InputFormatPrompt, opts.CallOptions.InputFormat)
				assert.Equal(t, "regular", opts.CallOptions.Mode.Type())
			},
		},
		{
			name: "Multiple options",
			opts: []GenerateOption{
				WithMaxTokens(100),
				WithTemperature(0.7),
				WithInputFormat(api.InputFormatMessages),
			},
			validate: func(t *testing.T, opts GenerateTextConfig) {
				assert.Equal(t, 100, opts.CallOptions.MaxTokens)
				assert.NotNil(t, opts.CallOptions.Temperature)
				assert.Equal(t, 0.7, *opts.CallOptions.Temperature)
				assert.Equal(t, api.InputFormatMessages, opts.CallOptions.InputFormat)
			},
		},
		{
			name: "Override defaults",
			opts: []GenerateOption{
				WithMode(api.ObjectJSONMode{Name: "test"}),
				WithInputFormat(api.InputFormatPrompt),
			},
			validate: func(t *testing.T, opts GenerateTextConfig) {
				assert.Equal(t, api.InputFormatPrompt, opts.CallOptions.InputFormat)
				mode, ok := opts.CallOptions.Mode.(api.ObjectJSONMode)
				assert.True(t, ok)
				assert.Equal(t, "test", mode.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := buildGenerateConfig(tt.opts)
			tt.validate(t, opts)
		})
	}
}

// mockLanguageModel implements api.LanguageModel for testing
type mockLanguageModel struct {
	name string
}

func (m *mockLanguageModel) Generate(ctx context.Context, prompt []api.Message, opts api.CallOptions) (api.Response, error) {
	return api.Response{}, nil
}

func (m *mockLanguageModel) Stream(ctx context.Context, prompt []api.Message, opts api.CallOptions) (api.StreamResponse, error) {
	return api.StreamResponse{}, nil
}

func (m *mockLanguageModel) DefaultObjectGenerationMode() api.ObjectGenerationMode {
	return "json"
}

func (m *mockLanguageModel) ModelID() string {
	return m.name
}

func (m *mockLanguageModel) ProviderName() string {
	return "mock-provider"
}

func (m *mockLanguageModel) SpecificationVersion() string {
	return "v1"
}

func (m *mockLanguageModel) SupportsImageURLs() bool {
	return true
}

func (m *mockLanguageModel) SupportsStructuredOutputs() bool {
	return false
}

func (m *mockLanguageModel) SupportsURL(u *url.URL) bool {
	return false
}
