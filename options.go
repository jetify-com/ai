package ai

import (
	"net/http"

	"go.jetify.com/ai/api"
)

type GenerateOptions struct {
	CallOptions api.CallOptions
	Model       api.LanguageModel
}

// GenerateOption is a function that modifies GenerateConfig.
type GenerateOption func(*GenerateOptions)

// WithModel sets the language model to use for generation
func WithModel(model api.LanguageModel) GenerateOption {
	return func(o *GenerateOptions) {
		o.Model = model
	}
}

// WithMaxOutputTokens specifies the maximum number of tokens to generate
func WithMaxOutputTokens(maxTokens int) GenerateOption {
	return func(o *GenerateOptions) {
		o.CallOptions.MaxOutputTokens = maxTokens
	}
}

// WithTemperature controls randomness in the model's output.
// It is recommended to set either Temperature or TopP, but not both.
func WithTemperature(temperature float64) GenerateOption {
	return func(o *GenerateOptions) {
		o.CallOptions.Temperature = &temperature
	}
}

// WithStopSequences specifies sequences that will stop generation when produced.
// Providers may have limits on the number of stop sequences.
func WithStopSequences(stopSequences ...string) GenerateOption {
	return func(o *GenerateOptions) {
		o.CallOptions.StopSequences = stopSequences
	}
}

// WithTopP controls nucleus sampling.
// It is recommended to set either Temperature or TopP, but not both.
func WithTopP(topP float64) GenerateOption {
	return func(o *GenerateOptions) {
		o.CallOptions.TopP = topP
	}
}

// WithTopK limits sampling to the top K options for each token.
// Used to remove "long tail" low probability responses.
// Recommended for advanced use cases only.
func WithTopK(topK int) GenerateOption {
	return func(o *GenerateOptions) {
		o.CallOptions.TopK = topK
	}
}

// WithPresencePenalty affects the likelihood of the model repeating
// information that is already in the prompt
func WithPresencePenalty(penalty float64) GenerateOption {
	return func(o *GenerateOptions) {
		o.CallOptions.PresencePenalty = penalty
	}
}

// WithFrequencyPenalty affects the likelihood of the model
// repeatedly using the same words or phrases
func WithFrequencyPenalty(penalty float64) GenerateOption {
	return func(o *GenerateOptions) {
		o.CallOptions.FrequencyPenalty = penalty
	}
}

// WithResponseFormat specifies whether the output should be text or JSON.
// For JSON output, a schema can optionally guide the model.
func WithResponseFormat(format *api.ResponseFormat) GenerateOption {
	return func(o *GenerateOptions) {
		o.CallOptions.ResponseFormat = format
	}
}

// WithSeed provides an integer seed for random sampling.
// If supported by the model, calls will generate deterministic results.
func WithSeed(seed int) GenerateOption {
	return func(o *GenerateOptions) {
		o.CallOptions.Seed = seed
	}
}

// WithHeaders specifies additional HTTP headers to send with the request.
// Only applicable for HTTP-based providers.
func WithHeaders(headers http.Header) GenerateOption {
	return func(o *GenerateOptions) {
		o.CallOptions.Headers = headers
	}
}

// WithTools specifies the tools available for the model to use during generation.
func WithTools(tools ...api.ToolDefinition) GenerateOption {
	return func(o *GenerateOptions) {
		o.CallOptions.Tools = tools
	}
}

// WithToolChoice specifies how the model should select which tool to use.
func WithToolChoice(toolChoice *api.ToolChoice) GenerateOption {
	return func(o *GenerateOptions) {
		o.CallOptions.ToolChoice = toolChoice
	}
}

// WithProviderMetadata sets additional provider-specific metadata.
// The metadata is passed through to the provider from the AI SDK and enables
// provider-specific functionality that can be fully encapsulated in the provider.
func WithProviderMetadata(providerName string, metadata any) GenerateOption {
	return func(o *GenerateOptions) {
		if o.CallOptions.ProviderMetadata == nil {
			o.CallOptions.ProviderMetadata = api.NewProviderMetadata(map[string]any{})
		}
		o.CallOptions.ProviderMetadata.Set(providerName, metadata)
	}
}

// buildGenerateConfig combines multiple generate options into a single GenerateConfig struct.
func buildGenerateConfig(opts []GenerateOption) GenerateOptions {
	config := GenerateOptions{
		CallOptions: api.CallOptions{
			ProviderMetadata: api.NewProviderMetadata(map[string]any{}),
		},
		Model: DefaultLanguageModel(),
	}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
