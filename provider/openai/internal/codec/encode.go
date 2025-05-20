package codec

import (
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
	"go.jetify.com/ai/api"
)

func Encode(modelID string, prompt []api.Message, opts api.CallOptions) (responses.ResponseNewParams, []api.CallWarning, error) {
	modelConfig := getModelConfig(modelID)

	params := responses.ResponseNewParams{}
	params.Model = modelID

	warnings, err := applyCallOptions(&params, opts, modelConfig)
	if err != nil {
		return responses.ResponseNewParams{}, warnings, err
	}

	// Handle Mode configuration (which can override ResponseFormat)
	if opts.Mode != nil {
		tools, err := EncodeToolMode(opts.Mode, opts)
		if err != nil {
			return responses.ResponseNewParams{}, warnings, err
		}

		// Apply tools if any
		if len(tools.Tools) > 0 {
			params.Tools = tools.Tools
		}

		// Apply tool choice if set
		if !param.IsOmitted(tools.ToolChoice.OfToolChoiceMode) ||
			tools.ToolChoice.OfFunctionTool != nil ||
			tools.ToolChoice.OfHostedTool != nil {
			params.ToolChoice = tools.ToolChoice
		}

		// Apply response format if set (ObjectJSONMode)
		// This will override any ResponseFormat set by applyCallOptions
		if tools.ResponseFormat != nil {
			params.Text = *tools.ResponseFormat
		}

		warnings = append(warnings, tools.Warnings...)
	}

	// Encode the prompt to OpenAI format
	openaiPrompt, err := EncodePrompt(prompt, modelConfig)
	if err != nil {
		return responses.ResponseNewParams{}, warnings, err
	}
	// Set input messages
	params.Input = responses.ResponseNewParamsInputUnion{
		OfInputItemList: openaiPrompt.Messages,
	}

	warnings = append(warnings, openaiPrompt.Warnings...)

	return params, warnings, nil
}

func getModelConfig(modelID string) modelConfig {
	// o series reasoning models
	if len(modelID) > 0 && modelID[0] == 'o' {
		// Check for specific o1 models that need special handling
		if strings.HasPrefix(modelID, "o1-mini") || strings.HasPrefix(modelID, "o1-preview") {
			return modelConfig{
				IsReasoningModel:       true,
				SystemMessageMode:      "remove",
				RequiredAutoTruncation: false,
			}
		}

		// All other o-series models
		return modelConfig{
			IsReasoningModel:       true,
			SystemMessageMode:      "developer",
			RequiredAutoTruncation: false,
		}
	}

	// gpt models (non o-series)
	return modelConfig{
		IsReasoningModel:       false,
		SystemMessageMode:      "system",
		RequiredAutoTruncation: false,
	}
}

func applyCallOptions(params *responses.ResponseNewParams, opts api.CallOptions, modelConfig modelConfig) ([]api.CallWarning, error) {
	warnings := unsupportedWarnings(opts)

	// Set base parameters
	if opts.Temperature != nil {
		params.Temperature = openai.Float(*opts.Temperature)
	}
	if opts.TopP > 0 {
		params.TopP = openai.Float(opts.TopP)
	}
	if opts.MaxTokens > 0 {
		params.MaxOutputTokens = openai.Int(int64(opts.MaxTokens))
	}

	// Handle JSON response format if specified
	if opts.ResponseFormat != nil && opts.ResponseFormat.Type == "json" {
		err := applyJSONResponseFormat(params, opts)
		if err != nil {
			return warnings, err
		}
	}

	// Apply provider options from metadata
	reasoningEffort := applyProviderMetadata(params, opts)

	// Apply model-specific settings
	if modelConfig.RequiredAutoTruncation {
		params.Truncation = "auto"
	}

	// Apply reasoning settings and handle unsupported options for reasoning models
	reasoningWarnings := applyReasoningSettings(params, opts, modelConfig, reasoningEffort)
	warnings = append(warnings, reasoningWarnings...)

	return warnings, nil
}

// applyJSONResponseFormat handles setting up JSON response format options
func applyJSONResponseFormat(params *responses.ResponseNewParams, opts api.CallOptions) error {
	isStrict := getIsStrict(opts)

	if opts.ResponseFormat.Schema != nil {
		schemaMap, err := jsonSchemaAsMap(opts.ResponseFormat.Schema)
		if err != nil {
			return fmt.Errorf("failed to convert JSON schema: %w", err)
		}
		params.Text = responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Type:   "json_schema",
					Name:   opts.ResponseFormat.Name,
					Schema: schemaMap,
					Strict: openai.Bool(isStrict),
				},
			},
		}
		if opts.ResponseFormat.Description != "" {
			params.Text.Format.OfJSONSchema.Description = openai.String(opts.ResponseFormat.Description)
		}
	} else {
		params.Text = responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONObject: &shared.ResponseFormatJSONObjectParam{
					Type: "json_object",
				},
			},
		}
	}

	return nil
}

// applyProviderMetadata applies metadata-specific options to the parameters
// and returns the reasoning effort if specified
func applyProviderMetadata(params *responses.ResponseNewParams, opts api.CallOptions) string {
	var reasoningEffort string

	if opts.ProviderMetadata != nil {
		metadata := GetMetadata(opts)
		if metadata != nil {
			if metadata.ParallelToolCalls != nil {
				params.ParallelToolCalls = openai.Bool(*metadata.ParallelToolCalls)
			}
			if metadata.PreviousResponseID != "" {
				params.PreviousResponseID = openai.String(metadata.PreviousResponseID)
			}
			if metadata.Store != nil {
				params.Store = openai.Bool(*metadata.Store)
			}
			if metadata.User != "" {
				params.User = openai.String(metadata.User)
			}
			if metadata.Instructions != "" {
				params.Instructions = openai.String(metadata.Instructions)
			}
			// Extract reasoningEffort if available
			reasoningEffort = metadata.ReasoningEffort
		}
	}

	return reasoningEffort
}

// applyReasoningSettings applies settings specific to reasoning models
// and handles unsupported options
func applyReasoningSettings(params *responses.ResponseNewParams, opts api.CallOptions,
	modelConfig modelConfig, reasoningEffort string,
) []api.CallWarning {
	var warnings []api.CallWarning

	// Apply reasoning settings for reasoning models
	if modelConfig.IsReasoningModel && reasoningEffort != "" {
		params.Reasoning = shared.ReasoningParam{
			Effort: shared.ReasoningEffort(reasoningEffort),
		}
	}

	// Handle unsupported settings for reasoning models
	// See https://platform.openai.com/docs/guides/reasoning#limitations
	if modelConfig.IsReasoningModel {
		// Check if Temperature is set
		if opts.Temperature != nil {
			params.Temperature = param.Opt[float64]{} // Omit the field
			warnings = append(warnings, api.CallWarning{
				Type:    "unsupported-setting",
				Setting: "Temperature",
				Details: "Temperature is not supported for reasoning models",
			})
		}

		// Check if TopP is set
		if opts.TopP != 0 {
			params.TopP = param.Opt[float64]{} // Omit the field
			warnings = append(warnings, api.CallWarning{
				Type:    "unsupported-setting",
				Setting: "TopP",
				Details: "TopP is not supported for reasoning models",
			})
		}
	}

	return warnings
}

func unsupportedWarnings(opts api.CallOptions) []api.CallWarning {
	var warnings []api.CallWarning

	// Check for frequency penalty
	if opts.FrequencyPenalty != 0 {
		warnings = append(warnings, api.CallWarning{
			Type:    "unsupported-setting",
			Setting: "FrequencyPenalty",
		})
	}

	// Check for presence penalty
	if opts.PresencePenalty != 0 {
		warnings = append(warnings, api.CallWarning{
			Type:    "unsupported-setting",
			Setting: "PresencePenalty",
		})
	}

	// Check for top-k (not directly supported by OpenAI)
	if opts.TopK > 0 {
		warnings = append(warnings, api.CallWarning{
			Type:    "unsupported-setting",
			Setting: "TopK",
		})
	}

	// Check for seed
	if opts.Seed != 0 {
		warnings = append(warnings, api.CallWarning{
			Type:    "unsupported-setting",
			Setting: "Seed",
		})
	}

	// Check for stop sequences
	if len(opts.StopSequences) > 0 {
		warnings = append(warnings, api.CallWarning{
			Type:    "unsupported-setting",
			Setting: "StopSequences",
		})
	}

	return warnings
}
