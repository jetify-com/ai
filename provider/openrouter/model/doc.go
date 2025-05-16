// Package model provides constants for specifying which model to use when using
// OpenRouter.
// OpenRouter is a unified API that provides access to various AI models from different
// providers including OpenAI, Anthropic, Google, Meta, and others. This package defines
// the model identifiers needed to specify which model to use when
// making requests through OpenRouter.
//
// Model identifiers follow the format "provider/model-name[:tag]", for example:
//   - openai/gpt-4
//   - anthropic/claude-3-opus
//   - meta-llama/llama-3-70b-instruct
//
// Some models have additional tags like ":free" or ":beta" that indicate special
// versions or pricing tiers.
//
// Example usage:
//
//	import "github.com/your-org/aisdk/provider/openrouter/model"
//
//	// Use a predefined model constant
//	modelID := model.O3MiniHigh
//
// For the most up-to-date list of available models and their capabilities,
// see https://openrouter.ai/docs#models
package model
