package ai

import (
	"sync/atomic"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/anthropic"
)

var defaultLanguageModel atomic.Value

func init() {
	model := anthropic.NewLanguageModel(anthropic.ModelClaude37Sonnet20250219)
	defaultLanguageModel.Store(model)
}

func SetDefaultLanguageModel(lm api.LanguageModel) {
	defaultLanguageModel.Store(lm)
}

func DefaultLanguageModel() api.LanguageModel {
	return defaultLanguageModel.Load().(api.LanguageModel)
}
