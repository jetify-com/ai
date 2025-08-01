package codec

import (
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/sashabaranov/go-openai/jsonschema"
	"go.jetify.com/ai/api"
)

type AnthropicTools struct {
	Tools      []anthropic.BetaToolUnionParam
	ToolChoice anthropic.BetaToolChoiceUnionParam
	Betas      []anthropic.AnthropicBeta
	Warnings   []api.CallWarning
}

// convertArgs is a generic helper function that converts tool.Args to the expected struct type.
// It handles both cases where Args is already the correct struct type or a map[string]any.
func convertArgs[T any](args any) (*T, error) {
	// Try direct type assertion first
	if typedArgs, ok := args.(*T); ok {
		return typedArgs, nil
	}
	if typedArgs, ok := args.(T); ok {
		return &typedArgs, nil
	}

	// If that fails, try converting from map[string]any via JSON marshaling/unmarshaling
	data, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal args: %w", err)
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal args to target type: %w", err)
	}

	return &result, nil
}

// encodeInputSchema converts the JSON schema definition to Anthropic's schema format
func encodeInputSchema(schema *jsonschema.Definition) (anthropic.BetaToolInputSchemaParam, error) {
	// Verify the schema type is "object"
	if schema.Type != "" && schema.Type != "object" {
		return anthropic.BetaToolInputSchemaParam{}, fmt.Errorf("unsupported schema type: %s, only 'object' is supported", schema.Type)
	}

	// Create the input schema with the type field
	inputSchema := anthropic.BetaToolInputSchemaParam{
		Type: "object",
	}

	// Add properties only if they exist
	if len(schema.Properties) > 0 {
		// Convert the properties to a map[string]any by marshaling and unmarshaling
		var properties map[string]any
		// Marshal to JSON
		propsJSON, err := json.Marshal(schema.Properties)
		if err != nil {
			return anthropic.BetaToolInputSchemaParam{}, fmt.Errorf("failed to marshal properties: %w", err)
		}

		// Unmarshal back to map[string]any
		if err := json.Unmarshal(propsJSON, &properties); err != nil {
			return anthropic.BetaToolInputSchemaParam{}, fmt.Errorf("failed to unmarshal properties: %w", err)
		}

		// Set the properties field
		inputSchema.Properties = properties
	}

	// Add the required field if it's present in the original schema
	if len(schema.Required) > 0 {
		inputSchema.Required = schema.Required
	}

	return inputSchema, nil
}

// EncodeFunctionTool converts an API FunctionTool to Anthropic's tool format
func EncodeFunctionTool(tool api.FunctionTool) (anthropic.BetaToolUnionParam, error) {
	inputSchema, err := encodeInputSchema(tool.InputSchema)
	if err != nil {
		return anthropic.BetaToolUnionParam{}, fmt.Errorf("error encoding input schema: %w", err)
	}

	return anthropic.BetaToolUnionParam{
		OfTool: &anthropic.BetaToolParam{
			Type:        "tool",
			Name:        tool.Name,
			Description: anthropic.String(tool.Description),
			InputSchema: inputSchema,
		},
	}, nil
}

// EncodeProviderDefinedTool converts provider-defined tools to Anthropic's format
// Returns the tool, required betas, and any warnings
func EncodeProviderDefinedTool(
	tool api.ProviderDefinedTool,
) (anthropic.BetaToolUnionParam, []anthropic.AnthropicBeta, []api.CallWarning, error) {
	var warnings []api.CallWarning
	var betas []anthropic.AnthropicBeta

	// TODO: Instead of specifying the tool version as part of the tool definition,
	// should we automatically figure it out based on the on the version of Sonnet selected?

	switch tool.ID {
	case "anthropic.computer":
		computerArgs, err := convertArgs[ComputerToolArgs](tool.Args)
		if err != nil {
			return anthropic.BetaToolUnionParam{}, betas, warnings, fmt.Errorf("failed to convert computer tool args: %w", err)
		}
		switch computerArgs.Version {
		case "", "20250124": // default to 20250124
			betas = append(betas, anthropic.AnthropicBetaComputerUse2025_01_24)
			toolParam := anthropic.BetaToolUnionParam{
				OfComputerUseTool20250124: &anthropic.BetaToolComputerUse20250124Param{
					Type:            "computer_20250124",
					DisplayHeightPx: int64(computerArgs.DisplayHeightPx),
					DisplayWidthPx:  int64(computerArgs.DisplayWidthPx),
				},
			}
			if computerArgs.DisplayNumber != 0 {
				toolParam.OfComputerUseTool20250124.DisplayNumber = anthropic.Int(int64(computerArgs.DisplayNumber))
			}
			return toolParam, betas, warnings, nil

		case "20241022":
			betas = append(betas, anthropic.AnthropicBetaComputerUse2024_10_22)
			toolParam := anthropic.BetaToolUnionParam{
				OfComputerUseTool20241022: &anthropic.BetaToolComputerUse20241022Param{
					Type:            "computer_20241022",
					DisplayHeightPx: int64(computerArgs.DisplayHeightPx),
					DisplayWidthPx:  int64(computerArgs.DisplayWidthPx),
				},
			}
			if computerArgs.DisplayNumber != 0 {
				toolParam.OfComputerUseTool20241022.DisplayNumber = anthropic.Int(int64(computerArgs.DisplayNumber))
			}
			return toolParam, betas, warnings, nil
		default:
			return anthropic.BetaToolUnionParam{}, betas, warnings, fmt.Errorf("unsupported computer tool version: %s", computerArgs.Version)
		}

	case "anthropic.text_editor":
		textEditorArgs, err := convertArgs[TextEditorToolArgs](tool.Args)
		if err != nil {
			return anthropic.BetaToolUnionParam{}, betas, warnings, fmt.Errorf("failed to convert text editor tool args: %w", err)
		}
		switch textEditorArgs.Version {
		case "", "20250124": // default to 20250124
			betas = append(betas, anthropic.AnthropicBetaComputerUse2025_01_24)
			return anthropic.BetaToolUnionParam{
				OfTextEditor20250124: &anthropic.BetaToolTextEditor20250124Param{
					Type: "text_editor_20250124",
				},
			}, betas, warnings, nil

		case "20241022":
			betas = append(betas, anthropic.AnthropicBetaComputerUse2024_10_22)
			return anthropic.BetaToolUnionParam{
				OfTextEditor20241022: &anthropic.BetaToolTextEditor20241022Param{
					Type: "text_editor_20241022",
				},
			}, betas, warnings, nil
		default:
			return anthropic.BetaToolUnionParam{}, betas, warnings, fmt.Errorf("unsupported text editor tool version: %s", textEditorArgs.Version)
		}

	case "anthropic.bash":
		bashArgs, err := convertArgs[BashToolArgs](tool.Args)
		if err != nil {
			return anthropic.BetaToolUnionParam{}, betas, warnings, fmt.Errorf("failed to convert bash tool args: %w", err)
		}
		switch bashArgs.Version {
		case "", "20250124": // default to 20250124
			betas = append(betas, anthropic.AnthropicBetaComputerUse2025_01_24)
			return anthropic.BetaToolUnionParam{
				OfBashTool20250124: &anthropic.BetaToolBash20250124Param{
					Type: "bash_20250124",
				},
			}, betas, warnings, nil

		case "20241022":
			betas = append(betas, anthropic.AnthropicBetaComputerUse2024_10_22)
			return anthropic.BetaToolUnionParam{
				OfBashTool20241022: &anthropic.BetaToolBash20241022Param{
					Type: "bash_20241022",
				},
			}, betas, warnings, nil
		default:
			return anthropic.BetaToolUnionParam{}, betas, warnings, fmt.Errorf("unsupported bash tool version: %s", bashArgs.Version)
		}

	default:
		warnings = append(warnings, api.CallWarning{
			Type: "unsupported-tool",
			Tool: tool,
		})
		return anthropic.BetaToolUnionParam{}, betas, warnings, nil
	}
}

// EncodeToolChoice converts API ToolChoice to Anthropic's format
func EncodeToolChoice(
	toolChoice *api.ToolChoice,
) (anthropic.BetaToolChoiceUnionParam, error) {
	if toolChoice == nil {
		return anthropic.BetaToolChoiceUnionParam{}, nil
	}

	switch toolChoice.Type {
	case "auto":
		return anthropic.BetaToolChoiceUnionParam{
			OfAuto: &anthropic.BetaToolChoiceAutoParam{
				Type: "auto",
			},
		}, nil
	case "required":
		return anthropic.BetaToolChoiceUnionParam{
			OfAny: &anthropic.BetaToolChoiceAnyParam{
				Type: "any",
			},
		}, nil
	case "none":
		// Handled in EncodeTools
		return anthropic.BetaToolChoiceUnionParam{}, nil
	case "tool":
		// Use the constructor function
		return anthropic.BetaToolChoiceParamOfTool(toolChoice.ToolName), nil
	default:
		return anthropic.BetaToolChoiceUnionParam{}, fmt.Errorf("unsupported tool choice type: %s", toolChoice.Type)
	}
}

// EncodeTools processes the top-level Tools and ToolChoice fields from CallOptions
// and returns all the anthropic-specific tool configuration
func EncodeTools(tools []api.ToolDefinition, toolChoice *api.ToolChoice) (AnthropicTools, error) {
	var warnings []api.CallWarning

	// Handle case where no tools are provided
	if len(tools) == 0 && toolChoice == nil {
		return AnthropicTools{}, nil
	}

	anthropicTools := make([]anthropic.BetaToolUnionParam, 0, len(tools))
	betasMap := make(map[anthropic.AnthropicBeta]bool)

	// Process each tool
	for _, toolItem := range tools {
		switch tool := toolItem.(type) {
		case api.FunctionTool:
			functionTool, err := EncodeFunctionTool(tool)
			if err != nil {
				return AnthropicTools{}, err
			}
			anthropicTools = append(anthropicTools, functionTool)
		case *api.FunctionTool:
			functionTool, err := EncodeFunctionTool(*tool)
			if err != nil {
				return AnthropicTools{}, err
			}
			anthropicTools = append(anthropicTools, functionTool)
		case api.ProviderDefinedTool:
			providerTool, toolBetas, toolWarnings, err := EncodeProviderDefinedTool(tool)
			if err != nil {
				return AnthropicTools{}, err
			}

			// Add tool if it's not empty (some tools might be unsupported)
			if providerTool.GetType() != nil {
				anthropicTools = append(anthropicTools, providerTool)
			}

			// Add betas to the map
			for _, beta := range toolBetas {
				betasMap[beta] = true
			}

			// Add warnings
			warnings = append(warnings, toolWarnings...)
		default:
			warnings = append(warnings, api.CallWarning{
				Type: "unsupported-tool",
				Tool: toolItem,
			})
		}
	}

	// Create beta slice from map keys
	betas := make([]anthropic.AnthropicBeta, 0, len(betasMap))
	for beta := range betasMap {
		betas = append(betas, beta)
	}

	// Process tool choice
	anthropicToolChoice, err := EncodeToolChoice(toolChoice)
	if err != nil {
		return AnthropicTools{}, err
	}

	// Special case for "none" tool choice
	if toolChoice != nil && toolChoice.Type == "none" {
		// Anthropic does not support 'none' tool choice, so we remove the tools
		// but still return warnings and betas
		return AnthropicTools{
			Tools:      nil,
			ToolChoice: anthropic.BetaToolChoiceUnionParam{},
			Betas:      betas,
			Warnings:   warnings,
		}, nil
	}

	return AnthropicTools{
		Tools:      anthropicTools,
		ToolChoice: anthropicToolChoice,
		Betas:      betas,
		Warnings:   warnings,
	}, nil
}
