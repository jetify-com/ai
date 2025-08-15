package codec

import (
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
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
		metadata := GetMetadata(&opts)
		if metadata != nil && metadata.StrictSchemas != nil {
			isStrict = *metadata.StrictSchemas
		}
	}
	return isStrict
}

func EncodeTools(tools []api.ToolDefinition, toolChoice *api.ToolChoice, opts api.CallOptions) (OpenAITools, error) {
	if len(tools) == 0 && toolChoice == nil {
		return OpenAITools{}, nil
	}

	result := OpenAITools{
		Tools:    make([]responses.ToolUnionParam, 0, len(tools)),
		Warnings: []api.CallWarning{},
	}

	// Process each tool
	for _, toolItem := range tools {
		tool, toolWarnings, err := encodeToolDefinition(toolItem)
		if err != nil {
			return OpenAITools{}, err
		}

		if len(toolWarnings) > 0 {
			result.Warnings = append(result.Warnings, toolWarnings...)
		}

		if tool != nil {
			result.Tools = append(result.Tools, *tool)
		}
	}

	// Process tool choice
	if toolChoice != nil {
		choice, err := encodeToolChoice(toolChoice)
		if err != nil {
			return OpenAITools{}, err
		}
		result.ToolChoice = choice
	}

	return result, nil
}

// encodeTool encodes a single tool into a ToolUnionParam
func encodeToolDefinition(toolItem api.ToolDefinition) (*responses.ToolUnionParam, []api.CallWarning, error) {
	switch tool := toolItem.(type) {
	case *api.FunctionTool:
		return encodeFunctionTool(*tool)
	case *api.ProviderDefinedTool:
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

	props, err := encodeSchema(tool.InputSchema)
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
			functionToolParam.Description = openai.String(tool.Description)
		}
	}

	return &result, nil, nil
}

// convertArgs is a generic helper function that converts tool.Args to the expected struct type.
// It handles both cases where Args is already the correct struct type or a map[string]any.
func convertArgs[T any](args any) (*T, error) {
	// Check if args is nil
	if args == nil {
		return nil, fmt.Errorf("args cannot be nil")
	}
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

// encodeProviderDefinedTool encodes a provider-defined tool
func encodeProviderDefinedTool(tool *api.ProviderDefinedTool) (*responses.ToolUnionParam, []api.CallWarning, error) {
	var result responses.ToolUnionParam
	var err error

	switch tool.ID {
	case "openai.file_search":
		result, err = encodeFileSearchTool(tool)
		if err != nil {
			return nil, nil, err
		}
	case "openai.web_search_preview":
		// TODO: Decide how to evolve handling of the web search tool.
		// So far, there are three types of tools:
		// 1. User provided tools (function calls)
		// 2. Built in or provider defined tools that require the client to perform the action
		// 3. Built in or provider defined tools where the LLM can perform the action by itself.
		//
		// Web search is an example where the LLM already performs the action automatically,
		// and it's already returning a list of sources, along with the text that it generated.
		//
		// We have a few options:
		// 1. Ignore it. The tool has already been executed. The text already contains the sources.
		//This is what Vercel does.
		// 2. Include it. But maybe we should mark it as already executed somehow so users can distinguish.
		// 3. Instead of including it as a ToolCall, include it as a ToolResult. Normally, ToolResults are
		//    sent by the client as part of the prompt, letting the LLM know that the tool it requested has
		//    been executed. But it might be OK, to allow the LLM to return a ToolResult as part of the response
		//    in cases when it executes a tool by itself.

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
func encodeFileSearchTool(tool *api.ProviderDefinedTool) (responses.ToolUnionParam, error) {
	fileSearchArgs, err := convertArgs[FileSearchToolArgs](tool.Args)
	if err != nil {
		return responses.ToolUnionParam{}, fmt.Errorf("failed to convert file search tool args: %w", err)
	}

	return responses.ToolParamOfFileSearch(fileSearchArgs.VectorStoreIDs), nil
}

// encodeWebSearchTool creates a web search tool parameter
func encodeWebSearchTool(tool *api.ProviderDefinedTool) (responses.ToolUnionParam, error) {
	webSearchArgs, err := convertArgs[WebSearchToolArgs](tool.Args)
	if err != nil {
		return responses.ToolUnionParam{}, fmt.Errorf("failed to convert web search tool args: %w", err)
	}

	// Create a web search tool param directly instead of using the helper function
	var webSearchParam responses.WebSearchToolParam
	webSearchParam.Type = responses.WebSearchToolTypeWebSearchPreview

	// Set search context size if provided
	if webSearchArgs.SearchContextSize != "" {
		webSearchParam.SearchContextSize = responses.WebSearchToolSearchContextSize(webSearchArgs.SearchContextSize)
	}

	// Set user location if provided
	if webSearchArgs.UserLocation != nil {
		userLocation := responses.WebSearchToolUserLocationParam{}

		if webSearchArgs.UserLocation.City != "" {
			userLocation.City = openai.String(webSearchArgs.UserLocation.City)
		}

		if webSearchArgs.UserLocation.Country != "" {
			userLocation.Country = openai.String(webSearchArgs.UserLocation.Country)
		}

		if webSearchArgs.UserLocation.Region != "" {
			userLocation.Region = openai.String(webSearchArgs.UserLocation.Region)
		}

		if webSearchArgs.UserLocation.Timezone != "" {
			userLocation.Timezone = openai.String(webSearchArgs.UserLocation.Timezone)
		}

		// Only set the UserLocation if at least one field was set
		if !param.IsOmitted(userLocation.City) || !param.IsOmitted(userLocation.Country) ||
			!param.IsOmitted(userLocation.Region) || !param.IsOmitted(userLocation.Timezone) {
			webSearchParam.UserLocation = userLocation
		}
	}

	return responses.ToolUnionParam{OfWebSearchPreview: &webSearchParam}, nil
}

// encodeComputerUseTool creates a computer use tool parameter
func encodeComputerUseTool(tool *api.ProviderDefinedTool) (responses.ToolUnionParam, error) {
	computerArgs, err := convertArgs[ComputerUseToolArgs](tool.Args)
	if err != nil {
		return responses.ToolUnionParam{}, fmt.Errorf("failed to convert computer use tool args: %w", err)
	}

	// Validate required parameters
	if computerArgs.DisplayHeight <= 0 {
		return responses.ToolUnionParam{}, fmt.Errorf("displayHeight is required and must be positive")
	}

	if computerArgs.DisplayWidth <= 0 {
		return responses.ToolUnionParam{}, fmt.Errorf("displayWidth is required and must be positive")
	}

	if computerArgs.Environment == "" {
		return responses.ToolUnionParam{}, fmt.Errorf("environment is required")
	}

	// Validate that environment is one of the allowed values
	env := computerArgs.Environment
	if env != "mac" && env != "windows" && env != "ubuntu" && env != "browser" {
		return responses.ToolUnionParam{}, fmt.Errorf("environment must be one of: mac, windows, ubuntu, browser; got %q", env)
	}

	return responses.ToolParamOfComputerUsePreview(
		int64(computerArgs.DisplayHeight),
		int64(computerArgs.DisplayWidth),
		responses.ComputerToolEnvironment(computerArgs.Environment),
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
