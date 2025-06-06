package api

import "encoding/json"

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
	Content string

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata
}

var _ Message = &SystemMessage{}

func (m SystemMessage) Role() MessageRole { return MessageRoleSystem }

func (m SystemMessage) GetProviderMetadata() *ProviderMetadata { return m.ProviderMetadata }

// UserMessage represents a user message that can contain text, images, and files
type UserMessage struct {
	// Content contains an array of content blocks (text, image, or file)
	Content []ContentBlock

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata
}

var _ Message = &UserMessage{}

func (m UserMessage) Role() MessageRole { return MessageRoleUser }

func (m UserMessage) GetProviderMetadata() *ProviderMetadata { return m.ProviderMetadata }

// AssistantMessage represents an assistant message that can contain text and tool calls
type AssistantMessage struct {
	// Content contains an array of content blocks (text or tool calls)
	Content []ContentBlock

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata
}

var _ Message = &AssistantMessage{}

func (m AssistantMessage) Role() MessageRole { return MessageRoleAssistant }

func (m AssistantMessage) GetProviderMetadata() *ProviderMetadata { return m.ProviderMetadata }

// ToolMessage represents a tool message containing tool results
type ToolMessage struct {
	// Content contains an array of tool results
	Content []ToolResultBlock

	// ProviderMetadata contains additional provider-specific metadata.
	// They are passed through to the provider from the AI SDK and enable
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata
}

var _ Message = &ToolMessage{}

func (m ToolMessage) Role() MessageRole { return MessageRoleTool }

func (m ToolMessage) GetProviderMetadata() *ProviderMetadata { return m.ProviderMetadata }

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
