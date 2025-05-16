package aisdk

import (
	"context"

	"go.jetify.com/ai/api"
)

// GenerateText generates a text response for a given prompt using a language model.
// This function does not stream its output.
//
// It returns a [api.Response] containing the generated text, the results of
// any tool calls, and additional information.
//
// It supports either a string argument, which will be converted to a
// [api.UserMessage] with
//
//	GenerateText(ctx, "Hello, world!")
//
// Or a series of [api.Message] arguments:
//
//	GenerateText(ctx,
//	  UserMessage("Show me a picture of a cat"),
//	  AssistantMessage(
//	    "Here is a picture of a cat",
//	    ImageBlock{URL: "https://example.com/cat.png"},
//	  ),
//	)
//
// The last argument can optionally be a series of [api.CallOption] arguments:
//
//	GenerateText(ctx, "Hello, world!", WithMaxTokens(100))
func GenerateText(ctx context.Context, args ...any) (api.Response, error) {
	llmArgs, err := toLLMArgs(args...)
	if err != nil {
		return api.Response{}, err
	}

	return generate(ctx, llmArgs.Prompt, llmArgs.Config)
}

func generate(ctx context.Context, prompt []api.Message, config GenerateTextConfig) (api.Response, error) {
	return config.Model.Generate(ctx, prompt, config.CallOptions)
}

type GenerateTextConfig struct {
	CallOptions api.CallOptions
	Model       api.LanguageModel
}

// GenerateOption is a function that modifies GenerateConfig.
type GenerateOption func(*GenerateTextConfig)

// WithModel sets the language model to use for generation
func WithModel(model api.LanguageModel) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.Model = model
	}
}

// WithMaxTokens specifies the maximum number of tokens to generate
func WithMaxTokens(maxTokens int) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.MaxTokens = maxTokens
	}
}

// WithTemperature controls randomness in the model's output.
// It is recommended to set either Temperature or TopP, but not both.
func WithTemperature(temperature float64) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.Temperature = &temperature
	}
}

// WithStopSequences specifies sequences that will stop generation when produced.
// Providers may have limits on the number of stop sequences.
func WithStopSequences(stopSequences ...string) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.StopSequences = stopSequences
	}
}

// WithTopP controls nucleus sampling.
// It is recommended to set either Temperature or TopP, but not both.
func WithTopP(topP float64) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.TopP = topP
	}
}

// WithTopK limits sampling to the top K options for each token.
// Used to remove "long tail" low probability responses.
// Recommended for advanced use cases only.
func WithTopK(topK int) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.TopK = topK
	}
}

// WithPresencePenalty affects the likelihood of the model repeating
// information that is already in the prompt
func WithPresencePenalty(penalty float64) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.PresencePenalty = penalty
	}
}

// WithFrequencyPenalty affects the likelihood of the model
// repeatedly using the same words or phrases
func WithFrequencyPenalty(penalty float64) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.FrequencyPenalty = penalty
	}
}

// WithResponseFormat specifies whether the output should be text or JSON.
// For JSON output, a schema can optionally guide the model.
func WithResponseFormat(format *api.ResponseFormat) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.ResponseFormat = format
	}
}

// WithSeed provides an integer seed for random sampling.
// If supported by the model, calls will generate deterministic results.
func WithSeed(seed int) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.Seed = seed
	}
}

// WithHeaders specifies additional HTTP headers to send with the request.
// Only applicable for HTTP-based providers.
func WithHeaders(headers map[string]string) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.Headers = headers
	}
}

// WithInputFormat specifies whether the user provided the input as messages or as a prompt.
// This can help guide non-chat models in the expansion, as different expansions
// may be needed for chat vs non-chat use cases.
func WithInputFormat(format api.InputFormat) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.InputFormat = format
	}
}

// WithMode affects the behavior of the language model. It is required to
// support provider-independent streaming and generation of structured objects.
// The model can take this information and e.g. configure json mode, the correct
// low level grammar, etc. It can also be used to optimize the efficiency of the
// streaming, e.g. tool-delta stream parts are only needed in the
// object-tool mode.
//
// Mode will be removed in v2, and at that point it will be deprecated.
// All necessary settings will be directly supported through the call settings,
// in particular responseFormat, toolChoice, and tools.
func WithMode(mode api.ModeConfig) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.Mode = mode
	}
}

// WithProviderMetadata sets additional provider-specific metadata.
// The metadata is passed through to the provider from the AI SDK and enables
// provider-specific functionality that can be fully encapsulated in the provider.
func WithProviderMetadata(providerName string, metadata any) GenerateOption {
	return func(o *GenerateTextConfig) {
		if o.CallOptions.ProviderMetadata == nil {
			o.CallOptions.ProviderMetadata = api.NewProviderMetadata(map[string]any{})
		}
		o.CallOptions.ProviderMetadata.Set(providerName, metadata)
	}
}

// WithTools specifies the tools available to the model during generation.
func WithTools(tools ...api.ToolDefinition) GenerateOption {
	return func(o *GenerateTextConfig) {
		o.CallOptions.Mode = api.RegularMode{Tools: tools}
	}
}

// buildGenerateConfig combines multiple generate options into a single GenerateConfig struct.
func buildGenerateConfig(opts []GenerateOption) GenerateTextConfig {
	config := GenerateTextConfig{
		CallOptions: api.CallOptions{
			InputFormat:      "prompt",          // default
			Mode:             api.RegularMode{}, // default
			ProviderMetadata: api.NewProviderMetadata(map[string]any{}),
		},
		Model: DefaultLanguageModel(),
	}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
