package codec

import (
	"fmt"

	"go.jetify.com/ai/api"
)

// TODO: We might want to combine consecutive messages of the same type for all
// providers, and not just Anthropic.
// In that case, we could move this to the SDK core before a provider is called.

// mergeMessages combines consecutive messages of the same type into single messages.
// This reduces the total number of messages while preserving the semantic meaning.
func mergeMessages(prompt []api.Message) []api.Message {
	if len(prompt) == 0 {
		return prompt
	}

	result := make([]api.Message, 0, len(prompt))
	var currentGroup []api.Message

	for _, msg := range prompt {
		if len(currentGroup) > 0 && currentGroup[0].Role() != msg.Role() {
			appendMessageGroup(&result, currentGroup)
			currentGroup = currentGroup[:0] // Clear slice but keep capacity
		}
		currentGroup = append(currentGroup, msg)
	}
	appendMessageGroup(&result, currentGroup)
	return result
}

// appendMessageGroup processes a group of messages and appends them to the result.
// If the messages can be combined, they are combined into a single message.
// If combination fails, the messages are appended individually.
func appendMessageGroup(result *[]api.Message, group []api.Message) {
	if len(group) == 0 {
		return
	}
	if len(group) == 1 {
		*result = append(*result, group[0])
		return
	}

	combined, err := mergeMessageGroup(group)
	if err != nil {
		*result = append(*result, group...)
	} else {
		*result = append(*result, combined)
	}
}

// mergeMessageGroup combines multiple messages of the same type into a single message.
// Metadata handling follows these rules:
//   - For system messages: preserves the last message's metadata
//   - For other messages:
//   - Each block keeps its own metadata if it has any
//   - For the last block of the last message only:
//     If the block has no metadata, it gets the message's metadata
func mergeMessageGroup(messages []api.Message) (api.Message, error) {
	if len(messages) == 0 {
		return nil, nil
	}

	switch first := messages[0].(type) {
	case *api.SystemMessage:
		return mergeSystemMessages(messages)
	case *api.UserMessage:
		return mergeUserMessages(messages)
	case *api.AssistantMessage:
		return mergeAssistantMessages(messages)
	case *api.ToolMessage:
		return mergeToolMessages(messages)
	default:
		return nil, fmt.Errorf("unsupported message type: %T", first)
	}
}

// mergeSystemMessages combines multiple system messages into a single message,
// joining their content with newlines and preserving the last message's metadata.
func mergeSystemMessages(messages []api.Message) (api.Message, error) {
	var combinedContent string
	for i, msg := range messages {
		sys, ok := msg.(*api.SystemMessage)
		if !ok {
			return nil, fmt.Errorf("expected SystemMessage, got %T", msg)
		}
		if i > 0 {
			combinedContent += "\n"
		}
		combinedContent += sys.Content
	}
	// Use the last message's metadata
	lastMsg, ok := messages[len(messages)-1].(*api.SystemMessage)
	if !ok {
		return nil, fmt.Errorf("expected SystemMessage, got %T", messages[len(messages)-1])
	}
	return &api.SystemMessage{
		Content:          combinedContent,
		ProviderMetadata: lastMsg.ProviderMetadata,
	}, nil
}

// mergeUserMessages combines multiple user messages into a single message.
// Metadata handling:
//   - Each block keeps its own metadata if it has any
//   - For the last block of the last message only:
//     If the block has no metadata, it gets the message's metadata
func mergeUserMessages(messages []api.Message) (api.Message, error) {
	var combinedContent []api.ContentBlock
	for i, msg := range messages {
		user, ok := msg.(*api.UserMessage)
		if !ok {
			return nil, fmt.Errorf("expected UserMessage, got %T", msg)
		}
		isLastMessage := i == len(messages)-1

		// For all blocks except the last one in the last message,
		// just append with their own metadata
		if !isLastMessage {
			combinedContent = append(combinedContent, user.Content...)
			continue
		}

		// For the last message, handle the last block specially
		for j, block := range user.Content {
			if j == len(user.Content)-1 {
				// For the last block of the last message, preserve message metadata if block has none
				switch b := block.(type) {
				case *api.TextBlock:
					if b.ProviderMetadata.IsZero() {
						b.ProviderMetadata = user.ProviderMetadata
					}
				case *api.ImageBlock:
					if b.ProviderMetadata.IsZero() {
						b.ProviderMetadata = user.ProviderMetadata
					}
				case *api.FileBlock:
					if b.ProviderMetadata.IsZero() {
						b.ProviderMetadata = user.ProviderMetadata
					}
				}
			}
			combinedContent = append(combinedContent, block)
		}
	}
	return &api.UserMessage{Content: combinedContent}, nil
}

// mergeAssistantMessages combines multiple assistant messages into a single message.
// Metadata handling:
//   - Each block keeps its own metadata if it has any
//   - For the last block of the last message only:
//     If the block has no metadata, it gets the message's metadata
func mergeAssistantMessages(messages []api.Message) (api.Message, error) {
	var combinedContent []api.ContentBlock
	for i, msg := range messages {
		assistant, ok := msg.(*api.AssistantMessage)
		if !ok {
			return nil, fmt.Errorf("expected AssistantMessage, got %T", msg)
		}
		isLastMessage := i == len(messages)-1

		if !isLastMessage {
			// For all blocks except the last message, just append with their own metadata
			combinedContent = append(combinedContent, assistant.Content...)
			continue
		}

		// Process the last message specially
		processLastAssistantMessage(&combinedContent, assistant)
	}
	return &api.AssistantMessage{Content: combinedContent}, nil
}

// processLastAssistantMessage handles the special case of the last message in a sequence
// where metadata needs to be preserved for the last block if it has none
func processLastAssistantMessage(combinedContent *[]api.ContentBlock, assistant *api.AssistantMessage) {
	for j, block := range assistant.Content {
		isLastBlock := j == len(assistant.Content)-1

		if isLastBlock {
			// For the last block of the last message, preserve message metadata if block has none
			preserveMetadataForLastBlock(block, assistant.ProviderMetadata)
		}

		*combinedContent = append(*combinedContent, block)
	}
}

// preserveMetadataForLastBlock applies message-level metadata to the last content block
// if the block doesn't already have metadata
func preserveMetadataForLastBlock(block api.ContentBlock, metadata *api.ProviderMetadata) {
	switch b := block.(type) {
	case *api.TextBlock:
		if b.ProviderMetadata.IsZero() {
			b.ProviderMetadata = metadata
		}
	case *api.ToolCallBlock:
		if b.ProviderMetadata.IsZero() {
			b.ProviderMetadata = metadata
		}
	case *api.ReasoningBlock:
		if b.ProviderMetadata.IsZero() {
			b.ProviderMetadata = metadata
		}
	case *api.RedactedReasoningBlock:
		if b.ProviderMetadata.IsZero() {
			b.ProviderMetadata = metadata
		}
	}
}

// mergeToolMessages combines multiple tool messages into a single message.
// Metadata handling:
//   - Each block keeps its own metadata if it has any
//   - For the last block of the last message only:
//     If the block has no metadata, it gets the message's metadata
func mergeToolMessages(messages []api.Message) (api.Message, error) {
	var combinedContent []api.ToolResultBlock
	for i, msg := range messages {
		tool, ok := msg.(*api.ToolMessage)
		if !ok {
			return nil, fmt.Errorf("expected ToolMessage, got %T", msg)
		}
		isLastMessage := i == len(messages)-1

		// For all blocks except the last one in the last message,
		// just append with their own metadata
		if !isLastMessage {
			combinedContent = append(combinedContent, tool.Content...)
			continue
		}

		// For the last message, handle the last block specially
		for j, block := range tool.Content {
			if j == len(tool.Content)-1 && block.ProviderMetadata.IsZero() {
				// For the last block of the last message, preserve message metadata if block has none
				block.ProviderMetadata = tool.ProviderMetadata
			}
			combinedContent = append(combinedContent, block)
		}
	}
	return &api.ToolMessage{Content: combinedContent}, nil
}
