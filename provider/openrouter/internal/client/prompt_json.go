package client

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON for SystemMessage
func (m *SystemMessage) UnmarshalJSON(data []byte) error {
	type Alias SystemMessage
	aux := struct {
		Role string `json:"role"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Role != RoleSystem {
		return fmt.Errorf("invalid role for SystemMessage: %s", aux.Role)
	}
	return nil
}

// UnmarshalJSON for AssistantMessage
func (m *AssistantMessage) UnmarshalJSON(data []byte) error {
	type Alias AssistantMessage
	aux := struct {
		Role string `json:"role"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Role != RoleAssistant {
		return fmt.Errorf("invalid role for AssistantMessage: %s", aux.Role)
	}
	return nil
}

// UnmarshalJSON for UserMessage
func (m *UserMessage) UnmarshalJSON(data []byte) error {
	type Alias UserMessage
	aux := struct {
		Role string `json:"role"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Role != RoleUser {
		return fmt.Errorf("invalid role for UserMessage: %s", aux.Role)
	}
	return nil
}

// UnmarshalJSON for ToolMessage
func (m *ToolMessage) UnmarshalJSON(data []byte) error {
	type Alias ToolMessage
	aux := struct {
		Role string `json:"role"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Role != RoleTool {
		return fmt.Errorf("invalid role for ToolMessage: %s", aux.Role)
	}
	return nil
}

// Add UnmarshalJSON methods for the content part types
func (p *TextPart) UnmarshalJSON(data []byte) error {
	type Alias TextPart
	aux := struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Type != ContentTypeText {
		return fmt.Errorf("invalid type for TextPart: %s", aux.Type)
	}
	return nil
}

func (p *ImagePart) UnmarshalJSON(data []byte) error {
	type Alias ImagePart
	aux := struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Type != ContentTypeImageURL {
		return fmt.Errorf("invalid type for ImagePart: %s", aux.Type)
	}
	return nil
}

// MarshalJSON for SystemMessage
func (m *SystemMessage) MarshalJSON() ([]byte, error) {
	type Alias SystemMessage
	return json.Marshal(struct {
		Role string `json:"role"`
		*Alias
	}{
		Role:  RoleSystem,
		Alias: (*Alias)(m),
	})
}

// MarshalJSON for AssistantMessage
func (m *AssistantMessage) MarshalJSON() ([]byte, error) {
	type Alias AssistantMessage
	return json.Marshal(struct {
		Role string `json:"role"`
		*Alias
	}{
		Role:  RoleAssistant,
		Alias: (*Alias)(m),
	})
}

// MarshalJSON for UserMessage
func (m *UserMessage) MarshalJSON() ([]byte, error) {
	type Alias UserMessage
	return json.Marshal(struct {
		Role string `json:"role"`
		*Alias
	}{
		Role:  RoleUser,
		Alias: (*Alias)(m),
	})
}

// MarshalJSON for ToolMessage
func (m *ToolMessage) MarshalJSON() ([]byte, error) {
	type Alias ToolMessage
	return json.Marshal(struct {
		Role string `json:"role"`
		*Alias
	}{
		Role:  RoleTool,
		Alias: (*Alias)(m),
	})
}

// MarshalJSON for TextPart
func (p *TextPart) MarshalJSON() ([]byte, error) {
	type Alias TextPart
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  ContentTypeText,
		Alias: (*Alias)(p),
	})
}

// MarshalJSON for ImagePart
func (p *ImagePart) MarshalJSON() ([]byte, error) {
	type Alias ImagePart
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  ContentTypeImageURL,
		Alias: (*Alias)(p),
	})
}

// MarshalMessage marshals a Message interface into JSON bytes
func MarshalMessage(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}

// UnmarshalMessage unmarshals JSON bytes into the appropriate Message type
func UnmarshalMessage(data []byte) (Message, error) {
	// First unmarshal just the role to determine the message type
	var roleCheck struct {
		Role string `json:"role"`
	}
	if err := json.Unmarshal(data, &roleCheck); err != nil {
		return nil, err
	}

	// Create and unmarshal into the appropriate type based on role
	var msg Message
	switch roleCheck.Role {
	case RoleSystem:
		msg = &SystemMessage{}
	case RoleUser:
		msg = &UserMessage{}
	case RoleAssistant:
		msg = &AssistantMessage{}
	case RoleTool:
		msg = &ToolMessage{}
	default:
		return nil, fmt.Errorf("unknown message role: %s", roleCheck.Role)
	}

	if err := json.Unmarshal(data, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

// MarshalJSON implements custom JSON marshaling for user content
func (c *UserMessageContent) MarshalJSON() ([]byte, error) {
	// If there are parts, marshal as an array of parts
	if c.Parts != nil {
		return json.Marshal(c.Parts)
	}

	// Otherwise marshal as string (empty string will be marshaled as "")
	return json.Marshal(c.Text)
}

// Add this function to handle unmarshaling of content parts
func unmarshalContentPart(data []byte) (ContentPart, error) {
	var typeCheck struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &typeCheck); err != nil {
		return nil, err
	}

	switch typeCheck.Type {
	case ContentTypeText:
		var text TextPart
		if err := json.Unmarshal(data, &text); err != nil {
			return nil, err
		}
		return &text, nil
	case ContentTypeImageURL:
		var image ImagePart
		if err := json.Unmarshal(data, &image); err != nil {
			return nil, err
		}
		return &image, nil
	default:
		return nil, fmt.Errorf("unknown content part type: %s", typeCheck.Type)
	}
}

// UnmarshalJSON for UserMessageContent to handle both string and array cases
func (c *UserMessageContent) UnmarshalJSON(data []byte) error {
	// Try unmarshaling as string first
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		c.Text = text
		c.Parts = nil
		return nil
	}

	// If that fails, try as array of content parts
	var rawParts []json.RawMessage
	if err := json.Unmarshal(data, &rawParts); err != nil {
		return err
	}

	c.Text = ""
	c.Parts = make([]ContentPart, len(rawParts))
	for i, raw := range rawParts {
		part, err := unmarshalContentPart(raw)
		if err != nil {
			return err
		}
		c.Parts[i] = part
	}
	return nil
}
