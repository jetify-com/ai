package client

// Role constants for message types
const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleTool      = "tool"
)

// Content type constants
const (
	ContentTypeText     = "text"
	ContentTypeFunction = "function"
	ContentTypeImageURL = "image_url"
)

// Prompt represents an array of chat messages in OpenRouter format.
// This matches OpenRouter's chat completion API message format.
type Prompt []Message

// Message is the interface that all message types must implement
type Message interface {
	// The role of the message. One of RoleSystem, RoleUser, RoleAssistant, or RoleTool
	Role() string
}

// SystemMessage represents a system message
type SystemMessage struct {
	Content string `json:"content"`
}

var _ Message = &SystemMessage{}

func (m *SystemMessage) Role() string {
	return RoleSystem
}

// AssistantMessage represents an assistant message
type AssistantMessage struct {
	Content string `json:"content,omitempty"`
	// TODO: figure out whether reasoning is just for assistant,
	// or if it's a general field for all messages.
	// Coming from: https://github.com/OpenRouterTeam/ai-sdk-provider/blob/main/src/openrouter-chat-language-model.ts#L514
	Reasoning string     `json:"reasoning,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

var _ Message = &AssistantMessage{}

func (m *AssistantMessage) Role() string {
	return RoleAssistant
}

// UserMessage represents a user message
type UserMessage struct {
	Content UserMessageContent `json:"content"`
}

var _ Message = &UserMessage{}

func (m *UserMessage) Role() string {
	return RoleUser
}

// UserMessageContent represents either a string or array of content parts.
// When sending to OpenRouter's API:
// - If Parts is non-nil, it will be sent as an array of content parts (e.g. for text + images)
// - Otherwise, Text will be sent as a simple string content
// Note: Text and Parts are mutually exclusive - only one should be set at a time
type UserMessageContent struct {
	Text  string        // if content is a simple string
	Parts []ContentPart // if content is an array of parts
}

// ContentPart is the interface for different content part types
type ContentPart interface {
	Type() string
}

// TextPart represents a text content part
type TextPart struct {
	Text string `json:"text"`
}

var _ ContentPart = &TextPart{}

func (p *TextPart) Type() string {
	return ContentTypeText
}

// ImagePart represents an image content part
type ImagePart struct {
	ImageURL struct {
		URL string `json:"url"`
	} `json:"image_url"`
}

var _ ContentPart = &ImagePart{}

func (p *ImagePart) Type() string {
	return ContentTypeImageURL
}

// ToolCall represents a tool call from the assistant
type ToolCall struct {
	Type     string `json:"type"` // always "function"
	ID       string `json:"id"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
		// TODO: Docs indicate there's an optional "description" field.
		// This differs from the TypeScript provider.
	} `json:"function"`
}

// ToolMessage represents a tool message
type ToolMessage struct {
	Content    string `json:"content"`
	ToolCallID string `json:"tool_call_id"`
}

var _ Message = &ToolMessage{}

func (m *ToolMessage) Role() string {
	return RoleTool
}
