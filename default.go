package ai

import (
	"sync/atomic"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/openai"
)

// modelWrapper wraps api.LanguageModel to ensure consistent type for atomic.Value
type modelWrapper struct {
	model api.LanguageModel
}

var defaultLanguageModel atomic.Value

func init() {
	model := openai.NewLanguageModel(openai.ChatModelGPT5)
	defaultLanguageModel.Store(&modelWrapper{model: model})
}

func SetDefaultLanguageModel(lm api.LanguageModel) {
	defaultLanguageModel.Store(&modelWrapper{model: lm})
}

func DefaultLanguageModel() api.LanguageModel {
	wrapper := defaultLanguageModel.Load().(*modelWrapper)
	return wrapper.model
}
