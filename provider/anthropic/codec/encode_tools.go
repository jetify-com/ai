package codec

import (
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/sashabaranov/go-openai/jsonschema"
	"go.jetify.com/ai/api"
)

type AnthropicTools struct {
	Tools      []anthropic.BetaToolUnionUnionParam
	ToolChoice []anthropic.BetaToolChoiceUnionParam
	Betas      []anthropic.AnthropicBeta
	Warnings   []api.CallWarning
}

// encodeInputSchema converts the JSON schema definition to Anthropic's schema format
func encodeInputSchema(schema *jsonschema.Definition) (anthropic.BetaToolInputSchemaParam, error) {
	// Verify the schema type is "object"
	if schema.Type != "" && schema.Type != "object" {
		return anthropic.BetaToolInputSchemaParam{}, fmt.Errorf("unsupported schema type: %s, only 'object' is supported", schema.Type)
	}

	// Create the input schema with the type field
	inputSchema := anthropic.BetaToolInputSchemaParam{
		Type: anthropic.F(anthropic.BetaToolInputSchemaTypeObject),
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
		inputSchema.Properties = anthropic.F[interface{}](properties)
	}

	// Add the required field if it's present in the original schema
	if len(schema.Required) > 0 {
		inputSchema.ExtraFields = map[string]interface{}{
			"required": schema.Required,
		}
	}

	return inputSchema, nil
}

// EncodeFunctionTool converts an API FunctionTool to Anthropic's tool format
func EncodeFunctionTool(tool api.FunctionTool) (anthropic.BetaToolUnionUnionParam, error) {
	inputSchema, err := encodeInputSchema(tool.InputSchema)
	if err != nil {
		return nil, fmt.Errorf("error encoding input schema: %w", err)
	}

	return anthropic.BetaToolParam{
		Name:        anthropic.String(tool.Name),
		Description: anthropic.String(tool.Description),
		InputSchema: anthropic.F(inputSchema),
	}, nil
}

// EncodeProviderDefinedTool converts provider-defined tools to Anthropic's format
// Returns the tool, required betas, and any warnings
func EncodeProviderDefinedTool(
	tool api.ProviderDefinedTool,
) (anthropic.BetaToolUnionUnionParam, []anthropic.AnthropicBeta, []api.CallWarning, error) {
	var warnings []api.CallWarning
	var betas []anthropic.AnthropicBeta

	// TODO: Instead of specifying the tool version as part of the tool definition,
	// should we automatically figure it out based on the on the version of Sonnet selected?

	switch tool.ID() {
	case "anthropic.computer":
		computerTool, ok := tool.(*ComputerUseTool)
		if !ok {
			return nil, betas, warnings, fmt.Errorf("computer tool must be of type ComputerUseTool, got %T", tool)
		}
		switch computerTool.Version {
		case "", "20250124": // default to 20250124
			betas = append(betas, anthropic.AnthropicBetaComputerUse2025_01_24)
			return anthropic.BetaToolComputerUse20250124Param{
				Name:            anthropic.F(anthropic.BetaToolComputerUse20250124Name(computerTool.Name())),
				Type:            anthropic.F(anthropic.BetaToolComputerUse20250124TypeComputer20250124),
				DisplayWidthPx:  anthropic.Int(int64(computerTool.DisplayWidthPx)),
				DisplayHeightPx: anthropic.Int(int64(computerTool.DisplayHeightPx)),
				DisplayNumber:   anthropic.Int(int64(computerTool.DisplayNumber)),
			}, betas, warnings, nil

		case "20241022":
			betas = append(betas, anthropic.AnthropicBetaComputerUse2024_10_22)
			return anthropic.BetaToolComputerUse20241022Param{
				Name:            anthropic.F(anthropic.BetaToolComputerUse20241022Name(computerTool.Name())),
				Type:            anthropic.F(anthropic.BetaToolComputerUse20241022TypeComputer20241022),
				DisplayWidthPx:  anthropic.Int(int64(computerTool.DisplayWidthPx)),
				DisplayHeightPx: anthropic.Int(int64(computerTool.DisplayHeightPx)),
				DisplayNumber:   anthropic.Int(int64(computerTool.DisplayNumber)),
			}, betas, warnings, nil
		default:
			return nil, betas, warnings, fmt.Errorf("unsupported computer tool version: %s", computerTool.Version)
		}

	case "anthropic.text_editor":
		textEditorTool, ok := tool.(*TextEditorTool)
		if !ok {
			return nil, betas, warnings, fmt.Errorf("text editor tool must be of type TextEditorTool, got %T", tool)
		}
		switch textEditorTool.Version {
		case "", "20250124": // default to 20250124
			betas = append(betas, anthropic.AnthropicBetaComputerUse2025_01_24)
			return anthropic.BetaToolTextEditor20250124Param{
				Name: anthropic.F(anthropic.BetaToolTextEditor20250124Name(textEditorTool.Name())),
				Type: anthropic.F(anthropic.BetaToolTextEditor20250124TypeTextEditor20250124),
			}, betas, warnings, nil

		case "20241022":
			betas = append(betas, anthropic.AnthropicBetaComputerUse2024_10_22)
			return anthropic.BetaToolTextEditor20241022Param{
				Name: anthropic.F(anthropic.BetaToolTextEditor20241022Name(textEditorTool.Name())),
				Type: anthropic.F(anthropic.BetaToolTextEditor20241022TypeTextEditor20241022),
			}, betas, warnings, nil
		default:
			return nil, betas, warnings, fmt.Errorf("unsupported text editor tool version: %s", textEditorTool.Version)
		}

	case "anthropic.bash":
		bashTool, ok := tool.(*BashTool)
		if !ok {
			return nil, betas, warnings, fmt.Errorf("bash tool must be of type BashTool, got %T", tool)
		}
		switch bashTool.Version {
		case "", "20250124": // default to 20250124
			betas = append(betas, anthropic.AnthropicBetaComputerUse2025_01_24)
			return anthropic.BetaToolBash20250124Param{
				Name: anthropic.F(anthropic.BetaToolBash20250124Name(bashTool.Name())),
				Type: anthropic.F(anthropic.BetaToolBash20250124TypeBash20250124),
			}, betas, warnings, nil

		case "20241022":
			betas = append(betas, anthropic.AnthropicBetaComputerUse2024_10_22)
			return anthropic.BetaToolBash20241022Param{
				Name: anthropic.F(anthropic.BetaToolBash20241022Name(bashTool.Name())),
				Type: anthropic.F(anthropic.BetaToolBash20241022TypeBash20241022),
			}, betas, warnings, nil
		default:
			return nil, betas, warnings, fmt.Errorf("unsupported bash tool version: %s", bashTool.Version)
		}

	default:
		warnings = append(warnings, api.CallWarning{
			Type: "unsupported-tool",
			Tool: tool,
		})
		return nil, betas, warnings, nil
	}
}

// EncodeToolChoice converts API ToolChoice to Anthropic's format
func EncodeToolChoice(
	toolChoice *api.ToolChoice,
) ([]anthropic.BetaToolChoiceUnionParam, error) {
	if toolChoice == nil {
		return nil, nil
	}

	switch toolChoice.Type {
	case "auto":
		return []anthropic.BetaToolChoiceUnionParam{
			anthropic.BetaToolChoiceAutoParam{
				Type: anthropic.F(anthropic.BetaToolChoiceAutoTypeAuto),
			},
		}, nil
	case "required":
		return []anthropic.BetaToolChoiceUnionParam{
			anthropic.BetaToolChoiceAnyParam{
				Type: anthropic.F(anthropic.BetaToolChoiceAnyTypeAny),
			},
		}, nil
	case "none":
		// Handled in EncodeTools
		return nil, nil
	case "tool":
		return []anthropic.BetaToolChoiceUnionParam{
			anthropic.BetaToolChoiceToolParam{
				Type: anthropic.F(anthropic.BetaToolChoiceToolTypeTool),
				Name: anthropic.String(toolChoice.ToolName),
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported tool choice type: %s", toolChoice.Type)
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

	anthropicTools := make([]anthropic.BetaToolUnionUnionParam, 0, len(tools))
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

			// Add tool if it's not nil (some tools might be unsupported)
			if providerTool != nil {
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
			ToolChoice: nil,
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
