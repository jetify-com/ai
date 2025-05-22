package codec

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
	"go.jetify.com/ai/api"
)

type modelConfig struct {
	IsReasoningModel       bool   `json:"isReasoningModel"`
	SystemMessageMode      string `json:"systemMessageMode"`
	RequiredAutoTruncation bool   `json:"requiredAutoTruncation"`
}

type OpenAIPrompt struct {
	Messages []responses.ResponseInputItemUnionParam
	Warnings []api.CallWarning
}

func EncodePrompt(prompt []api.Message, modelConfig modelConfig) (*OpenAIPrompt, error) {
	messages := make([]responses.ResponseInputItemUnionParam, 0, len(prompt))
	warnings := []api.CallWarning{}

	for _, message := range prompt {
		switch m := message.(type) {
		case *api.SystemMessage:
			item, itemWarnings, err := EncodeSystemMessage(m, modelConfig)
			if err != nil {
				return nil, fmt.Errorf("encoding system message: %w", err)
			}
			if item != nil {
				messages = append(messages, *item)
			}
			warnings = append(warnings, itemWarnings...)

		case *api.UserMessage:
			item, err := EncodeUserMessage(m)
			if err != nil {
				return nil, fmt.Errorf("encoding user message: %w", err)
			}
			messages = append(messages, item)

		case *api.AssistantMessage:
			items, err := EncodeAssistantMessage(m)
			if err != nil {
				return nil, fmt.Errorf("encoding assistant message: %w", err)
			}
			messages = append(messages, items...)

		case *api.ToolMessage:
			items, err := EncodeToolMessage(m)
			if err != nil {
				return nil, fmt.Errorf("encoding tool message: %w", err)
			}
			messages = append(messages, items...)

		default:
			return nil, fmt.Errorf("unsupported message type: %T", message)
		}
	}

	return &OpenAIPrompt{
		Messages: messages,
		Warnings: warnings,
	}, nil
}

// EncodeSystemMessage converts a system message to OpenAI format
func EncodeSystemMessage(message *api.SystemMessage, config modelConfig) (*responses.ResponseInputItemUnionParam, []api.CallWarning, error) {
	warnings := []api.CallWarning{}

	switch config.SystemMessageMode {
	case "", "system":
		// Create the message with role "system" using direct string content
		item := responses.ResponseInputItemParamOfMessage(message.Content, responses.EasyInputMessageRoleSystem)
		return &item, warnings, nil

	case "developer":
		// Create the message with role "developer" using direct string content
		item := responses.ResponseInputItemParamOfMessage(message.Content, responses.EasyInputMessageRoleDeveloper)
		return &item, warnings, nil

	case "remove":
		// Don't add the message, just add a warning
		warnings = append(warnings, api.CallWarning{
			Type:    "other",
			Message: "system messages are removed for this model",
		})
		return nil, warnings, nil

	default:
		return nil, warnings, fmt.Errorf("unsupported system message mode: %s", config.SystemMessageMode)
	}
}

// EncodeUserMessage converts a user message to OpenAI format
func EncodeUserMessage(message *api.UserMessage) (responses.ResponseInputItemUnionParam, error) {
	content := make([]responses.ResponseInputContentUnionParam, 0, len(message.Content))

	for _, block := range message.Content {
		item, err := EncodeUserContentBlock(block)
		if err != nil {
			return responses.ResponseInputItemUnionParam{}, fmt.Errorf("failed to encode content block: %w", err)
		}
		if item != nil {
			content = append(content, *item)
		}
	}

	return responses.ResponseInputItemParamOfInputMessage(content, "user"), nil
}

// EncodeUserContentBlock converts a content block to OpenAI format
func EncodeUserContentBlock(block api.ContentBlock) (*responses.ResponseInputContentUnionParam, error) {
	switch b := block.(type) {
	case *api.TextBlock:
		return EncodeInputTextBlock(b)

	case *api.ImageBlock:
		return EncodeImageBlock(b)

	case *api.FileBlock:
		return EncodeFileBlock(b)

	default:
		return nil, fmt.Errorf("unsupported content block type: %T", block)
	}
}

// EncodeInputTextBlock converts a text block to OpenAI format
func EncodeInputTextBlock(block *api.TextBlock) (*responses.ResponseInputContentUnionParam, error) {
	if block == nil {
		return nil, fmt.Errorf("text block cannot be nil")
	}
	param := responses.ResponseInputContentParamOfInputText(block.Text)
	return &param, nil
}

// EncodeOutputTextBlock converts a text block to OpenAI format for assistant messages
func EncodeOutputTextBlock(block *api.TextBlock) (*responses.ResponseInputItemUnionParam, error) {
	if block == nil {
		return nil, fmt.Errorf("text block cannot be nil")
	}
	textParam := &responses.ResponseOutputTextParam{
		Text: block.Text,
		Type: "output_text",
	}

	// Create content array with the output text
	content := []responses.ResponseOutputMessageContentUnionParam{
		{
			OfOutputText: textParam,
		},
	}

	// Create a message with this content and role "assistant"
	// Passing empty string for ID and no status
	item := responses.ResponseInputItemParamOfOutputMessage(content, "", "")
	return &item, nil
}

// EncodeFileBlock converts a file block to OpenAI format
func EncodeFileBlock(block *api.FileBlock) (*responses.ResponseInputContentUnionParam, error) {
	if block == nil {
		return nil, fmt.Errorf("file block cannot be nil")
	}

	// Create the file param
	fileParam := responses.ResponseInputFileParam{}

	// Check if we have a URL or Data
	if block.URL != "" {
		// For OpenAI we can't include a URL directly. We would need to upload a
		// file using the file API and then use a file ID that needs to happen as
		// pre-processing before we call the provider if we've made it here with a
		// URL it's an error
		return nil, fmt.Errorf("file URLs in user messages are not supported")
	}

	if block.Data == nil {
		return nil, fmt.Errorf("file block must have either URL or Data")
	}

	// For file data, check the mime type
	if block.MediaType != "application/pdf" {
		return nil, fmt.Errorf("only PDF files are supported in user messages")
	}

	// Encode the PDF data as base64
	base64Data := base64.StdEncoding.EncodeToString(block.Data)

	// Generate filename - use metadata if available, otherwise default
	// TODO: If we could set the file name based on the index of the block that would help debug
	// Or maybe we should make sure every block always has an ID.
	filename := "file.pdf"
	metadata := GetMetadata(block)
	if metadata != nil && metadata.Filename != "" {
		filename = metadata.Filename
	}

	// Set the file data with data URL format
	fileParam.Filename = openai.String(filename)
	dataURL := fmt.Sprintf("data:%s;base64,%s", block.MediaType, base64Data)
	fileParam.FileData = openai.String(dataURL)

	// Create content union param with the file
	contentParam := responses.ResponseInputContentUnionParam{
		OfInputFile: &fileParam,
	}

	return &contentParam, nil
}

// EncodeImageBlock converts an image block to OpenAI format
func EncodeImageBlock(block *api.ImageBlock) (*responses.ResponseInputContentUnionParam, error) {
	if block == nil {
		return nil, fmt.Errorf("image block cannot be nil")
	}

	// Get OpenAI-specific metadata for image detail level
	metadata := GetMetadata(block)

	// Create the image param
	imageParam := responses.ResponseInputImageParam{}

	// Only set the detail if explicitly provided in metadata
	if metadata != nil && metadata.ImageDetail != "" {
		detailLevel := metadata.ImageDetail
		switch detailLevel {
		case "high":
			imageParam.Detail = responses.ResponseInputImageDetailHigh
		case "low":
			imageParam.Detail = responses.ResponseInputImageDetailLow
		case "auto":
			imageParam.Detail = responses.ResponseInputImageDetailAuto
		default:
			return nil, fmt.Errorf("invalid image detail level: %s (must be one of 'high', 'low', or 'auto')", detailLevel)
		}
	}

	// Check if we have a URL or Data
	if block.URL != "" {
		// For URL-based images
		imageParam.ImageURL = openai.String(block.URL)
	} else if block.Data != nil {
		// For base64 data, create a data URL
		mimeType := "image/jpeg" // Default mime type
		if block.MediaType != "" {
			mimeType = block.MediaType
		}

		// Properly encode the binary data as base64
		base64Data := base64.StdEncoding.EncodeToString(block.Data)

		// Create the data URL with the encoded data
		dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)
		imageParam.ImageURL = openai.String(dataURL)
	} else {
		return nil, fmt.Errorf("image block must have either URL or Data")
	}

	// Create content union param with the image
	contentParam := responses.ResponseInputContentUnionParam{
		OfInputImage: &imageParam,
	}

	return &contentParam, nil
}

// EncodeToolCallBlock converts a tool call block to OpenAI format
func EncodeToolCallBlock(block *api.ToolCallBlock) (*responses.ResponseInputItemUnionParam, error) {
	if block == nil {
		return nil, fmt.Errorf("tool call block cannot be nil")
	}

	if block.ToolName == "" {
		return nil, fmt.Errorf("tool call is missing tool name")
	}

	// Marshal the arguments to JSON
	argsJSON, err := json.Marshal(block.Args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tool call arguments: %v", err)
	}
	arguments := string(argsJSON)

	// Create a function call item
	item := responses.ResponseInputItemParamOfFunctionCall(
		arguments,
		block.ToolCallID,
		block.ToolName,
	)

	return &item, nil
}

// EncodeAssistantMessage converts an assistant message to OpenAI format
func EncodeAssistantMessage(message *api.AssistantMessage) ([]responses.ResponseInputItemUnionParam, error) {
	items := make([]responses.ResponseInputItemUnionParam, 0, len(message.Content))

	for _, block := range message.Content {
		switch b := block.(type) {
		case *api.TextBlock:
			item, err := EncodeOutputTextBlock(b)
			if err != nil {
				return nil, fmt.Errorf("encoding text block: %w", err)
			}
			items = append(items, *item)
		case *api.ToolCallBlock:
			// Handle the tool call as a separate item
			item, err := EncodeToolCallBlock(b)
			if err != nil {
				return nil, fmt.Errorf("encoding tool call block: %w", err)
			}
			if item != nil {
				items = append(items, *item)
			}
		// TODO: Handle reasoning blocks in assistant messages
		default:
			return nil, fmt.Errorf("unsupported content block type in assistant message: %T", block)
		}
	}

	return items, nil
}

// EncodeToolMessage converts a tool message to OpenAI format
func EncodeToolMessage(message *api.ToolMessage) ([]responses.ResponseInputItemUnionParam, error) {
	items := make([]responses.ResponseInputItemUnionParam, 0, len(message.Content))

	for i := range message.Content {
		item, err := EncodeToolResultBlock(&message.Content[i])
		if err != nil {
			return nil, fmt.Errorf("encoding tool result block: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// encodeComputerToolResult handles encoding the result of a computer use tool call
func encodeComputerToolResult(result *api.ToolResultBlock) (responses.ResponseInputItemUnionParam, error) {
	if len(result.Content) != 1 {
		return responses.ResponseInputItemUnionParam{}, fmt.Errorf("expected 1 content block for computer use tool result, got %d", len(result.Content))
	}

	content := result.Content[0]
	imageBlock, ok := content.(api.ImageBlock)
	if !ok {
		return responses.ResponseInputItemUnionParam{}, fmt.Errorf("expected image block for computer use tool result, got %T", content)
	}

	// Create data URL from image data
	// TODO: Add helper methods to the image and file blocks to make this easier
	dataURL := "data:" + imageBlock.MediaType + ";base64," + base64.StdEncoding.EncodeToString(imageBlock.Data)

	screenshot := responses.ResponseComputerToolCallOutputScreenshotParam{
		Type:     "computer_screenshot",
		ImageURL: openai.String(dataURL),
	}

	// Extract safety checks from provider metadata if available
	var acknowledgedSafetyChecks []responses.ResponseInputItemComputerCallOutputAcknowledgedSafetyCheckParam
	if metadata := GetMetadata(result); metadata != nil {
		for _, check := range metadata.ComputerSafetyChecks {
			acknowledgedSafetyChecks = append(acknowledgedSafetyChecks, responses.ResponseInputItemComputerCallOutputAcknowledgedSafetyCheckParam{
				ID:      check.ID,
				Code:    openai.String(check.Code),
				Message: openai.String(check.Message),
			})
		}
	}

	// Create the computer call output parameter
	output := responses.ResponseInputItemComputerCallOutputParam{
		CallID:                   result.ToolCallID,
		Output:                   screenshot,
		AcknowledgedSafetyChecks: acknowledgedSafetyChecks,
	}

	return responses.ResponseInputItemUnionParam{
		OfComputerCallOutput: &output,
	}, nil
}

func EncodeToolResultBlock(result *api.ToolResultBlock) (responses.ResponseInputItemUnionParam, error) {
	if result.ToolCallID == "" {
		return responses.ResponseInputItemUnionParam{}, fmt.Errorf("tool result is missing tool call ID")
	}

	// Handle computer use tool output
	if result.ToolCallID == "openai.computer_use_preview" {
		return encodeComputerToolResult(result)
	}

	// Handle regular function tool output
	resultJSON, err := json.Marshal(result.Result)
	if err != nil {
		return responses.ResponseInputItemUnionParam{}, fmt.Errorf("failed to marshal tool result: %v", err)
	}
	output := string(resultJSON)

	// Create a function call output item
	return responses.ResponseInputItemParamOfFunctionCallOutput(
		result.ToolCallID,
		output,
	), nil
}
