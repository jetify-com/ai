package api

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

// MessageRole represents the role of a message sender
type MessageRole string

const (
	// MessageRoleSystem represents a system message
	MessageRoleSystem MessageRole = "system"
	// MessageRoleUser represents a user message
	MessageRoleUser MessageRole = "user"
	// MessageRoleAssistant represents an assistant message
	MessageRoleAssistant MessageRole = "assistant"
	// MessageRoleTool represents a tool message
	MessageRoleTool MessageRole = "tool"
)

// Message represents a message in a prompt sequence.
// Note: Not all models and prompt formats support multi-modal inputs and
// tool calls. The validation happens at runtime.
//
// Note: This is not a user-facing prompt. The AI SDK methods will map the
// user-facing prompt types such as chat or instruction prompts to this format.
//
// Note: there could be additional blocks for each role in the future,
// e.g. when the assistant can return images or the user can share files
// such as PDFs.
type Message interface {
	// Role returns the role of the message sender.
	// Valid roles are: "system", "user", "assistant", or "tool".
	Role() MessageRole
	GetProviderMetadata() *ProviderMetadata
	// TODO: should we add a Content() method?

	// TODO: Decide if we should "flatten" the different message types into a single
	// Message struct (a concrete type instead of an interface).
}

// SystemMessage represents a system message with plain text content
type SystemMessage struct {
	// Content contains the message text
	Content string `json:"content"`

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

var _ Message = &SystemMessage{}

func (m SystemMessage) Role() MessageRole { return MessageRoleSystem }

func (m SystemMessage) GetProviderMetadata() *ProviderMetadata { return m.ProviderMetadata }

// MarshalJSON includes the role field when marshaling SystemMessage
func (m SystemMessage) MarshalJSON() ([]byte, error) {
	type Alias SystemMessage
	return json.Marshal(struct {
		Role string `json:"role"`
		*Alias
	}{
		Role:  string(MessageRoleSystem),
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for SystemMessage
func (m *SystemMessage) UnmarshalJSON(data []byte) error {
	// Use gjson to check role first
	roleResult := gjson.GetBytes(data, "role")
	if !roleResult.Exists() {
		return fmt.Errorf("system message missing required 'role' field")
	}
	if roleResult.String() != string(MessageRoleSystem) {
		return fmt.Errorf("invalid role for SystemMessage: expected 'system', got '%s'", roleResult.String())
	}

	// Use alias to unmarshal the rest
	type Alias SystemMessage
	aux := (*Alias)(m)
	return json.Unmarshal(data, aux)
}

// UserMessage represents a user message that can contain text, images, and files
type UserMessage struct {
	// Content contains an array of content blocks (text, image, or file)
	Content []ContentBlock `json:"content"`

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

var _ Message = &UserMessage{}

func (m UserMessage) Role() MessageRole { return MessageRoleUser }

func (m UserMessage) GetProviderMetadata() *ProviderMetadata { return m.ProviderMetadata }

// MarshalJSON includes the role field when marshaling UserMessage
func (m UserMessage) MarshalJSON() ([]byte, error) {
	type Alias UserMessage
	return json.Marshal(struct {
		Role string `json:"role"`
		*Alias
	}{
		Role:  string(MessageRoleUser),
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for UserMessage
func (m *UserMessage) UnmarshalJSON(data []byte) error {
	// Use gjson to check role first
	roleResult := gjson.GetBytes(data, "role")
	if !roleResult.Exists() {
		return fmt.Errorf("user message missing required 'role' field")
	}
	if roleResult.String() != string(MessageRoleUser) {
		return fmt.Errorf("invalid role for UserMessage: expected 'user', got '%s'", roleResult.String())
	}

	// Use a temporary struct to unmarshal everything except content
	type UserMessageAlias UserMessage
	temp := struct {
		*UserMessageAlias
		Content []json.RawMessage `json:"content,omitempty"`
	}{
		UserMessageAlias: (*UserMessageAlias)(m),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Handle content blocks
	if temp.Content != nil {
		content, err := unmarshalContentBlocks(temp.Content)
		if err != nil {
			return err
		}
		m.Content = content
	}

	return nil
}

// AssistantMessage represents an assistant message that can contain text and tool calls
type AssistantMessage struct {
	// Content contains an array of content blocks (text or tool calls)
	Content []ContentBlock `json:"content"`

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

var _ Message = &AssistantMessage{}

func (m AssistantMessage) Role() MessageRole { return MessageRoleAssistant }

func (m AssistantMessage) GetProviderMetadata() *ProviderMetadata { return m.ProviderMetadata }

// MarshalJSON includes the role field when marshaling AssistantMessage
func (m AssistantMessage) MarshalJSON() ([]byte, error) {
	type Alias AssistantMessage
	return json.Marshal(struct {
		Role string `json:"role"`
		*Alias
	}{
		Role:  string(MessageRoleAssistant),
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for AssistantMessage
func (m *AssistantMessage) UnmarshalJSON(data []byte) error {
	// Use gjson to check role first
	roleResult := gjson.GetBytes(data, "role")
	if !roleResult.Exists() {
		return fmt.Errorf("assistant message missing required 'role' field")
	}
	if roleResult.String() != string(MessageRoleAssistant) {
		return fmt.Errorf("invalid role for AssistantMessage: expected 'assistant', got '%s'", roleResult.String())
	}

	// Use a temporary struct to unmarshal everything except content
	type AssistantMessageAlias AssistantMessage
	temp := struct {
		*AssistantMessageAlias
		Content []json.RawMessage `json:"content,omitempty"`
	}{
		AssistantMessageAlias: (*AssistantMessageAlias)(m),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Handle content blocks
	if temp.Content != nil {
		content, err := unmarshalContentBlocks(temp.Content)
		if err != nil {
			return err
		}
		m.Content = content
	}

	return nil
}

// ToolMessage represents a tool message containing tool results
type ToolMessage struct {
	// Content contains an array of tool results
	Content []ToolResultBlock `json:"content"`

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

var _ Message = &ToolMessage{}

func (m ToolMessage) Role() MessageRole { return MessageRoleTool }

func (m ToolMessage) GetProviderMetadata() *ProviderMetadata { return m.ProviderMetadata }

// MarshalJSON includes the role field when marshaling ToolMessage
func (m ToolMessage) MarshalJSON() ([]byte, error) {
	type Alias ToolMessage
	return json.Marshal(struct {
		Role string `json:"role"`
		*Alias
	}{
		Role:  string(MessageRoleTool),
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for ToolMessage
func (m *ToolMessage) UnmarshalJSON(data []byte) error {
	// Use gjson to check role first
	roleResult := gjson.GetBytes(data, "role")
	if !roleResult.Exists() {
		return fmt.Errorf("tool message missing required 'role' field")
	}
	if roleResult.String() != string(MessageRoleTool) {
		return fmt.Errorf("invalid role for ToolMessage: expected 'tool', got '%s'", roleResult.String())
	}

	// Use a temporary struct to unmarshal everything except content
	type ToolMessageAlias ToolMessage
	temp := struct {
		*ToolMessageAlias
		Content []json.RawMessage `json:"content,omitempty"`
	}{
		ToolMessageAlias: (*ToolMessageAlias)(m),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Handle content blocks - tool messages contain ToolResultBlock
	if temp.Content != nil {
		content := make([]ToolResultBlock, len(temp.Content))
		for i, blockData := range temp.Content {
			// Use gjson to extract type
			typeResult := gjson.GetBytes(blockData, "type")
			if !typeResult.Exists() {
				return fmt.Errorf("tool result block at index %d missing required 'type' field", i)
			}

			blockType := typeResult.String()
			if blockType != string(ContentBlockTypeToolResult) {
				return fmt.Errorf("invalid type for tool result block at index %d: expected 'tool-result', got '%s'", i, blockType)
			}

			var toolResult ToolResultBlock
			if err := json.Unmarshal(blockData, &toolResult); err != nil {
				return fmt.Errorf("failed to unmarshal tool result block at index %d: %w", i, err)
			}
			content[i] = toolResult
		}
		m.Content = content
	}

	return nil
}

// ContentFromText creates a slice of content blocks with a single text block.
func ContentFromText(text string) []ContentBlock {
	return []ContentBlock{
		&TextBlock{Text: text},
	}
}

// ContentBlockType represents the type of content block in a message
type ContentBlockType string

const (
	// ContentBlockTypeText represents text content
	ContentBlockTypeText ContentBlockType = "text"
	// ContentBlockTypeImage represents image content
	ContentBlockTypeImage ContentBlockType = "image"
	// ContentBlockTypeFile represents file content
	ContentBlockTypeFile ContentBlockType = "file"
	// ContentBlockTypeToolCall represents a tool call
	ContentBlockTypeToolCall ContentBlockType = "tool-call"
	// ContentBlockTypeToolResult represents a tool result
	ContentBlockTypeToolResult ContentBlockType = "tool-result"
	// ContentBlockTypeReasoning represents reasoning text from the model
	ContentBlockTypeReasoning ContentBlockType = "reasoning"
	// ContentBlockTypeRedactedReasoning represents redacted reasoning data
	ContentBlockTypeRedactedReasoning ContentBlockType = "redacted-reasoning"
	// ContentBlockTypeSource represents a source block
	ContentBlockTypeSource ContentBlockType = "source"
)

// ContentBlock represents a block of content in a message
type ContentBlock interface {
	// Type returns the type of the content block.
	// Valid types are: "text", "image", "file", "tool-call", "tool-result",
	// "reasoning", "redacted-reasoning".
	Type() ContentBlockType
	// ProviderMetadata returns the provider-specific metadata for the content block.
	GetProviderMetadata() *ProviderMetadata
}

// Reasoning represents a reasoning content block in a message.
// It can be either a ReasoningBlock or a RedactedReasoningBlock.
type Reasoning interface {
	ContentBlock
	isReasoning()
}

// TextBlock represents text content in a message
type TextBlock struct {
	// Text contains the text content
	Text string `json:"text"`

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

var _ ContentBlock = &TextBlock{}

func (b TextBlock) Type() ContentBlockType { return ContentBlockTypeText }

func (b TextBlock) GetProviderMetadata() *ProviderMetadata { return b.ProviderMetadata }

// MarshalJSON includes the type field when marshaling TextBlock
func (b TextBlock) MarshalJSON() ([]byte, error) {
	type Alias TextBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(ContentBlockTypeText),
		Alias: (*Alias)(&b),
	})
}

// ReasoningBlock represents a reasoning content block of a prompt.
// It contains a string of reasoning text.
type ReasoningBlock struct {
	// Text contains the reasoning text.
	Text string `json:"text"`

	// Signature is an optional signature for verifying that the reasoning originated from the model.
	Signature string `json:"signature,omitzero"`

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

var (
	_ ContentBlock = &ReasoningBlock{}
	_ Reasoning    = &ReasoningBlock{}
)

func (b ReasoningBlock) Type() ContentBlockType { return ContentBlockTypeReasoning }

func (b ReasoningBlock) isReasoning() {}

func (b ReasoningBlock) GetProviderMetadata() *ProviderMetadata { return b.ProviderMetadata }

// MarshalJSON includes the type field when marshaling ReasoningBlock
func (b ReasoningBlock) MarshalJSON() ([]byte, error) {
	type Alias ReasoningBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(ContentBlockTypeReasoning),
		Alias: (*Alias)(&b),
	})
}

// RedactedReasoningBlock represents a redacted reasoning content block of a prompt.
type RedactedReasoningBlock struct {
	// Data contains redacted reasoning data.
	Data string `json:"data"`

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

var (
	_ ContentBlock = &RedactedReasoningBlock{}
	_ Reasoning    = &RedactedReasoningBlock{}
)

func (b RedactedReasoningBlock) Type() ContentBlockType { return ContentBlockTypeRedactedReasoning }

func (b RedactedReasoningBlock) isReasoning() {}

func (b RedactedReasoningBlock) GetProviderMetadata() *ProviderMetadata { return b.ProviderMetadata }

// MarshalJSON includes the type field when marshaling RedactedReasoningBlock
func (b RedactedReasoningBlock) MarshalJSON() ([]byte, error) {
	type Alias RedactedReasoningBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(ContentBlockTypeRedactedReasoning),
		Alias: (*Alias)(&b),
	})
}

// ImageBlock represents an image in a message.
// Either URL or Data should be set, but not both.
// TODO: merge into file block in language model v2
type ImageBlock struct {
	// URL is the external URL of the image.
	URL string `json:"url,omitzero"`

	// Data contains the image data as raw bytes.
	// If this is set, also set the MimeType so that the AI SDK knows
	// how to interpret the data.
	Data []byte `json:"data,omitempty"`

	// MediaType is the IANA media type (mime type) of the image
	MediaType string `json:"media_type,omitzero"`

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

var _ ContentBlock = &ImageBlock{}

func (b ImageBlock) Type() ContentBlockType { return ContentBlockTypeImage }

func (b ImageBlock) GetProviderMetadata() *ProviderMetadata { return b.ProviderMetadata }

// MarshalJSON includes the type field when marshaling ImageBlock
func (b ImageBlock) MarshalJSON() ([]byte, error) {
	type Alias ImageBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(ContentBlockTypeImage),
		Alias: (*Alias)(&b),
	})
}

// ImageBlockFromURL creates a new image block from a URL.
func ImageBlockFromURL(url string) *ImageBlock {
	return &ImageBlock{
		URL: url,
	}
}

// ImageBlockFromData creates a new image block from raw bytes.
func ImageBlockFromData(data []byte, mediaType string) *ImageBlock {
	return &ImageBlock{
		Data:      data,
		MediaType: mediaType,
	}
}

// FileBlock represents a file in a message.
// Either URL or Data should be set, but not both.
type FileBlock struct {
	// Filename is the filename of the file. Optional.
	Filename string `json:"filename,omitzero"` // TODO: input only

	// URL is the external URL of the file.
	URL string `json:"url,omitzero"`

	// Data contains the file data as raw bytes.
	// If this is set, also set the MimeType so that the AI SDK knows
	// how to interpret the data.
	Data []byte `json:"data,omitempty"`

	// MediaType is the IANA media type (mime type) of the file.
	// It can support wildcards, e.g. `image/*` (in which case the provider needs to take appropriate action).
	// See: https://www.iana.org/assignments/media-types/media-types.xhtml
	MediaType string `json:"media_type,omitzero"`

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

var _ ContentBlock = &FileBlock{}

func (b FileBlock) Type() ContentBlockType { return ContentBlockTypeFile }

func (b FileBlock) GetProviderMetadata() *ProviderMetadata { return b.ProviderMetadata }

// MarshalJSON includes the type field when marshaling FileBlock
func (b FileBlock) MarshalJSON() ([]byte, error) {
	type Alias FileBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(ContentBlockTypeFile),
		Alias: (*Alias)(&b),
	})
}

// FileBlockFromURL creates a new file block from a URL.
func FileBlockFromURL(url string) *FileBlock {
	return &FileBlock{
		URL: url,
	}
}

// FileBlockFromData creates a new file block from raw bytes.
func FileBlockFromData(data []byte, mediaType string) *FileBlock {
	return &FileBlock{
		Data:      data,
		MediaType: mediaType,
	}
}

// ToolCallBlock represents a tool call in a message (usually generated by the AI model)
type ToolCallBlock struct {
	// ToolCallID is the ID of the tool call. This ID is used to match the tool call with the tool result.
	ToolCallID string `json:"tool_call_id"`

	// ToolName is the name of the tool that is being called
	ToolName string `json:"tool_name"`

	// Args contains the arguments of the tool call as a JSON payload matching
	// the tool's input schema.
	// Note that args are often generated by the language model and may be
	// malformed.
	Args json.RawMessage `json:"args"` // TODO: decide if this is the right type for this field

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

var _ ContentBlock = &ToolCallBlock{}

func (b ToolCallBlock) Type() ContentBlockType { return ContentBlockTypeToolCall }

func (b ToolCallBlock) GetProviderMetadata() *ProviderMetadata { return b.ProviderMetadata }

// MarshalJSON includes the type field when marshaling ToolCallBlock
func (b ToolCallBlock) MarshalJSON() ([]byte, error) {
	type Alias ToolCallBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(ContentBlockTypeToolCall),
		Alias: (*Alias)(&b),
	})
}

// ToolResultBlock represents a tool result in a message. Usually sent back to the model as input,
// after it requested a tool call with a matching ToolCallID.
type ToolResultBlock struct {
	// ToolCallID is the ID of the tool call that this result is associated with
	ToolCallID string `json:"tool_call_id"`

	// ToolName is the name of the tool that generated this result
	ToolName string `json:"tool_name"`

	// Result contains the result of the tool call
	Result any `json:"result"`

	// IsError indicates if the result is an error or an error message
	IsError bool `json:"is_error,omitzero"`

	// Content contains tool results as an array of blocks.
	// This enables advanced tool results including images.
	// When this is used, the Result field should be ignored
	// (if the provider supports content).
	Content []ContentBlock `json:"content,omitempty"`

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

var _ ContentBlock = &ToolResultBlock{}

func (b ToolResultBlock) Type() ContentBlockType { return ContentBlockTypeToolResult }

func (b ToolResultBlock) GetProviderMetadata() *ProviderMetadata { return b.ProviderMetadata }

// MarshalJSON includes the type field when marshaling ToolResultBlock
func (b ToolResultBlock) MarshalJSON() ([]byte, error) {
	type Alias ToolResultBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(ContentBlockTypeToolResult),
		Alias: (*Alias)(&b),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for ToolResultBlock
func (b *ToolResultBlock) UnmarshalJSON(data []byte) error {
	// Use gjson to check type first
	typeResult := gjson.GetBytes(data, "type")
	if !typeResult.Exists() {
		return fmt.Errorf("tool result block missing required 'type' field")
	}
	if typeResult.String() != string(ContentBlockTypeToolResult) {
		return fmt.Errorf("invalid type for ToolResultBlock: expected 'tool-result', got '%s'", typeResult.String())
	}

	// Use a temporary struct to unmarshal everything except content
	type ToolResultBlockAlias ToolResultBlock
	temp := struct {
		*ToolResultBlockAlias
		Content []json.RawMessage `json:"content,omitempty"`
	}{
		ToolResultBlockAlias: (*ToolResultBlockAlias)(b),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Handle content blocks
	if temp.Content != nil {
		content, err := unmarshalContentBlocks(temp.Content)
		if err != nil {
			return err
		}
		b.Content = content
	}

	return nil
}

// SourceBlock represents a source that has been used as input to generate the response.
type SourceBlock struct {
	// ID is the ID of the source.
	ID string `json:"id"`

	// URL is the external URL of the source.
	URL string `json:"url"`

	// Title is the title of the source.
	Title string `json:"title,omitzero"`

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}

var _ ContentBlock = &SourceBlock{}

func (b SourceBlock) Type() ContentBlockType { return ContentBlockTypeSource }

func (b SourceBlock) GetProviderMetadata() *ProviderMetadata { return b.ProviderMetadata }

// MarshalJSON includes the type field when marshaling SourceBlock
func (b SourceBlock) MarshalJSON() ([]byte, error) {
	type Alias SourceBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  string(ContentBlockTypeSource),
		Alias: (*Alias)(&b),
	})
}

// unmarshalContentBlocks is a helper function to unmarshal an array of ContentBlock
func unmarshalContentBlocks(rawBlocks []json.RawMessage) ([]ContentBlock, error) {
	blocks := make([]ContentBlock, len(rawBlocks))
	for i, blockData := range rawBlocks {
		// Use gjson to extract type without full unmarshaling
		typeResult := gjson.GetBytes(blockData, "type")
		if !typeResult.Exists() {
			return nil, fmt.Errorf("content block at index %d missing required 'type' field", i)
		}

		blockType := typeResult.String()

		// Based on type, unmarshal into appropriate concrete type
		switch blockType {
		case string(ContentBlockTypeText):
			var textBlock TextBlock
			if err := json.Unmarshal(blockData, &textBlock); err != nil {
				return nil, fmt.Errorf("failed to unmarshal text block at index %d: %w", i, err)
			}
			blocks[i] = &textBlock
		case string(ContentBlockTypeImage):
			var imageBlock ImageBlock
			if err := json.Unmarshal(blockData, &imageBlock); err != nil {
				return nil, fmt.Errorf("failed to unmarshal image block at index %d: %w", i, err)
			}
			blocks[i] = &imageBlock
		case string(ContentBlockTypeFile):
			var fileBlock FileBlock
			if err := json.Unmarshal(blockData, &fileBlock); err != nil {
				return nil, fmt.Errorf("failed to unmarshal file block at index %d: %w", i, err)
			}
			blocks[i] = &fileBlock
		case string(ContentBlockTypeToolCall):
			var toolCallBlock ToolCallBlock
			if err := json.Unmarshal(blockData, &toolCallBlock); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tool call block at index %d: %w", i, err)
			}
			blocks[i] = &toolCallBlock
		case string(ContentBlockTypeToolResult):
			var toolResultBlock ToolResultBlock
			if err := json.Unmarshal(blockData, &toolResultBlock); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tool result block at index %d: %w", i, err)
			}
			blocks[i] = &toolResultBlock
		case string(ContentBlockTypeReasoning):
			var reasoningBlock ReasoningBlock
			if err := json.Unmarshal(blockData, &reasoningBlock); err != nil {
				return nil, fmt.Errorf("failed to unmarshal reasoning block at index %d: %w", i, err)
			}
			blocks[i] = &reasoningBlock
		case string(ContentBlockTypeRedactedReasoning):
			var redactedReasoningBlock RedactedReasoningBlock
			if err := json.Unmarshal(blockData, &redactedReasoningBlock); err != nil {
				return nil, fmt.Errorf("failed to unmarshal redacted reasoning block at index %d: %w", i, err)
			}
			blocks[i] = &redactedReasoningBlock
		case string(ContentBlockTypeSource):
			var sourceBlock SourceBlock
			if err := json.Unmarshal(blockData, &sourceBlock); err != nil {
				return nil, fmt.Errorf("failed to unmarshal source block at index %d: %w", i, err)
			}
			blocks[i] = &sourceBlock
		default:
			return nil, fmt.Errorf("unknown content block type '%s' at index %d", blockType, i)
		}
	}
	return blocks, nil
}
