// TODO: This package should be internal, but leaving it public for now
// because as we transition to the new AI SDK, some of the functions in
// here are useful.

package codec

// Functions that convert AI SDK types to Anthropic types.

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"go.jetify.com/ai/api"
)

type AnthropicPrompt struct {
	System   []anthropic.BetaTextBlockParam `json:"system,omitempty"`
	Messages []anthropic.BetaMessageParam   `json:"messages"`

	Betas []anthropic.AnthropicBeta `json:"betas,omitempty"`
}

// EncodePrompt converts an AI SDK prompt into Anthropic's message format
func EncodePrompt(prompt []api.Message) (*AnthropicPrompt, error) {
	// Pre-process the prompt to merge messages of the same type together.
	prompt = mergeMessages(prompt)

	messages := make([]anthropic.BetaMessageParam, 0, len(prompt))
	var systemBlocks []anthropic.BetaTextBlockParam
	hasSeenNonSystem := false
	var betas []anthropic.AnthropicBeta

	for _, msg := range prompt {
		switch m := msg.(type) {
		case *api.SystemMessage:
			if hasSeenNonSystem && len(systemBlocks) > 0 {
				return nil, fmt.Errorf("multiple system messages separated by user/assistant messages are not supported")
			}
			block, err := EncodeSystemMessage(m)
			if err != nil {
				return nil, err
			}
			systemBlocks = append(systemBlocks, block)

		case *api.UserMessage:
			hasSeenNonSystem = true
			msg, msgBetas, err := EncodeUserMessage(m)
			if err != nil {
				return nil, err
			}
			messages = append(messages, msg)
			betas = append(betas, msgBetas...)

		case *api.AssistantMessage:
			hasSeenNonSystem = true
			msg, err := EncodeAssistantMessage(m)
			if err != nil {
				return nil, err
			}
			messages = append(messages, msg)

		case *api.ToolMessage:
			hasSeenNonSystem = true
			msg, err := EncodeToolMessage(m)
			if err != nil {
				return nil, err
			}
			messages = append(messages, msg)

		default:
			return nil, fmt.Errorf("unsupported message type: %T", msg)
		}
	}

	return &AnthropicPrompt{
		System:   systemBlocks,
		Messages: messages,
		Betas:    betas,
	}, nil
}

// TODO
// Below this line functions should become private. Temporarily making them
// public while we transition to the new AI SDK.

func EncodeUserMessage(msg *api.UserMessage) (anthropic.BetaMessageParam, []anthropic.AnthropicBeta, error) {
	params := make([]anthropic.BetaContentBlockParamUnion, 0, len(msg.Content))
	var betas []anthropic.AnthropicBeta
	for _, block := range msg.Content {
		param, partBetas, err := EncodeUserContentBlock(block)
		if err != nil {
			return anthropic.BetaMessageParam{}, nil, err
		}
		params = append(params, param)
		betas = append(betas, partBetas...)
	}
	return NewUserMessage(params...), betas, nil
}

func EncodeTextBlock(block *api.TextBlock) (anthropic.BetaTextBlockParam, error) {
	if block == nil {
		return anthropic.BetaTextBlockParam{}, fmt.Errorf("text block cannot be nil")
	}
	param := anthropic.BetaTextBlockParam{
		Type: anthropic.F(anthropic.BetaTextBlockParamTypeText),
		Text: anthropic.F(block.Text),
	}
	if cacheControl := getCacheControl(block); cacheControl != nil {
		param.CacheControl = anthropic.F(*cacheControl)
	}
	return param, nil
}

func EncodeImageBlock(block *api.ImageBlock) (anthropic.BetaImageBlockParam, error) {
	if block == nil {
		return anthropic.BetaImageBlockParam{}, fmt.Errorf("image block cannot be nil")
	}
	var param anthropic.BetaImageBlockParam
	if block.URL != "" {
		urlSource := anthropic.BetaURLImageSourceParam{
			Type: anthropic.F(anthropic.BetaURLImageSourceTypeURL),
			URL:  anthropic.F(block.URL),
		}
		param = anthropic.BetaImageBlockParam{
			Type:   anthropic.F(anthropic.BetaImageBlockParamTypeImage),
			Source: anthropic.F[anthropic.BetaImageBlockParamSourceUnion](urlSource),
		}
	} else if block.Data != nil {
		mimeType := block.MimeType
		if mimeType == "" {
			mimeType = "image/jpeg" // Default to JPEG if no mime type specified
		}
		base64Data := base64.StdEncoding.EncodeToString(block.Data)
		param = NewImageBlockBase64(mimeType, base64Data)
	} else {
		return anthropic.BetaImageBlockParam{}, fmt.Errorf("image block must have either URL or Data")
	}
	if cacheControl := getCacheControl(block); cacheControl != nil {
		param.CacheControl = anthropic.F(*cacheControl)
	}
	return param, nil
}

func EncodeFileBlock(block *api.FileBlock) (anthropic.BetaBase64PDFBlockParam, []anthropic.AnthropicBeta, error) {
	if block == nil {
		return anthropic.BetaBase64PDFBlockParam{}, []anthropic.AnthropicBeta{}, fmt.Errorf("file block cannot be nil")
	}

	var param anthropic.BetaBase64PDFBlockParam
	param.Type = anthropic.F(anthropic.BetaBase64PDFBlockTypeDocument)
	var betas []anthropic.AnthropicBeta
	var isPDF bool

	if block.URL != "" {
		// Check if it's a PDF URL by looking at the extension or mime type
		isPDF = block.MimeType == "application/pdf"

		if !isPDF && block.URL != "" {
			// Parse the URL to extract the path component
			parsedURL, err := url.Parse(block.URL)
			if err == nil {
				// Get the filename from the path and check its extension
				filename := path.Base(parsedURL.Path)
				isPDF = strings.ToLower(path.Ext(filename)) == ".pdf"
			}
		}

		if isPDF {
			// PDF URL source
			urlSource := anthropic.BetaURLPDFSourceParam{
				Type: anthropic.F(anthropic.BetaURLPDFSourceTypeURL),
				URL:  anthropic.F(block.URL),
			}
			param.Source = anthropic.F[anthropic.BetaBase64PDFBlockSourceUnionParam](urlSource)
		} else {
			// Plain text URL handling for non-PDFs
			textSource := anthropic.BetaPlainTextSourceParam{
				Type:      anthropic.F(anthropic.BetaPlainTextSourceTypeText),
				Data:      anthropic.F(block.URL),
				MediaType: anthropic.F(anthropic.BetaPlainTextSourceMediaTypeTextPlain),
			}
			param.Source = anthropic.F[anthropic.BetaBase64PDFBlockSourceUnionParam](textSource)
		}
	} else if block.Data != nil {
		mimeType := block.MimeType
		if mimeType == "application/pdf" {
			// Base64 PDF source
			base64Data := base64.StdEncoding.EncodeToString(block.Data)
			base64Source := anthropic.BetaBase64PDFSourceParam{
				Type:      anthropic.F(anthropic.BetaBase64PDFSourceTypeBase64),
				Data:      anthropic.F(base64Data),
				MediaType: anthropic.F(anthropic.BetaBase64PDFSourceMediaTypeApplicationPDF),
			}
			param.Source = anthropic.F[anthropic.BetaBase64PDFBlockSourceUnionParam](base64Source)
			isPDF = true
		} else if mimeType == "text/plain" || mimeType == "" {
			// Plain text source
			textData := string(block.Data)
			textSource := anthropic.BetaPlainTextSourceParam{
				Type:      anthropic.F(anthropic.BetaPlainTextSourceTypeText),
				Data:      anthropic.F(textData),
				MediaType: anthropic.F(anthropic.BetaPlainTextSourceMediaTypeTextPlain),
			}
			param.Source = anthropic.F[anthropic.BetaBase64PDFBlockSourceUnionParam](textSource)
		} else {
			// Unsupported mime type
			return anthropic.BetaBase64PDFBlockParam{}, []anthropic.AnthropicBeta{}, fmt.Errorf("unsupported mime type for file block: %s", mimeType)
		}
	} else {
		// Handle empty file block (no URL or Data)
		return anthropic.BetaBase64PDFBlockParam{}, []anthropic.AnthropicBeta{}, fmt.Errorf("file block must have either URL or Data")
	}

	// TODO: Add support for content-type source which would allow structured content blocks
	// TODO: Add citation support: https://docs.anthropic.com/en/docs/build-with-claude/citations

	// Handle cache control if present
	if cacheControl := getCacheControl(block); cacheControl != nil {
		param.CacheControl = anthropic.F(*cacheControl)
	}

	// Add PDF beta if needed
	if isPDF {
		betas = append(betas, anthropic.AnthropicBetaPDFs2024_09_25)
	}

	return param, betas, nil
}

func EncodeUserContentBlock(block api.ContentBlock) (anthropic.BetaContentBlockParamUnion, []anthropic.AnthropicBeta, error) {
	switch block := block.(type) {
	case *api.TextBlock:
		param, err := EncodeTextBlock(block)
		if err != nil {
			return nil, nil, err
		}
		return param, nil, nil
	case *api.ImageBlock:
		param, err := EncodeImageBlock(block)
		if err != nil {
			return nil, nil, err
		}
		return param, nil, nil
	case *api.FileBlock:
		param, betas, err := EncodeFileBlock(block)
		if err != nil {
			return nil, nil, err
		}
		return param, betas, nil
	default:
		return nil, nil, fmt.Errorf("unsupported content block type: %T", block)
	}
}

func EncodeToolCallBlock(block *api.ToolCallBlock) (anthropic.BetaToolUseBlockParam, error) {
	if block == nil {
		return anthropic.BetaToolUseBlockParam{}, fmt.Errorf("tool call block cannot be nil")
	}

	param := anthropic.BetaToolUseBlockParam{
		Type:  anthropic.F(anthropic.BetaToolUseBlockParamTypeToolUse),
		ID:    anthropic.F(block.ToolCallID),
		Name:  anthropic.F(block.ToolName),
		Input: anthropic.F[any](block.Args),
	}
	if cacheControl := getCacheControl(block); cacheControl != nil {
		param.CacheControl = anthropic.F(*cacheControl)
	}
	return param, nil
}

func EncodeAssistantMessage(msg *api.AssistantMessage) (anthropic.BetaMessageParam, error) {
	params := make([]anthropic.BetaContentBlockParamUnion, 0, len(msg.Content))
	for _, block := range msg.Content {
		var param anthropic.BetaContentBlockParamUnion
		var err error
		switch block := block.(type) {
		case *api.TextBlock:
			param, err = EncodeTextBlock(block)
		case *api.ToolCallBlock:
			toolParam, err := EncodeToolCallBlock(block)
			if err != nil {
				return anthropic.BetaMessageParam{}, err
			}
			param = toolParam
		case *api.ReasoningBlock:
			// TODO: in the Vercel API there is a "sendReasoning" option that needs to
			// be enabled for reasoning blocks to be sent. Do we want to add a similar
			// option?
			// See here: https://github.com/vercel/ai/blob/main/packages/anthropic/src/convert-to-anthropic-messages-prompt.ts#L257
			reasoningParam, err := EncodeReasoningBlock(block)
			if err != nil {
				return anthropic.BetaMessageParam{}, err
			}
			param = reasoningParam
		case *api.RedactedReasoningBlock:
			redactedParam, err := EncodeRedactedReasoningBlock(block)
			if err != nil {
				return anthropic.BetaMessageParam{}, err
			}
			param = redactedParam
		default:
			return anthropic.BetaMessageParam{}, fmt.Errorf("unsupported assistant content block type: %T", block)
		}
		if err != nil {
			return anthropic.BetaMessageParam{}, err
		}
		params = append(params, param)
	}
	return NewAssistantMessage(params...), nil
}

func EncodeReasoningBlock(block *api.ReasoningBlock) (anthropic.BetaThinkingBlockParam, error) {
	if block == nil {
		return anthropic.BetaThinkingBlockParam{}, fmt.Errorf("reasoning block cannot be nil")
	}
	return anthropic.BetaThinkingBlockParam{
		Type:      anthropic.F(anthropic.BetaThinkingBlockParamTypeThinking),
		Thinking:  anthropic.F(block.Text),
		Signature: anthropic.F(block.Signature),
	}, nil
}

func EncodeRedactedReasoningBlock(block *api.RedactedReasoningBlock) (anthropic.BetaRedactedThinkingBlockParam, error) {
	if block == nil {
		return anthropic.BetaRedactedThinkingBlockParam{}, fmt.Errorf("redacted reasoning block cannot be nil")
	}
	return anthropic.BetaRedactedThinkingBlockParam{
		Type: anthropic.F(anthropic.BetaRedactedThinkingBlockParamTypeRedactedThinking),
		Data: anthropic.F(block.Data),
	}, nil
}

func EncodeToolMessage(msg *api.ToolMessage) (anthropic.BetaMessageParam, error) {
	blocks := make([]anthropic.BetaContentBlockParamUnion, 0, len(msg.Content))
	for _, result := range msg.Content {
		if result.ToolCallID == "" {
			return anthropic.BetaMessageParam{}, fmt.Errorf("tool call ID cannot be empty")
		}

		var block anthropic.BetaContentBlockParamUnion
		if result.Content != nil {
			// Handle rich content if present
			content := make([]anthropic.BetaToolResultBlockParamContentUnion, 0, len(result.Content))
			for _, part := range result.Content {
				var param anthropic.BetaToolResultBlockParamContentUnion
				var err error

				switch block := part.(type) {
				case *api.TextBlock:
					param, err = EncodeTextBlock(block)
				case *api.ImageBlock:
					param, err = EncodeImageBlock(block)
				default:
					return anthropic.BetaMessageParam{}, fmt.Errorf("unsupported tool result content type: %T", block)
				}
				if err != nil {
					return anthropic.BetaMessageParam{}, fmt.Errorf("failed to encode %T block: %v", part, err)
				}
				content = append(content, param)
			}

			param := anthropic.BetaToolResultBlockParam{
				Type:      anthropic.F(anthropic.BetaToolResultBlockParamTypeToolResult),
				ToolUseID: anthropic.F(result.ToolCallID),
				Content:   anthropic.F(content),
				IsError:   anthropic.F(result.IsError),
			}
			if cacheControl := getCacheControl(result); cacheControl != nil {
				param.CacheControl = anthropic.F(*cacheControl)
			}
			block = param
		} else {
			// Fallback to JSON encoding the Result field
			resultJSON, err := json.Marshal(result.Result)
			if err != nil {
				return anthropic.BetaMessageParam{}, fmt.Errorf("failed to marshal tool result: %v", err)
			}
			block = NewToolResultBlock(result.ToolCallID, string(resultJSON), result.IsError)
			if cacheControl := getCacheControl(result); cacheControl != nil {
				param := block.(anthropic.BetaToolResultBlockParam)
				param.CacheControl = anthropic.F(*cacheControl)
				block = param
			}
		}
		blocks = append(blocks, block)
	}
	return NewUserMessage(blocks...), nil
}

func EncodeSystemMessage(msg *api.SystemMessage) (anthropic.BetaTextBlockParam, error) {
	block := anthropic.BetaTextBlockParam{
		Type: anthropic.F(anthropic.BetaTextBlockParamTypeText),
		Text: anthropic.F(msg.Content),
	}
	if cacheControl := getCacheControl(msg); cacheControl != nil {
		block.CacheControl = anthropic.F(*cacheControl)
	}
	return block, nil
}

// getCacheControl extracts cache control settings from provider metadata
func getCacheControl(source api.MetadataSource) *anthropic.BetaCacheControlEphemeralParam {
	metadata := GetMetadata(source)
	if metadata == nil {
		return nil
	}

	if metadata.CacheControl == "ephemeral" {
		return &anthropic.BetaCacheControlEphemeralParam{
			Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
		}
	}
	return nil
}

func NewImageBlockBase64(mediaType string, encodedData string) anthropic.BetaImageBlockParam {
	return anthropic.BetaImageBlockParam{
		Type: anthropic.F(anthropic.BetaImageBlockParamTypeImage),
		Source: anthropic.F(anthropic.BetaImageBlockParamSourceUnion(anthropic.BetaImageBlockParamSource{
			Type:      anthropic.F(anthropic.BetaImageBlockParamSourceTypeBase64),
			Data:      anthropic.F(encodedData),
			MediaType: anthropic.F(anthropic.BetaImageBlockParamSourceMediaType(mediaType)),
		})),
	}
}

func NewToolResultBlock(toolUseID string, content string, isError bool) anthropic.BetaToolResultBlockParam {
	return anthropic.BetaToolResultBlockParam{
		Type:      anthropic.F(anthropic.BetaToolResultBlockParamTypeToolResult),
		ToolUseID: anthropic.F(toolUseID),
		Content:   anthropic.F([]anthropic.BetaToolResultBlockParamContentUnion{NewTextBlock(content)}),
		IsError:   anthropic.F(isError),
	}
}

func NewToolUseBlockParam(id string, name string, input interface{}) anthropic.BetaToolUseBlockParam {
	return anthropic.BetaToolUseBlockParam{
		ID:    anthropic.F(id),
		Input: anthropic.F(input),
		Name:  anthropic.F(name),
		Type:  anthropic.F(anthropic.BetaToolUseBlockParamTypeToolUse),
	}
}

func NewTextBlock(text string) anthropic.BetaTextBlockParam {
	return anthropic.BetaTextBlockParam{
		Text: anthropic.F(text),
		Type: anthropic.F(anthropic.BetaTextBlockParamTypeText),
	}
}

func NewUserMessage(blocks ...anthropic.BetaContentBlockParamUnion) anthropic.BetaMessageParam {
	return anthropic.BetaMessageParam{
		Role:    anthropic.F(anthropic.BetaMessageParamRoleUser),
		Content: anthropic.F(blocks),
	}
}

func NewAssistantMessage(blocks ...anthropic.BetaContentBlockParamUnion) anthropic.BetaMessageParam {
	return anthropic.BetaMessageParam{
		Role:    anthropic.F(anthropic.BetaMessageParamRoleAssistant),
		Content: anthropic.F(blocks),
	}
}
