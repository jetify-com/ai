package codec

import (
	"github.com/anthropics/anthropic-sdk-go"
	"go.jetify.com/ai/api"
)

// For now we are using a single type for all metadata.
// TODO: Decide if we will need different types for different metadata.
type Metadata struct {
	// --- Used in requests ---
	CacheControl string `json:"cache_control,omitempty"`

	// --- Used in responses ---

	Thinking ThinkingConfig `json:"thinking,omitzero"`
	Usage    Usage          `json:"usage,omitempty"`
}

func GetMetadata(source api.MetadataSource) *Metadata {
	return api.GetMetadata[Metadata]("anthropic", source)
}

// Anthropic-specific Usage information like CacheCreationInputTokens
// and CacheReadInputTokens.
type Usage anthropic.BetaUsage

// ThinkingConfig represents the configuration for thinking behavior
type ThinkingConfig struct {
	// Whether to enable extended thinking.
	//
	// See
	// [extended thinking](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
	// for details.
	Enabled bool `json:"enabled,omitzero"`

	// Determines how many tokens Claude can use for its internal reasoning process.
	// Larger budgets can enable more thorough analysis for complex problems, improving
	// response quality.
	//
	// Must be ≥1024 and less than `max_tokens`.
	BudgetTokens int `json:"budgetTokens,omitzero"`
}
