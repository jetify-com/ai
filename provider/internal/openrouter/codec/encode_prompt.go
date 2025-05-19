package codec

// Functions that convert AI SDK types to OpenRouter types.

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/internal/openrouter/client"
)

// EncodePrompt converts an AI SDK prompt into OpenRouter's chat message format
func EncodePrompt(prompt []api.Message) (client.Prompt, error) {
	// Pre-allocate with extra space for potential tool message expansion
	messages := make(client.Prompt, 0, len(prompt)*2)

	for _, msg := range prompt {
		encodedMsgs, err := encodeMessage(msg)
		if err != nil {
			return nil, err
		}
		messages = append(messages, encodedMsgs...)
	}

	return messages, nil
}

func encodeMessage(msg api.Message) ([]client.Message, error) {
	switch msg := msg.(type) {
	case *api.SystemMessage:
		encoded := encodeSystemMessage(msg)
		return []client.Message{encoded}, nil
	case *api.UserMessage:
		encoded, err := encodeUserMessage(msg)
		if err != nil {
			return nil, err
		}
		return []client.Message{encoded}, nil
	case *api.AssistantMessage:
		encoded, err := encodeAssistantMessage(msg)
		if err != nil {
			return nil, err
		}
		return []client.Message{encoded}, nil
	case *api.ToolMessage:
		return encodeToolMessage(msg)
	default:
		// TODO: use a more specific error type from api.
		return nil, fmt.Errorf("unsupported message type: %T", msg)
	}
}

func encodeSystemMessage(msg *api.SystemMessage) *client.SystemMessage {
	return &client.SystemMessage{
		Content: msg.Content,
	}
}

func encodeUserMessage(msg *api.UserMessage) (*client.UserMessage, error) {
	// Special case: If there's exactly one text block, use a simpler format
	// This optimization avoids creating an unnecessary array of blocks
	if len(msg.Content) == 1 {
		if textBlock, ok := msg.Content[0].(*api.TextBlock); ok {
			return &client.UserMessage{
				Content: client.UserMessageContent{
					Text: textBlock.Text,
				},
			}, nil
		}
	}

	// Otherwise encode all blocks
	parts, err := encodeUserContent(msg.Content)
	if err != nil {
		return nil, err
	}

	return &client.UserMessage{
		Content: client.UserMessageContent{
			Parts: parts,
		},
	}, nil
}

func encodeUserContent(content []api.ContentBlock) ([]client.ContentPart, error) {
	parts := make([]client.ContentPart, 0, len(content))
	for _, block := range content {
		encodedPart, err := encodeUserContentBlock(block)
		if err != nil {
			return nil, err
		}
		parts = append(parts, encodedPart)
	}
	return parts, nil
}

func encodeUserContentBlock(block api.ContentBlock) (client.ContentPart, error) {
	switch block := block.(type) {
	case *api.TextBlock:
		return encodeTextBlock(block), nil
	case *api.ImageBlock:
		return encodeImageBlock(block), nil
	case *api.FileBlock:
		return encodeFileBlock(block), nil
	default:
		return nil, fmt.Errorf("unsupported content block type: %T", block)
	}
}

func encodeTextBlock(block *api.TextBlock) *client.TextPart {
	return &client.TextPart{
		Text: block.Text,
	}
}

func encodeImageBlock(block *api.ImageBlock) *client.ImagePart {
	url := block.URL
	// If no URL is provided but we have raw image data,
	// convert it to a data URL (e.g., "data:image/jpeg;base64,/9j/4AAQ...")
	if url == "" && block.Data != nil {
		mimeType := block.MimeType
		if mimeType == "" {
			mimeType = "image/jpeg" // Default to JPEG if no mime type specified
		}
		url = fmt.Sprintf("data:%s;base64,%s",
			mimeType,
			base64.StdEncoding.EncodeToString(block.Data),
		)
	}

	imagePart := &client.ImagePart{}
	imagePart.ImageURL.URL = url
	return imagePart
}

func encodeFileBlock(block *api.FileBlock) *client.TextPart {
	text := block.URL
	if text == "" && block.Data != nil {
		if block.MimeType == "" {
			// If no mime type, treat the data as plain text
			text = string(block.Data)
		} else {
			// This extra functionality of encoding when a mime type is set is
			// not part of the TypeScript OpenRouter implementation. We've added it.
			// Double check that this is beneficial.
			text = fmt.Sprintf("data:%s;base64,%s",
				block.MimeType,
				base64.StdEncoding.EncodeToString(block.Data),
			)
		}
	}
	return &client.TextPart{
		Text: text,
	}
}

func encodeAssistantMessage(msg *api.AssistantMessage) (*client.AssistantMessage, error) {
	text := ""
	toolCalls := []client.ToolCall{}

	// Combine all text parts into a single string and collect tool calls
	for _, part := range msg.Content {
		switch block := part.(type) {
		case *api.TextBlock:
			encoded := encodeTextBlock(block)
			text += encoded.Text // Concatenate all text parts
		case *api.ToolCallBlock:
			toolCall, err := encodeToolCallBlock(block)
			if err != nil {
				return nil, err
			}
			toolCalls = append(toolCalls, toolCall)
		default:
			return nil, fmt.Errorf("unsupported assistant content block type: %T", block)
		}
	}

	return &client.AssistantMessage{
		Content:   text,
		ToolCalls: toolCalls,
	}, nil
}

func encodeToolCallBlock(block *api.ToolCallBlock) (client.ToolCall, error) {
	args, err := json.Marshal(block.Args)
	if err != nil {
		return client.ToolCall{}, fmt.Errorf("failed to marshal tool call args: %w", err)
	}

	toolCall := client.ToolCall{
		Type: "function",
		ID:   block.ToolCallID,
	}
	toolCall.Function.Name = block.ToolName
	toolCall.Function.Arguments = string(args)
	return toolCall, nil
}

func encodeToolMessage(msg *api.ToolMessage) ([]client.Message, error) {
	// Convert each tool result into a separate ToolMessage
	// Note: OpenRouter expects tool messages to be flattened into separate messages
	messages := make([]client.Message, 0, len(msg.Content))
	for _, result := range msg.Content {
		resultJSON, err := json.Marshal(result.Result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tool result: %w", err)
		}
		messages = append(messages, &client.ToolMessage{
			Content:    string(resultJSON),
			ToolCallID: result.ToolCallID,
		})
	}

	// If no results, return empty slice
	if len(messages) == 0 {
		return []client.Message{}, nil
	}

	return messages, nil
}
