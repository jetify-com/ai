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

	Betas    []anthropic.AnthropicBeta
	Warnings []api.CallWarning
}

func EncodeToolMode(mode api.ModeConfig) (AnthropicTools, error) {
	switch m := mode.(type) {
	case api.RegularMode:
		return encodeRegularMode(m)
	case *api.RegularMode:
		return encodeRegularMode(*m)

	case api.ObjectJSONMode:
		return AnthropicTools{}, fmt.Errorf("json-mode object generation is not supported")

	case *api.ObjectJSONMode:
		return AnthropicTools{}, fmt.Errorf("json-mode object generation is not supported")

	case api.ObjectToolMode:
		return encodeObjectToolMode(m)

	case *api.ObjectToolMode:
		return encodeObjectToolMode(*m)

	default:
		return AnthropicTools{}, fmt.Errorf("unsupported mode type: %T", mode)
	}
}

func encodeRegularMode(mode api.RegularMode) (AnthropicTools, error) {
	if len(mode.Tools) == 0 {
		return AnthropicTools{}, nil
	}

	tools := make([]anthropic.BetaToolUnionUnionParam, 0, len(mode.Tools))
	warnings := make([]api.CallWarning, 0)
	betasMap := make(map[anthropic.AnthropicBeta]bool)

	// Process each tool
	for _, toolItem := range mode.Tools {
		switch tool := toolItem.(type) {
		case api.FunctionTool:
			functionTool, err := EncodeFunctionTool(tool)
			if err != nil {
				return AnthropicTools{}, err
			}
			tools = append(tools, functionTool)
		case *api.FunctionTool:
			functionTool, err := EncodeFunctionTool(*tool)
			if err != nil {
				return AnthropicTools{}, err
			}
			tools = append(tools, functionTool)
		case api.ProviderDefinedTool:
			providerTool, toolBetas, toolWarnings, err := EncodeProviderDefinedTool(tool)
			if err != nil {
				return AnthropicTools{}, err
			}

			// Add tool if it's not nil (some tools might be unsupported)
			if providerTool != nil {
				tools = append(tools, providerTool)
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
	toolChoice, err := EncodeToolChoice(mode.ToolChoice)
	if err != nil {
		return AnthropicTools{}, err
	}

	// Special case for "none" tool choice
	if mode.ToolChoice != nil && mode.ToolChoice.Type == "none" {
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
		Tools:      tools,
		ToolChoice: toolChoice,
		Betas:      betas,
		Warnings:   warnings,
	}, nil
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

func encodeObjectToolMode(mode api.ObjectToolMode) (AnthropicTools, error) {
	inputSchema, err := encodeInputSchema(mode.Tool.InputSchema)
	if err != nil {
		return AnthropicTools{}, fmt.Errorf("error encoding input schema: %w", err)
	}

	tool := anthropic.BetaToolParam{
		Name:        anthropic.String(mode.Tool.Name),
		Description: anthropic.String(mode.Tool.Description),
		InputSchema: anthropic.F(inputSchema),
	}
	tools := []anthropic.BetaToolUnionUnionParam{tool}

	toolChoice := []anthropic.BetaToolChoiceUnionParam{
		anthropic.BetaToolChoiceToolParam{
			Type: anthropic.F(anthropic.BetaToolChoiceToolTypeTool),
			Name: anthropic.String(mode.Tool.Name),
		},
	}

	return AnthropicTools{
		Tools:      tools,
		ToolChoice: toolChoice,
		Betas:      []anthropic.AnthropicBeta{},
		Warnings:   []api.CallWarning{},
	}, nil
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
