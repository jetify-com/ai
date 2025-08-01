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

type anthropicMessage struct {
	message     *anthropic.BetaMessageParam   // nil if not applicable
	systemBlock *anthropic.BetaTextBlockParam // nil if not applicable
	betas       []anthropic.AnthropicBeta
}

// EncodePrompt converts an AI SDK prompt into Anthropic's message format
func EncodePrompt(prompt []api.Message) (*AnthropicPrompt, error) {
	// Pre-process the prompt to merge messages of the same type together.
	// TODO: maybe this should move out of the provider, and be done by the framework
	// for all providers?
	prompt = mergeMessages(prompt)

	messages := make([]anthropic.BetaMessageParam, 0, len(prompt))
	var systemBlocks []anthropic.BetaTextBlockParam
	hasSeenNonSystem := false
	var betas []anthropic.AnthropicBeta

	for _, msg := range prompt {
		result, err := processMessage(msg)
		if err != nil {
			return nil, err
		}

		if result.message != nil {
			messages = append(messages, *result.message)
			hasSeenNonSystem = true
		}
		if result.systemBlock != nil {
			if hasSeenNonSystem && len(systemBlocks) > 0 {
				return nil, fmt.Errorf("multiple system messages separated by user/assistant messages are not supported")
			}
			systemBlocks = append(systemBlocks, *result.systemBlock)
		}
		betas = append(betas, result.betas...)
	}

	return &AnthropicPrompt{
		System:   systemBlocks,
		Messages: messages,
		Betas:    betas,
	}, nil
}

// processMessage handles a single message and returns the result
func processMessage(msg api.Message) (*anthropicMessage, error) {
	switch typedMsg := msg.(type) {
	case *api.SystemMessage:
		block, err := EncodeSystemMessage(typedMsg)
		if err != nil {
			return nil, err
		}
		return &anthropicMessage{
			systemBlock: &block,
		}, nil

	case api.SystemMessage:
		block, err := EncodeSystemMessage(&typedMsg)
		if err != nil {
			return nil, err
		}
		return &anthropicMessage{
			systemBlock: &block,
		}, nil

	case *api.UserMessage:
		msgParam, msgBetas, err := EncodeUserMessage(typedMsg)
		if err != nil {
			return nil, err
		}
		return &anthropicMessage{
			message: &msgParam,
			betas:   msgBetas,
		}, nil

	case api.UserMessage:
		msgParam, msgBetas, err := EncodeUserMessage(&typedMsg)
		if err != nil {
			return nil, err
		}
		return &anthropicMessage{
			message: &msgParam,
			betas:   msgBetas,
		}, nil

	case *api.AssistantMessage:
		msgParam, err := EncodeAssistantMessage(typedMsg)
		if err != nil {
			return nil, err
		}
		return &anthropicMessage{
			message: &msgParam,
		}, nil

	case api.AssistantMessage:
		msgParam, err := EncodeAssistantMessage(&typedMsg)
		if err != nil {
			return nil, err
		}
		return &anthropicMessage{
			message: &msgParam,
		}, nil

	case *api.ToolMessage:
		msgParam, err := EncodeToolMessage(typedMsg)
		if err != nil {
			return nil, err
		}
		return &anthropicMessage{
			message: &msgParam,
		}, nil

	case api.ToolMessage:
		msgParam, err := EncodeToolMessage(&typedMsg)
		if err != nil {
			return nil, err
		}
		return &anthropicMessage{
			message: &msgParam,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported message type: %T", msg)
	}
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
		Type: "text",
		Text: block.Text,
	}
	if cacheControl := getCacheControl(block); cacheControl != nil {
		param.CacheControl = *cacheControl
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
			Type: "url",
			URL:  block.URL,
		}
		param = anthropic.BetaImageBlockParam{
			Type: "image",
			Source: anthropic.BetaImageBlockParamSourceUnion{
				OfURL: &urlSource,
			},
		}
	} else if block.Data != nil {
		mimeType := block.MediaType
		if mimeType == "" {
			mimeType = "image/jpeg" // Default to JPEG if no mime type specified
		}
		base64Data := base64.StdEncoding.EncodeToString(block.Data)
		param = NewImageBlockBase64(mimeType, base64Data)
	} else {
		return anthropic.BetaImageBlockParam{}, fmt.Errorf("image block must have either URL or Data")
	}
	if cacheControl := getCacheControl(block); cacheControl != nil {
		param.CacheControl = *cacheControl
	}
	return param, nil
}

func EncodeFileBlock(block *api.FileBlock) (anthropic.BetaRequestDocumentBlockParam, []anthropic.AnthropicBeta, error) {
	if block == nil {
		return anthropic.BetaRequestDocumentBlockParam{}, []anthropic.AnthropicBeta{}, fmt.Errorf("file block cannot be nil")
	}

	var param anthropic.BetaRequestDocumentBlockParam
	param.Type = "document"
	var betas []anthropic.AnthropicBeta
	var isPDF bool

	if block.URL != "" {
		// Check if it's a PDF URL by looking at the extension or mime type
		isPDF = block.MediaType == "application/pdf"

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
				Type: "url",
				URL:  block.URL,
			}
			param.Source = anthropic.BetaRequestDocumentBlockSourceUnionParam{
				OfURL: &urlSource,
			}
		} else {
			// Plain text URL handling for non-PDFs
			textSource := anthropic.BetaPlainTextSourceParam{
				Type:      "text",
				Data:      block.URL,
				MediaType: "text/plain",
			}
			param.Source = anthropic.BetaRequestDocumentBlockSourceUnionParam{
				OfText: &textSource,
			}
		}
	} else if block.Data != nil {
		mimeType := block.MediaType
		if mimeType == "application/pdf" {
			// Base64 PDF source
			base64Data := base64.StdEncoding.EncodeToString(block.Data)
			base64Source := anthropic.BetaBase64PDFSourceParam{
				Type:      "base64",
				Data:      base64Data,
				MediaType: "application/pdf",
			}
			param.Source = anthropic.BetaRequestDocumentBlockSourceUnionParam{
				OfBase64: &base64Source,
			}
			isPDF = true
		} else if mimeType == "text/plain" || mimeType == "" {
			// Plain text source
			textData := string(block.Data)
			textSource := anthropic.BetaPlainTextSourceParam{
				Type:      "text",
				Data:      textData,
				MediaType: "text/plain",
			}
			param.Source = anthropic.BetaRequestDocumentBlockSourceUnionParam{
				OfText: &textSource,
			}
		} else {
			// Unsupported mime type
			return anthropic.BetaRequestDocumentBlockParam{}, []anthropic.AnthropicBeta{}, fmt.Errorf("unsupported mime type for file block: %s", mimeType)
		}
	} else {
		// Handle empty file block (no URL or Data)
		return anthropic.BetaRequestDocumentBlockParam{}, []anthropic.AnthropicBeta{}, fmt.Errorf("file block must have either URL or Data")
	}

	// TODO: Add support for content-type source which would allow structured content blocks
	// TODO: Add citation support: https://docs.anthropic.com/en/docs/build-with-claude/citations

	// Handle cache control if present
	if cacheControl := getCacheControl(block); cacheControl != nil {
		param.CacheControl = *cacheControl
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
			return anthropic.BetaContentBlockParamUnion{}, nil, err
		}
		return anthropic.BetaContentBlockParamUnion{
			OfText: &param,
		}, nil, nil
	case api.TextBlock:
		param, err := EncodeTextBlock(&block)
		if err != nil {
			return anthropic.BetaContentBlockParamUnion{}, nil, err
		}
		return anthropic.BetaContentBlockParamUnion{
			OfText: &param,
		}, nil, nil
	case *api.ImageBlock:
		param, err := EncodeImageBlock(block)
		if err != nil {
			return anthropic.BetaContentBlockParamUnion{}, nil, err
		}
		return anthropic.BetaContentBlockParamUnion{
			OfImage: &param,
		}, nil, nil
	case api.ImageBlock:
		param, err := EncodeImageBlock(&block)
		if err != nil {
			return anthropic.BetaContentBlockParamUnion{}, nil, err
		}
		return anthropic.BetaContentBlockParamUnion{
			OfImage: &param,
		}, nil, nil
	case *api.FileBlock:
		param, betas, err := EncodeFileBlock(block)
		if err != nil {
			return anthropic.BetaContentBlockParamUnion{}, nil, err
		}
		return anthropic.BetaContentBlockParamUnion{
			OfDocument: &param,
		}, betas, nil
	case api.FileBlock:
		param, betas, err := EncodeFileBlock(&block)
		if err != nil {
			return anthropic.BetaContentBlockParamUnion{}, nil, err
		}
		return anthropic.BetaContentBlockParamUnion{
			OfDocument: &param,
		}, betas, nil
	default:
		return anthropic.BetaContentBlockParamUnion{}, nil, fmt.Errorf("unsupported content block type: %T", block)
	}
}

func EncodeToolCallBlock(block *api.ToolCallBlock) (anthropic.BetaToolUseBlockParam, error) {
	if block == nil {
		return anthropic.BetaToolUseBlockParam{}, fmt.Errorf("tool call block cannot be nil")
	}

	param := anthropic.BetaToolUseBlockParam{
		Type:  "tool_use",
		ID:    block.ToolCallID,
		Name:  block.ToolName,
		Input: block.Args,
	}
	if cacheControl := getCacheControl(block); cacheControl != nil {
		param.CacheControl = *cacheControl
	}
	return param, nil
}

func EncodeAssistantMessage(msg *api.AssistantMessage) (anthropic.BetaMessageParam, error) {
	params := make([]anthropic.BetaContentBlockParamUnion, 0, len(msg.Content))
	for _, block := range msg.Content {
		var param anthropic.BetaContentBlockParamUnion
		switch block := block.(type) {
		case *api.TextBlock:
			textParam, err := EncodeTextBlock(block)
			if err != nil {
				return anthropic.BetaMessageParam{}, err
			}
			param = anthropic.BetaContentBlockParamUnion{
				OfText: &textParam,
			}
		case api.TextBlock:
			textParam, err := EncodeTextBlock(&block)
			if err != nil {
				return anthropic.BetaMessageParam{}, err
			}
			param = anthropic.BetaContentBlockParamUnion{
				OfText: &textParam,
			}
		case *api.ToolCallBlock:
			toolParam, err := EncodeToolCallBlock(block)
			if err != nil {
				return anthropic.BetaMessageParam{}, err
			}
			param = anthropic.BetaContentBlockParamUnion{
				OfToolUse: &toolParam,
			}
		case api.ToolCallBlock:
			toolParam, err := EncodeToolCallBlock(&block)
			if err != nil {
				return anthropic.BetaMessageParam{}, err
			}
			param = anthropic.BetaContentBlockParamUnion{
				OfToolUse: &toolParam,
			}
		case *api.ReasoningBlock:
			// TODO: in the Vercel API there is a "sendReasoning" option that needs to
			// be enabled for reasoning blocks to be sent. Do we want to add a similar
			// option?
			// See here: https://github.com/vercel/ai/blob/main/packages/anthropic/src/convert-to-anthropic-messages-prompt.ts#L257
			reasoningParam, err := EncodeReasoningBlock(block)
			if err != nil {
				return anthropic.BetaMessageParam{}, err
			}
			param = anthropic.BetaContentBlockParamUnion{
				OfThinking: &reasoningParam,
			}
		case api.ReasoningBlock:
			reasoningParam, err := EncodeReasoningBlock(&block)
			if err != nil {
				return anthropic.BetaMessageParam{}, err
			}
			param = anthropic.BetaContentBlockParamUnion{
				OfThinking: &reasoningParam,
			}
		case *api.RedactedReasoningBlock:
			redactedParam, err := EncodeRedactedReasoningBlock(block)
			if err != nil {
				return anthropic.BetaMessageParam{}, err
			}
			param = anthropic.BetaContentBlockParamUnion{
				OfRedactedThinking: &redactedParam,
			}
		case api.RedactedReasoningBlock:
			redactedParam, err := EncodeRedactedReasoningBlock(&block)
			if err != nil {
				return anthropic.BetaMessageParam{}, err
			}
			param = anthropic.BetaContentBlockParamUnion{
				OfRedactedThinking: &redactedParam,
			}
		default:
			return anthropic.BetaMessageParam{}, fmt.Errorf("unsupported assistant content block type: %T", block)
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
		Type:      "thinking",
		Thinking:  block.Text,
		Signature: block.Signature,
	}, nil
}

func EncodeRedactedReasoningBlock(block *api.RedactedReasoningBlock) (anthropic.BetaRedactedThinkingBlockParam, error) {
	if block == nil {
		return anthropic.BetaRedactedThinkingBlockParam{}, fmt.Errorf("redacted reasoning block cannot be nil")
	}
	return anthropic.BetaRedactedThinkingBlockParam{
		Type: "redacted_thinking",
		Data: block.Data,
	}, nil
}

func EncodeToolMessage(msg *api.ToolMessage) (anthropic.BetaMessageParam, error) {
	blocks := make([]anthropic.BetaContentBlockParamUnion, 0, len(msg.Content))
	for _, result := range msg.Content {
		if result.ToolCallID == "" {
			return anthropic.BetaMessageParam{}, fmt.Errorf("tool call ID cannot be empty")
		}

		block, err := encodeToolResult(result)
		if err != nil {
			return anthropic.BetaMessageParam{}, err
		}
		blocks = append(blocks, block)
	}
	return NewUserMessage(blocks...), nil
}

func encodeToolResult(result api.ToolResultBlock) (anthropic.BetaContentBlockParamUnion, error) {
	if result.Content != nil {
		return encodeToolResultContent(result)
	}
	return encodeToolResultJSON(result)
}

func encodeToolResultContent(result api.ToolResultBlock) (anthropic.BetaContentBlockParamUnion, error) {
	content := make([]anthropic.BetaToolResultBlockParamContentUnion, 0, len(result.Content))
	for _, part := range result.Content {
		param, err := encodeToolResultPart(part)
		if err != nil {
			return anthropic.BetaContentBlockParamUnion{}, err
		}
		content = append(content, param)
	}

	param := anthropic.BetaToolResultBlockParam{
		Type:      "tool_result",
		ToolUseID: result.ToolCallID,
		Content:   content,
		IsError:   anthropic.Bool(result.IsError),
	}
	if cacheControl := getCacheControl(result); cacheControl != nil {
		param.CacheControl = *cacheControl
	}
	return anthropic.BetaContentBlockParamUnion{
		OfToolResult: &param,
	}, nil
}

func encodeToolResultJSON(result api.ToolResultBlock) (anthropic.BetaContentBlockParamUnion, error) {
	resultJSON, err := json.Marshal(result.Result)
	if err != nil {
		return anthropic.BetaContentBlockParamUnion{}, fmt.Errorf("failed to marshal tool result: %v", err)
	}
	toolResultParam := NewToolResultBlock(result.ToolCallID, string(resultJSON), result.IsError)
	if cacheControl := getCacheControl(result); cacheControl != nil {
		toolResultParam.CacheControl = *cacheControl
	}
	return anthropic.BetaContentBlockParamUnion{
		OfToolResult: &toolResultParam,
	}, nil
}

func encodeToolResultPart(part api.ContentBlock) (anthropic.BetaToolResultBlockParamContentUnion, error) {
	switch block := part.(type) {
	case *api.TextBlock:
		return encodeToolResultTextPart(block)
	case api.TextBlock:
		return encodeToolResultTextPart(&block)
	case *api.ImageBlock:
		return encodeToolResultImagePart(block)
	case api.ImageBlock:
		return encodeToolResultImagePart(&block)
	default:
		return anthropic.BetaToolResultBlockParamContentUnion{}, fmt.Errorf("unsupported tool result content type: %T", block)
	}
}

func encodeToolResultTextPart(block *api.TextBlock) (anthropic.BetaToolResultBlockParamContentUnion, error) {
	textParam, err := EncodeTextBlock(block)
	if err != nil {
		return anthropic.BetaToolResultBlockParamContentUnion{}, fmt.Errorf("failed to encode text block: %v", err)
	}
	return anthropic.BetaToolResultBlockParamContentUnion{
		OfText: &textParam,
	}, nil
}

func encodeToolResultImagePart(block *api.ImageBlock) (anthropic.BetaToolResultBlockParamContentUnion, error) {
	imageParam, err := EncodeImageBlock(block)
	if err != nil {
		return anthropic.BetaToolResultBlockParamContentUnion{}, fmt.Errorf("failed to encode image block: %v", err)
	}
	return anthropic.BetaToolResultBlockParamContentUnion{
		OfImage: &imageParam,
	}, nil
}

func EncodeSystemMessage(msg *api.SystemMessage) (anthropic.BetaTextBlockParam, error) {
	block := anthropic.BetaTextBlockParam{
		Type: "text",
		Text: msg.Content,
	}
	if cacheControl := getCacheControl(msg); cacheControl != nil {
		block.CacheControl = *cacheControl
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
			Type: "ephemeral",
		}
	}
	return nil
}

func NewImageBlockBase64(mediaType string, encodedData string) anthropic.BetaImageBlockParam {
	base64Source := anthropic.BetaBase64ImageSourceParam{
		Type:      "base64",
		Data:      encodedData,
		MediaType: anthropic.BetaBase64ImageSourceMediaType(mediaType),
	}
	return anthropic.BetaImageBlockParam{
		Type: "image",
		Source: anthropic.BetaImageBlockParamSourceUnion{
			OfBase64: &base64Source,
		},
	}
}

func NewToolResultBlock(toolUseID string, content string, isError bool) anthropic.BetaToolResultBlockParam {
	return anthropic.BetaToolResultBlockParam{
		Type:      "tool_result",
		ToolUseID: toolUseID,
		Content: []anthropic.BetaToolResultBlockParamContentUnion{
			{
				OfText: &anthropic.BetaTextBlockParam{
					Text: content,
					Type: "text",
				},
			},
		},
		IsError: anthropic.Bool(isError),
	}
}

func NewToolUseBlockParam(id string, name string, input interface{}) anthropic.BetaToolUseBlockParam {
	return anthropic.BetaToolUseBlockParam{
		ID:    id,
		Input: input,
		Name:  name,
		Type:  "tool_use",
	}
}

func NewTextBlock(text string) anthropic.BetaTextBlockParam {
	return anthropic.BetaTextBlockParam{
		Text: text,
		Type: "text",
	}
}

func NewUserMessage(blocks ...anthropic.BetaContentBlockParamUnion) anthropic.BetaMessageParam {
	return anthropic.BetaMessageParam{
		Role:    "user",
		Content: blocks,
	}
}

func NewAssistantMessage(blocks ...anthropic.BetaContentBlockParamUnion) anthropic.BetaMessageParam {
	return anthropic.BetaMessageParam{
		Role:    "assistant",
		Content: blocks,
	}
}
