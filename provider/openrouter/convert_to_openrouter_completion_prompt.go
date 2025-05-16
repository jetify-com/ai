package openrouter

import (
	"strings"

	"go.jetify.com/ai/api"
)

type InputFormat string

const (
	InputFormatPrompt   InputFormat = "prompt"
	InputFormatMessages InputFormat = "messages"
)

type CompletionPromptOptions struct {
	Prompt      []api.Message
	InputFormat InputFormat
	User        string // defaults to "user" if empty
	Assistant   string // defaults to "assistant" if empty
}

// ConvertToOpenRouterCompletionPrompt converts an AI SDK prompt into OpenRouter's completion format.
// It returns the formatted prompt string and optional stop sequences.
func ConvertToOpenRouterCompletionPrompt(opts CompletionPromptOptions) (string, []string, error) {
	if opts.User == "" {
		opts.User = "user"
	}
	if opts.Assistant == "" {
		opts.Assistant = "assistant"
	}

	// Handle direct prompt case
	if opts.InputFormat == InputFormatPrompt {
		if text, ok := isDirectUserPrompt(opts.Prompt); ok {
			return text, nil, nil
		}
	}

	var b strings.Builder

	// Handle system message prefix if present
	if len(opts.Prompt) > 0 {
		if sys, ok := opts.Prompt[0].(*api.SystemMessage); ok {
			b.WriteString(sys.Content)
			b.WriteString("\n\n")
			opts.Prompt = opts.Prompt[1:]
		}
	}

	// Process remaining messages
	for _, msg := range opts.Prompt {
		if err := writeMessage(&b, msg, opts.User, opts.Assistant); err != nil {
			return "", nil, err
		}
	}

	// Add final assistant prefix
	b.WriteString(opts.Assistant)
	b.WriteString(":\n")

	return b.String(), []string{"\n" + opts.User + ":"}, nil
}

// isDirectUserPrompt checks if the prompt is a single user message with a single text block
func isDirectUserPrompt(prompt []api.Message) (string, bool) {
	if len(prompt) != 1 {
		return "", false
	}

	um, ok := prompt[0].(*api.UserMessage)
	if !ok {
		return "", false
	}

	if len(um.Content) != 1 {
		return "", false
	}

	textBlock, ok := um.Content[0].(*api.TextBlock)
	if !ok {
		return "", false
	}

	return textBlock.Text, true
}

// writeMessage formats and writes a single message to the string builder
func writeMessage(b *strings.Builder, msg api.Message, user, assistant string) error {
	switch m := msg.(type) {
	case *api.SystemMessage:
		return api.NewInvalidPromptError(msg, "unexpected system message in prompt", nil)

	case *api.UserMessage:
		text, err := gatherUserText(m)
		if err != nil {
			return err
		}
		b.WriteString(user)
		b.WriteString(":\n")
		b.WriteString(text)
		b.WriteString("\n\n")

	case *api.AssistantMessage:
		text, err := gatherAssistantText(m)
		if err != nil {
			return err
		}
		b.WriteString(assistant)
		b.WriteString(":\n")
		b.WriteString(text)
		b.WriteString("\n\n")

	case *api.ToolMessage:
		return api.NewUnsupportedFunctionalityError("tool messages", "")

	default:
		return api.NewInvalidPromptError(msg, "unknown message type", nil)
	}

	return nil
}

// gatherUserText collects text from user message blocks, rejecting unsupported content
func gatherUserText(msg *api.UserMessage) (string, error) {
	var b strings.Builder

	for _, block := range msg.Content {
		switch block := block.(type) {
		case *api.TextBlock:
			b.WriteString(block.Text)
		case *api.ImageBlock:
			return "", api.NewUnsupportedFunctionalityError("images", "")
		case *api.FileBlock:
			return "", api.NewUnsupportedFunctionalityError("file attachments", "")
		default:
			return "", api.NewUnsupportedFunctionalityError("unknown content type", "")
		}
	}

	return b.String(), nil
}

// gatherAssistantText collects text from assistant message blocks, rejecting tool calls
func gatherAssistantText(msg *api.AssistantMessage) (string, error) {
	var b strings.Builder

	for _, block := range msg.Content {
		switch block := block.(type) {
		case *api.TextBlock:
			b.WriteString(block.Text)
		case *api.ToolCallBlock:
			return "", api.NewUnsupportedFunctionalityError("tool-call messages", "")
		default:
			return "", api.NewUnsupportedFunctionalityError("unknown content type", "")
		}
	}

	return b.String(), nil
}
