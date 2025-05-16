package codec

import (
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
	"github.com/sashabaranov/go-openai/jsonschema"
	"go.jetify.com/ai/api"
)

type OpenAITools struct {
	ToolChoice     responses.ResponseNewParamsToolChoiceUnion
	Tools          []responses.ToolUnionParam
	Warnings       []api.CallWarning
	ResponseFormat *responses.ResponseTextConfigParam
}

// TODO: support provider metadata in tool mode information.
// getIsStrict determines if schema validation should be strict based on provider metadata
func getIsStrict(opts api.CallOptions) bool {
	isStrict := true // Default to true
	if opts.ProviderMetadata != nil {
		metadata := GetMetadata(opts)
		if metadata != nil && metadata.StrictSchemas != nil {
			isStrict = *metadata.StrictSchemas
		}
	}
	return isStrict
}

func EncodeToolMode(mode api.ModeConfig, opts api.CallOptions) (OpenAITools, error) {
	switch m := mode.(type) {
	case api.RegularMode:
		return EncodeRegularMode(m, opts)
	case *api.RegularMode:
		return EncodeRegularMode(*m, opts)

	case api.ObjectJSONMode:
		return encodeObjectJSONMode(m, opts)
	case *api.ObjectJSONMode:
		return encodeObjectJSONMode(*m, opts)

	case api.ObjectToolMode:
		return encodeObjectToolMode(m, opts)
	case *api.ObjectToolMode:
		return encodeObjectToolMode(*m, opts)

	default:
		return OpenAITools{}, fmt.Errorf("unsupported mode type: %T", mode)
	}
}

func EncodeRegularMode(mode api.RegularMode, opts api.CallOptions) (OpenAITools, error) {
	if len(mode.Tools) == 0 {
		return OpenAITools{}, nil
	}

	tools := make([]responses.ToolUnionParam, 0, len(mode.Tools))
	warnings := []api.CallWarning{}

	// Process each tool
	for _, toolItem := range mode.Tools {
		tool, toolWarnings, err := encodeToolDefinition(toolItem)
		if err != nil {
			return OpenAITools{}, err
		}

		if len(toolWarnings) > 0 {
			warnings = append(warnings, toolWarnings...)
		}

		if tool != nil {
			tools = append(tools, *tool)
		}
	}

	// Process tool choice
	toolChoice, err := encodeToolChoice(mode.ToolChoice)
	if err != nil {
		return OpenAITools{}, err
	}

	return OpenAITools{
		Tools:      tools,
		ToolChoice: toolChoice,
		Warnings:   warnings,
	}, nil
}

// encodeTool encodes a single tool into a ToolUnionParam
func encodeToolDefinition(toolItem api.ToolDefinition) (*responses.ToolUnionParam, []api.CallWarning, error) {
	switch tool := toolItem.(type) {
	case api.FunctionTool:
		return encodeFunctionTool(tool)
	case *api.FunctionTool:
		return encodeFunctionTool(*tool)
	case api.ProviderDefinedTool:
		return encodeProviderDefinedTool(tool)
	default:
		warning := api.CallWarning{
			Type: "unsupported-tool",
			Tool: toolItem,
		}
		return nil, []api.CallWarning{warning}, nil
	}
}

// encodeFunctionTool encodes a function tool
func encodeFunctionTool(tool api.FunctionTool) (*responses.ToolUnionParam, []api.CallWarning, error) {
	name := tool.Name

	props, err := jsonSchemaAsMap(tool.InputSchema)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert tool parameters: %w", err)
	}

	result := responses.ToolParamOfFunction(
		name,
		props,
		// TODO: allow passing the strict flag to the function
		true, // strict mode enabled
	)

	// Add description if provided
	if tool.Description != "" {
		if functionToolParam := result.OfFunction; functionToolParam != nil {
			functionToolParam.Description = openai.Opt(tool.Description)
		}
	}

	return &result, nil, nil
}

// encodeProviderDefinedTool encodes a provider-defined tool
func encodeProviderDefinedTool(tool api.ProviderDefinedTool) (*responses.ToolUnionParam, []api.CallWarning, error) {
	var result responses.ToolUnionParam
	var err error

	switch tool.ID() {
	case "openai.file_search":
		result, err = encodeFileSearchTool(tool)
		if err != nil {
			return nil, nil, err
		}
	case "openai.web_search_preview":
		result, err = encodeWebSearchTool(tool)
		if err != nil {
			return nil, nil, err
		}
	case "openai.computer_use_preview":
		result, err = encodeComputerUseTool(tool)
		if err != nil {
			return nil, nil, err
		}
	default:
		warning := api.CallWarning{
			Type: "unsupported-tool",
			Tool: tool,
		}
		return nil, []api.CallWarning{warning}, nil
	}

	return &result, nil, nil
}

// encodeFileSearchTool creates a file search tool parameter
func encodeFileSearchTool(tool api.ProviderDefinedTool) (responses.ToolUnionParam, error) {
	fileSearchTool, ok := tool.(*FileSearchTool)
	if !ok {
		return responses.ToolUnionParam{}, fmt.Errorf("expected FileSearchTool but got %T", tool)
	}

	return responses.ToolParamOfFileSearch(fileSearchTool.VectorStoreIDs), nil
}

// encodeWebSearchTool creates a web search tool parameter
func encodeWebSearchTool(tool api.ProviderDefinedTool) (responses.ToolUnionParam, error) {
	webSearchTool, ok := tool.(*WebSearchTool)
	if !ok {
		return responses.ToolUnionParam{}, fmt.Errorf("expected WebSearchTool but got %T", tool)
	}

	// Create a web search tool param directly instead of using the helper function
	var webSearchParam responses.WebSearchToolParam
	webSearchParam.Type = responses.WebSearchToolTypeWebSearchPreview

	// Set search context size if provided
	if webSearchTool.SearchContextSize != "" {
		webSearchParam.SearchContextSize = responses.WebSearchToolSearchContextSize(webSearchTool.SearchContextSize)
	}

	// Set user location if provided
	if webSearchTool.UserLocation != nil {
		userLocation := responses.WebSearchToolUserLocationParam{}

		if webSearchTool.UserLocation.City != "" {
			userLocation.City = openai.Opt(webSearchTool.UserLocation.City)
		}

		if webSearchTool.UserLocation.Country != "" {
			userLocation.Country = openai.Opt(webSearchTool.UserLocation.Country)
		}

		if webSearchTool.UserLocation.Region != "" {
			userLocation.Region = openai.Opt(webSearchTool.UserLocation.Region)
		}

		if webSearchTool.UserLocation.Timezone != "" {
			userLocation.Timezone = openai.Opt(webSearchTool.UserLocation.Timezone)
		}

		// Only set the UserLocation if at least one field was set
		if userLocation.City.IsPresent() || userLocation.Country.IsPresent() ||
			userLocation.Region.IsPresent() || userLocation.Timezone.IsPresent() {
			webSearchParam.UserLocation = userLocation
		}
	}

	return responses.ToolUnionParam{OfWebSearch: &webSearchParam}, nil
}

// encodeComputerUseTool creates a computer use tool parameter
func encodeComputerUseTool(tool api.ProviderDefinedTool) (responses.ToolUnionParam, error) {
	computerUseTool, ok := tool.(*ComputerUseTool)
	if !ok {
		return responses.ToolUnionParam{}, fmt.Errorf("expected ComputerUseTool but got %T", tool)
	}

	// Validate required parameters
	if computerUseTool.DisplayHeight <= 0 {
		return responses.ToolUnionParam{}, fmt.Errorf("displayHeight is required and must be positive")
	}

	if computerUseTool.DisplayWidth <= 0 {
		return responses.ToolUnionParam{}, fmt.Errorf("displayWidth is required and must be positive")
	}

	if computerUseTool.Environment == "" {
		return responses.ToolUnionParam{}, fmt.Errorf("environment is required")
	}

	// Validate that environment is one of the allowed values
	env := computerUseTool.Environment
	if env != "mac" && env != "windows" && env != "ubuntu" && env != "browser" {
		return responses.ToolUnionParam{}, fmt.Errorf("environment must be one of: mac, windows, ubuntu, browser; got %q", env)
	}

	return responses.ToolParamOfComputerUsePreview(
		computerUseTool.DisplayHeight,
		computerUseTool.DisplayWidth,
		responses.ComputerToolEnvironment(computerUseTool.Environment),
	), nil
}

// encodeToolChoice encodes a tool choice
func encodeToolChoice(toolChoice *api.ToolChoice) (responses.ResponseNewParamsToolChoiceUnion, error) {
	var result responses.ResponseNewParamsToolChoiceUnion

	if toolChoice == nil {
		return result, nil
	}

	switch toolChoice.Type {
	case "auto":
		result = responses.ResponseNewParamsToolChoiceUnion{
			OfToolChoiceMode: openai.Opt(responses.ToolChoiceOptionsAuto),
		}
	case "none":
		result = responses.ResponseNewParamsToolChoiceUnion{
			OfToolChoiceMode: openai.Opt(responses.ToolChoiceOptionsNone),
		}
	case "required":
		result = responses.ResponseNewParamsToolChoiceUnion{
			OfToolChoiceMode: openai.Opt(responses.ToolChoiceOptionsRequired),
		}
	case "tool":
		// Check if it's a provider-defined tool or a function tool
		switch toolChoice.ToolName {
		case "file_search", "web_search_preview", "computer_use_preview":
			// It's a provider-defined tool (hosted tool in OpenAI's terminology)
			result = responses.ResponseNewParamsToolChoiceUnion{
				OfHostedTool: &responses.ToolChoiceTypesParam{
					Type: responses.ToolChoiceTypesType(toolChoice.ToolName),
				},
			}
		default:
			//  It's a function tool
			result = responses.ResponseNewParamsToolChoiceUnion{
				OfFunctionTool: &responses.ToolChoiceFunctionParam{
					Name: toolChoice.ToolName,
				},
			}
		}
	default:
		return responses.ResponseNewParamsToolChoiceUnion{}, fmt.Errorf("unsupported tool choice type: %s", toolChoice.Type)
	}

	return result, nil
}

// TODO: promote to a framework-level function
func jsonSchemaAsMap(schema *jsonschema.Definition) (map[string]any, error) {
	if schema == nil {
		return nil, nil
	}

	// Marshal to JSON and unmarshal back to interface{} to convert the types
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal properties: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal properties: %w", err)
	}

	return result, nil
}

func encodeObjectToolMode(mode api.ObjectToolMode, opts api.CallOptions) (OpenAITools, error) {
	// Convert the schema to a map
	// TODO: the schema struct itself should have an .AsMap() method to convert to a map
	props, err := jsonSchemaAsMap(mode.Tool.InputSchema)
	if err != nil {
		return OpenAITools{}, fmt.Errorf("failed to convert tool parameters: %w", err)
	}

	// Get isStrict from opts
	isStrict := getIsStrict(opts)

	// Convert the tool to a function tool format
	tool := responses.ToolParamOfFunction(
		mode.Tool.Name,
		props,
		isStrict,
	)

	// Add description if provided
	if mode.Tool.Description != "" {
		if functionToolParam := tool.OfFunction; functionToolParam != nil {
			functionToolParam.Description = openai.Opt(mode.Tool.Description)
		}
	}

	tools := []responses.ToolUnionParam{tool}

	// Set tool choice to require this specific function
	toolChoice := responses.ResponseNewParamsToolChoiceUnion{
		OfFunctionTool: &responses.ToolChoiceFunctionParam{
			Name: mode.Tool.Name,
		},
	}

	return OpenAITools{
		Tools:      tools,
		ToolChoice: toolChoice,
		Warnings:   []api.CallWarning{},
	}, nil
}

func encodeObjectJSONMode(mode api.ObjectJSONMode, opts api.CallOptions) (OpenAITools, error) {
	isStrict := getIsStrict(opts)
	var responseFormat *responses.ResponseTextConfigParam

	if mode.Schema != nil {
		// Convert schema to map
		schemaMap, err := jsonSchemaAsMap(mode.Schema)
		if err != nil {
			return OpenAITools{}, fmt.Errorf("failed to convert JSON schema: %w", err)
		}

		// Create response format with JSON schema
		responseFormat = &responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Type:   "json_schema",
					Name:   mode.Name,
					Schema: schemaMap,
					Strict: openai.Bool(isStrict),
				},
			},
		}

		// Add description if provided
		if mode.Description != "" {
			responseFormat.Format.OfJSONSchema.Description = openai.String(mode.Description)
		}
	} else {
		// Generic JSON object format
		responseFormat = &responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONObject: &shared.ResponseFormatJSONObjectParam{
					Type: "json_object",
				},
			},
		}
	}

	return OpenAITools{
		ResponseFormat: responseFormat,
		Warnings:       []api.CallWarning{},
	}, nil
}
