package ai

import (
	"context"

	"go.jetify.com/ai/api"
)

// GenerateText uses a language model to generate a text response from a given prompt.
//
// This function does not stream its output.
//
// It returns a [api.Response] containing the generated text, the results of
// any tool calls, and additional information.
//
// A prompt is a sequence of [api.Message]s:
//
//	GenerateText(ctx, []api.Message{
//		&api.UserMessage{
//			Content: []api.ContentBlock{
//				&api.TextBlock{Text: "Show me a picture of a cat"},
//			},
//		},
//		&api.AssistantMessage{
//			Content: []api.ContentBlock{
//				&api.TextBlock{Text: "Here is a picture of a cat"},
//				&api.ImageBlock{URL: "https://example.com/cat.png"},
//			},
//		},
//	})
//
// The last argument can optionally be a series of [GenerateOption] arguments:
//
//	GenerateText(ctx, messages, WithMaxTokens(100))
func GenerateText(ctx context.Context, prompt []api.Message, opts ...GenerateOption) (*api.Response, error) {
	config := buildGenerateConfig(opts)
	return generate(ctx, prompt, config)
}

// GenerateTextStr uses a language model to generate a text response from a given string prompt.
//
// It is a convenience wrapper around GenerateText for simple string-based prompts.
//
// Example usage:
//
//	GenerateTextStr(ctx, "Write a brief summary of the benefits of renewable energy")
//
// The function can optionally take [GenerateOption] arguments:
//
//	GenerateTextStr(ctx, "Explain the key differences between REST and GraphQL APIs", WithMaxTokens(500))
//
// The string prompt is automatically converted to a [api.UserMessage] before
// being passed to GenerateText.
func GenerateTextStr(ctx context.Context, prompt string, opts ...GenerateOption) (*api.Response, error) {
	msg := &api.UserMessage{
		Content: []api.ContentBlock{&api.TextBlock{Text: prompt}},
	}
	return GenerateText(ctx, []api.Message{msg}, opts...)
}

func generate(ctx context.Context, prompt []api.Message, opts GenerateOptions) (*api.Response, error) {
	return opts.Model.Generate(ctx, prompt, opts.CallOptions)
}

// StreamText uses a language model to generate a streaming text response from a given prompt.
//
// This function streams its output as a sequence of events.
//
// It returns a [api.StreamResponse] containing a stream of events from the model,
// including partial text, tool calls, and additional information.
//
// A prompt is a sequence of [api.Message]s:
//
//	StreamText(ctx, []api.Message{
//		&api.UserMessage{
//			Content: []api.ContentBlock{
//				&api.TextBlock{Text: "Show me a picture of a cat"},
//			},
//		},
//		&api.AssistantMessage{
//			Content: []api.ContentBlock{
//				&api.TextBlock{Text: "Here is a picture of a cat"},
//				&api.ImageBlock{URL: "https://example.com/cat.png"},
//			},
//		},
//	})
//
// The last argument can optionally be a series of [GenerateOption] arguments:
//
//	StreamText(ctx, messages, WithMaxTokens(100))
func StreamText(ctx context.Context, prompt []api.Message, opts ...GenerateOption) (*api.StreamResponse, error) {
	config := buildGenerateConfig(opts)
	return stream(ctx, prompt, config)
}

// StreamTextStr uses a language model to generate a streaming text response from a given string prompt.
//
// It is a convenience wrapper around StreamText for simple string-based prompts.
//
// Example usage:
//
//	StreamTextStr(ctx, "Write a brief summary of the benefits of renewable energy")
//
// The function can optionally take [GenerateOption] arguments:
//
//	StreamTextStr(ctx, "Explain the key differences between REST and GraphQL APIs", WithMaxTokens(500))
//
// The string prompt is automatically converted to a [api.UserMessage] before
// being passed to StreamText.
func StreamTextStr(ctx context.Context, prompt string, opts ...GenerateOption) (*api.StreamResponse, error) {
	msg := &api.UserMessage{
		Content: []api.ContentBlock{&api.TextBlock{Text: prompt}},
	}
	return StreamText(ctx, []api.Message{msg}, opts...)
}

func stream(ctx context.Context, prompt []api.Message, opts GenerateOptions) (*api.StreamResponse, error) {
	return opts.Model.Stream(ctx, prompt, opts.CallOptions)
}
