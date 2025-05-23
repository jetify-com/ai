package ai

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
func GenerateText(ctx context.Context, prompt []api.Message, opts ...GenerateOption) (api.Response, error) {
	config := buildGenerateConfig(opts)
	return generate(ctx, prompt, config)
}

func generate(ctx context.Context, prompt []api.Message, opts GenerateOptions) (api.Response, error) {
	return opts.Model.Generate(ctx, prompt, opts.CallOptions)
}
