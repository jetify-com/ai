package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/provider/anthropic"
	"go.jetify.com/ai/provider/openai"
)

func TestDefaultLanguageModel(t *testing.T) {
	// Get current model and verify provider and model ID match expected values
	originalModel := DefaultLanguageModel()
	assert.Equal(t, "openai", originalModel.ProviderName())
	assert.Equal(t, openai.ChatModelGPT5, originalModel.ModelID())

	// Change model to different provider (Anthropic)
	anthropicModel := anthropic.NewLanguageModel(anthropic.ModelClaude3_5SonnetLatest)
	SetDefaultLanguageModel(anthropicModel)

	currentModel := DefaultLanguageModel()
	assert.Equal(t, "anthropic", currentModel.ProviderName())
	assert.Equal(t, anthropic.ModelClaude3_5SonnetLatest, currentModel.ModelID())

	// Restore model to original
	SetDefaultLanguageModel(originalModel)

	restoredModel := DefaultLanguageModel()
	assert.Equal(t, "openai", restoredModel.ProviderName())
	assert.Equal(t, openai.ChatModelGPT5, restoredModel.ModelID())
}
