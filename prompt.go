package aisdk

import (
	"go.jetify.com/ai/api"
	"go.jetify.com/pkg/try"
)

// UserMessage creates an [api.UserMessage] from the provided arguments.
//
// The arguments can be either a string:
//
//	UserMessage("Hello, world!")
//
// or a series of [api.ContentBlock]:
//
//	UserMessage(
//	  api.TextBlock{Text: "Hello, world!"},
//	  api.ImageBlock{URL: "https://example.com/image.png"},
//	)
//
// The last argument can optionally be a [api.ProviderMetadata].
func UserMessage(args ...any) try.Try[api.UserMessage] {
	blocks, metadata, err := processContentArgs(args...)
	if err != nil {
		return try.Errf[api.UserMessage]("error creating UserMessage: %w", err)
	}

	return try.Ok(api.UserMessage{
		Content:          blocks,
		ProviderMetadata: metadata,
	})
}

// AssistantMessage creates an [api.AssistantMessage] from the provided arguments.
//
// The arguments can be either a string:
//
//	AssistantMessage("Hello, world!")
//
// or a series of [api.ContentBlock]:
//
//	AssistantMessage(
//	  api.TextBlock{Text: "Hello, world!"},
//	  api.ImageBlock{URL: "https://example.com/image.png"},
//	)
//
// The last argument can optionally be a [api.ProviderMetadata].
func AssistantMessage(args ...any) try.Try[api.AssistantMessage] {
	blocks, metadata, err := processContentArgs(args...)
	if err != nil {
		return try.Errf[api.AssistantMessage]("error creating AssistantMessage: %w", err)
	}

	return try.Ok(api.AssistantMessage{
		Content:          blocks,
		ProviderMetadata: metadata,
	})
}

// SystemMessage creates an [api.SystemMessage] from the provided string content:
//
//	SystemMessage("You are a helpful assistant.")
//
// The last argument can optionally be a [api.ProviderMetadata].
func SystemMessage(args ...any) try.Try[api.SystemMessage] {
	var content string
	var metadata *api.ProviderMetadata

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			// TODO: should we concatenate multiple strings instead?
			// SystemMessage only supports a single string content
			if content != "" {
				return try.Errf[api.SystemMessage]("multiple string contents provided for SystemMessage")
			}
			content = v
		case *api.ProviderMetadata:
			if metadata != nil {
				return try.Errf[api.SystemMessage]("duplicate metadata provided: metadata can only be specified once")
			}
			metadata = v
		default:
			return try.Errf[api.SystemMessage]("unsupported argument type for SystemMessage: %T", arg)
		}
	}

	if content == "" {
		return try.Errf[api.SystemMessage]("no content provided for SystemMessage")
	}

	return try.Ok(api.SystemMessage{
		Content:          content,
		ProviderMetadata: metadata,
	})
}

// TODO: do we need a helper for the provider metadata? Or do we need to
// improve the api.NewProviderMetadata() constructor?
