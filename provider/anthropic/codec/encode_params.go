package codec

import (
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"go.jetify.com/ai/api"
)

func EncodeParams(
	prompt []api.Message, opts api.CallOptions,
) (anthropic.BetaMessageNewParams, []api.CallWarning, error) {
	anthropicPrompt, err := EncodePrompt(prompt)
	if err != nil {
		return anthropic.BetaMessageNewParams{}, []api.CallWarning{}, err
	}

	params, warnings, err := encodeCallOptions(opts)
	if err != nil {
		return anthropic.BetaMessageNewParams{}, warnings, err
	}

	if len(anthropicPrompt.System) > 0 {
		params.System = anthropic.F(anthropicPrompt.System)
	}
	if len(anthropicPrompt.Messages) > 0 {
		params.Messages = anthropic.F(anthropicPrompt.Messages)
	}

	params.Betas = anthropic.F(append(params.Betas.Value, anthropicPrompt.Betas...))

	return params, warnings, nil
}

func encodeCallOptions(opts api.CallOptions) (anthropic.BetaMessageNewParams, []api.CallWarning, error) {
	params := anthropic.BetaMessageNewParams{
		MaxTokens: anthropic.F(int64(4096)), // Default max tokens
	}

	// Set basic parameters
	if opts.MaxOutputTokens > 0 {
		params.MaxTokens = anthropic.F(int64(opts.MaxOutputTokens))
	}
	if opts.Temperature != nil {
		params.Temperature = anthropic.F(*opts.Temperature)
	}
	if opts.TopP > 0 {
		params.TopP = anthropic.F(opts.TopP)
	}
	if opts.TopK > 0 {
		params.TopK = anthropic.F(int64(opts.TopK))
	}
	if len(opts.StopSequences) > 0 {
		params.StopSequences = anthropic.F(opts.StopSequences)
	}

	// Handle unsupported settings
	warnings := unsupportedWarnings(opts)

	// Handle thinking-specific configuration
	thinkingWarnings, err := encodeThinking(&params, opts)
	if err != nil {
		return params, warnings, err
	}
	warnings = append(warnings, thinkingWarnings...)

	// Handle tool configuration
	if opts.Mode != nil {
		tools, err := EncodeToolMode(opts.Mode)
		if err != nil {
			return params, warnings, err
		}
		params.Betas = anthropic.F(append(params.Betas.Value, tools.Betas...))
		warnings = append(warnings, tools.Warnings...)

		if len(tools.Tools) > 0 {
			params.Tools = anthropic.F(tools.Tools)
		}
		if len(tools.ToolChoice) > 0 {
			params.ToolChoice = anthropic.F(tools.ToolChoice[0])
		}
	}
	return params, warnings, nil
}

func unsupportedWarnings(opts api.CallOptions) []api.CallWarning {
	var warnings []api.CallWarning

	if opts.FrequencyPenalty != 0 {
		warnings = append(warnings, api.CallWarning{
			Type:    "unsupported-setting",
			Setting: "FrequencyPenalty",
		})
	}

	if opts.PresencePenalty != 0 {
		warnings = append(warnings, api.CallWarning{
			Type:    "unsupported-setting",
			Setting: "PresencePenalty",
		})
	}

	if opts.Seed != 0 {
		warnings = append(warnings, api.CallWarning{
			Type:    "unsupported-setting",
			Setting: "Seed",
		})
	}

	if opts.ResponseFormat != nil && opts.ResponseFormat.Type != "text" {
		warnings = append(warnings, api.CallWarning{
			Type:    "unsupported-setting",
			Setting: "ResponseFormat",
			Details: "JSON response format is not supported.",
		})
	}

	return warnings
}

func encodeThinking(params *anthropic.BetaMessageNewParams, opts api.CallOptions) ([]api.CallWarning, error) {
	var warnings []api.CallWarning

	metadata := GetMetadata(opts)
	thinkingEnabled := metadata != nil && metadata.Thinking.Enabled

	if !thinkingEnabled {
		return warnings, nil
	}

	if metadata.Thinking.BudgetTokens == 0 {
		return warnings, fmt.Errorf("thinking requires a budget")
	}

	// Configure thinking parameters
	params.Thinking = anthropic.F[anthropic.BetaThinkingConfigParamUnion](
		anthropic.BetaThinkingConfigEnabledParam{
			Type:         anthropic.F(anthropic.BetaThinkingConfigEnabledTypeEnabled),
			BudgetTokens: anthropic.F(int64(metadata.Thinking.BudgetTokens)),
		})

	// Adjust max tokens to account for thinking budget
	params.MaxTokens = anthropic.F(params.MaxTokens.Value + int64(metadata.Thinking.BudgetTokens))

	// Add warnings for unsupported settings when thinking is enabled
	if opts.Temperature != nil {
		warnings = append(warnings, api.CallWarning{
			Type:    "unsupported-setting",
			Setting: "Temperature",
			Details: "Temperature is not supported when thinking is enabled",
		})
	}

	if opts.TopK > 0 {
		warnings = append(warnings, api.CallWarning{
			Type:    "unsupported-setting",
			Setting: "TopK",
			Details: "TopK is not supported when thinking is enabled",
		})
	}

	if opts.TopP > 0 {
		warnings = append(warnings, api.CallWarning{
			Type:    "unsupported-setting",
			Setting: "TopP",
			Details: "TopP is not supported when thinking is enabled",
		})
	}

	return warnings, nil
}
