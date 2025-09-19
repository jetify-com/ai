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

	// Don't merge system messages
	if len(group) > 0 {
		if _, isSystem := group[0].(*api.SystemMessage); isSystem {
			*result = append(*result, group...)
			return
		}
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
//   - For the last block of EACH message:
//     Message metadata is merged with block metadata (block takes precedence)
func mergeUserMessages(messages []api.Message) (api.Message, error) {
	var combinedContent []api.ContentBlock
	for _, msg := range messages {
		user, ok := msg.(*api.UserMessage)
		if !ok {
			return nil, fmt.Errorf("expected UserMessage, got %T", msg)
		}

		// Process every message for metadata precedence
		for j, block := range user.Content {
			isLastBlock := j == len(user.Content)-1

			if isLastBlock {
				// For the last block of each message, apply metadata precedence
				applyMetadataPrecedenceToBlock(block, user.ProviderMetadata)
			}
			combinedContent = append(combinedContent, block)
		}
	}
	return &api.UserMessage{Content: combinedContent}, nil
}

// mergeAssistantMessages combines multiple assistant messages into a single message.
// Metadata handling:
//   - Each block keeps its own metadata if it has any
//   - For the last block of EACH message:
//     Message metadata is merged with block metadata (block takes precedence)
func mergeAssistantMessages(messages []api.Message) (api.Message, error) {
	var combinedContent []api.ContentBlock
	for _, msg := range messages {
		assistant, ok := msg.(*api.AssistantMessage)
		if !ok {
			return nil, fmt.Errorf("expected AssistantMessage, got %T", msg)
		}

		// Process every message for metadata precedence
		for j, block := range assistant.Content {
			isLastBlock := j == len(assistant.Content)-1

			if isLastBlock {
				// For the last block of each message, apply metadata precedence
				applyMetadataPrecedenceToBlock(block, assistant.ProviderMetadata)
			}
			combinedContent = append(combinedContent, block)
		}
	}
	return &api.AssistantMessage{Content: combinedContent}, nil
}

// applyMetadataPrecedenceToBlock applies metadata precedence:
// - If block has no metadata, use message metadata
// - If block has metadata, block metadata takes complete precedence (no merging)
func applyMetadataPrecedenceToBlock(block api.ContentBlock, messageMetadata *api.ProviderMetadata) {
	if messageMetadata == nil || messageMetadata.IsZero() {
		return
	}

	switch b := block.(type) {
	case *api.TextBlock:
		if b.ProviderMetadata.IsZero() {
			b.ProviderMetadata = messageMetadata
		}
		// If block has metadata, don't apply message metadata (precedence)
	case *api.ToolCallBlock:
		if b.ProviderMetadata.IsZero() {
			b.ProviderMetadata = messageMetadata
		}
		// If block has metadata, don't apply message metadata (precedence)
	case *api.ReasoningBlock:
		if b.ProviderMetadata.IsZero() {
			b.ProviderMetadata = messageMetadata
		}
		// If block has metadata, don't apply message metadata (precedence)
	case *api.ImageBlock:
		if b.ProviderMetadata.IsZero() {
			b.ProviderMetadata = messageMetadata
		}
		// If block has metadata, don't apply message metadata (precedence)
	case *api.FileBlock:
		if b.ProviderMetadata.IsZero() {
			b.ProviderMetadata = messageMetadata
		}
		// If block has metadata, don't apply message metadata (precedence)
	}
}

// mergeToolMessages combines multiple tool messages into a single message.
// Metadata handling:
//   - Each block keeps its own metadata if it has any
//   - For the last block of EACH message:
//     Message metadata is merged with block metadata (block takes precedence)
func mergeToolMessages(messages []api.Message) (api.Message, error) {
	var combinedContent []api.ToolResultBlock
	for _, msg := range messages {
		tool, ok := msg.(*api.ToolMessage)
		if !ok {
			return nil, fmt.Errorf("expected ToolMessage, got %T", msg)
		}

		// Process every message for metadata precedence
		for j, block := range tool.Content {
			isLastBlock := j == len(tool.Content)-1

			if isLastBlock {
				// For the last block of each message, apply metadata precedence
				if messageMetadata := tool.ProviderMetadata; messageMetadata != nil && !messageMetadata.IsZero() {
					if block.ProviderMetadata.IsZero() {
						block.ProviderMetadata = messageMetadata
					}
					// If block has metadata, don't apply message metadata (precedence)
				}
			}
			combinedContent = append(combinedContent, block)
		}
	}
	return &api.ToolMessage{Content: combinedContent}, nil
}
