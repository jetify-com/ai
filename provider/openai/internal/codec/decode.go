package codec

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openai/openai-go/responses"
	"go.jetify.com/ai/api"
)

// responseContent holds the parsed content from an OpenAI message
type responseContent struct {
	Content  []api.ContentBlock
	HasTools bool
}

// DecodeResponse converts an OpenAI Response to the AI SDK Response type
func DecodeResponse(msg *responses.Response) (api.Response, error) {
	if msg == nil {
		return api.Response{
			Content:  []api.ContentBlock{},
			Warnings: []api.CallWarning{},
		}, nil
	}

	content, err := decodeContent(msg)
	if err != nil {
		return api.Response{}, err
	}

	// Create the response with the extracted fields
	resp := api.Response{
		Content:          content.Content,
		Usage:            decodeUsage(msg.Usage),
		ProviderMetadata: decodeProviderMetadata(msg),
		Warnings:         []api.CallWarning{},
	}

	resp.FinishReason = decodeFinishReason(msg.IncompleteDetails.Reason, content.HasTools)

	return resp, nil
}

// decodeContent processes all content from the response in sequential order
func decodeContent(msg *responses.Response) (responseContent, error) {
	content := responseContent{
		Content: []api.ContentBlock{},
	}

	if msg == nil {
		return content, nil
	}

	for _, outputItem := range msg.Output {
		switch outputItem.Type {
		case "function_call", "file_search_call", "web_search_call", "computer_call":
			toolCall, err := decodeToolCall(outputItem)
			if err != nil {
				return responseContent{}, fmt.Errorf("failed to decode tool call: %w", err)
			}
			content.Content = append(content.Content, &toolCall)
			content.HasTools = true
		case "message":
			message := outputItem.AsMessage()
			for _, contentPart := range message.Content {
				// TODO: handle refusals.
				if contentPart.Type != "output_text" {
					continue
				}

				textOutput := contentPart.AsOutputText()

				// Add text block first (if non-empty)
				if textOutput.Text != "" {
					content.Content = append(content.Content, &api.TextBlock{
						Text: textOutput.Text,
					})
				}

				// Add source blocks immediately after the text
				sourceBlocks := decodeAnnotations(textOutput.Annotations)
				content.Content = append(content.Content, sourceBlocks...)
			}
		case "reasoning":
			reasoning, err := decodeReasoning(outputItem)
			if err != nil {
				return responseContent{}, fmt.Errorf("failed to decode reasoning: %w", err)
			}
			if reasoning != nil {
				content.Content = append(content.Content, reasoning)
			}
		default:
			return responseContent{}, fmt.Errorf("unknown output item type: %s", outputItem.Type)
		}
	}

	return content, nil
}

// decodeReasoning processes a reasoning output item and returns a reasoning block
func decodeReasoning(item responses.ResponseOutputItemUnion) (api.Reasoning, error) {
	if item.Type != "reasoning" {
		return nil, fmt.Errorf("unexpected item type for reasoning: %s", item.Type)
	}

	reasoningItem := item.AsReasoning()
	if len(reasoningItem.Summary) == 0 {
		return nil, fmt.Errorf("reasoning item has no summary")
	}

	// For now, we'll concatenate all summary texts with newlines if there are multiple
	// Another option would be to return a slice of ReasoningBlocks.
	// TODO: decide on the best approach.
	var texts []string
	for i, summary := range reasoningItem.Summary {
		if summary.Text == "" {
			return nil, fmt.Errorf("empty text in reasoning summary at index %d", i)
		}
		texts = append(texts, summary.Text)
	}

	if len(texts) == 0 {
		return nil, fmt.Errorf("no valid text found in reasoning summaries")
	}

	return &api.ReasoningBlock{
		Text: strings.Join(texts, "\n"),
	}, nil
}

// decodeFunctionCall processes a function call output item
func decodeFunctionCall(functionCall responses.ResponseFunctionToolCall) (api.ToolCallBlock, error) {
	if functionCall.Name == "" {
		return api.ToolCallBlock{}, fmt.Errorf("function call missing name")
	}
	if functionCall.CallID == "" {
		return api.ToolCallBlock{}, fmt.Errorf("function call missing call_id")
	}
	args := json.RawMessage(functionCall.Arguments)
	return api.ToolCallBlock{
		ToolCallID: functionCall.CallID,
		ToolName:   functionCall.Name,
		Args:       args,
	}, nil
}

// decodeFileSearchCall processes a file search call output item
func decodeFileSearchCall(fileSearch responses.ResponseFileSearchToolCall) (api.ToolCallBlock, error) {
	return api.ToolCallBlock{
		ToolCallID: fileSearch.ID,
		ToolName:   "openai.file_search",
		Args:       json.RawMessage(fileSearch.RawJSON()),
	}, nil
}

// decodeWebSearchCall processes a web search call output item
func decodeWebSearchCall(webSearch responses.ResponseFunctionWebSearch) (api.ToolCallBlock, error) {
	return api.ToolCallBlock{
		ToolCallID: webSearch.ID,
		ToolName:   "openai.web_search_preview",
		Args:       json.RawMessage(webSearch.RawJSON()),
	}, nil
}

// decodeComputerCall processes a computer call output item
func decodeComputerCall(computerCall responses.ResponseComputerToolCall) (api.ToolCallBlock, error) {
	// Convert safety checks to our internal type - initialize as empty slice
	safetyChecks := make([]ComputerSafetyCheck, 0, len(computerCall.PendingSafetyChecks))
	for _, check := range computerCall.PendingSafetyChecks {
		safetyChecks = append(safetyChecks, ComputerSafetyCheck{
			ID:      check.ID,
			Code:    check.Code,
			Message: check.Message,
		})
	}

	// Create provider metadata with safety checks
	metadata := &Metadata{
		ComputerSafetyChecks: safetyChecks,
	}

	// Create tool call block with provider metadata
	return api.ToolCallBlock{
		ToolCallID:       computerCall.ID,
		ToolName:         "openai.computer_use_preview",
		Args:             json.RawMessage(computerCall.RawJSON()),
		ProviderMetadata: api.NewProviderMetadata(map[string]any{"openai": metadata}),
	}, nil
}

// decodeToolCall processes a tool call output items
// Note that there are several types of tool calls provided by OpenAI.
func decodeToolCall(item responses.ResponseOutputItemUnion) (api.ToolCallBlock, error) {
	switch item.Type {
	case "function_call":
		return decodeFunctionCall(item.AsFunctionCall())
	case "file_search_call":
		return decodeFileSearchCall(item.AsFileSearchCall())
	case "web_search_call":
		return decodeWebSearchCall(item.AsWebSearchCall())
	case "computer_call":
		return decodeComputerCall(item.AsComputerCall())
	default:
		return api.ToolCallBlock{}, fmt.Errorf("unknown tool call type: %s", item.Type)
	}
}

// decodeAnnotations converts annotations to source blocks
func decodeAnnotations(annotations []responses.ResponseOutputTextAnnotationUnion) []api.ContentBlock {
	var sourceBlocks []api.ContentBlock
	for i, annotation := range annotations {
		if annotation.Type != "url_citation" {
			// TODO: handle other annotation types: file_citation, file_path.
			continue
		}

		urlCitation := annotation.AsURLCitation()
		sourceBlocks = append(sourceBlocks, &api.SourceBlock{
			ID:    fmt.Sprintf("source-%d", i),
			URL:   urlCitation.URL,
			Title: urlCitation.Title,
		})
	}
	return sourceBlocks
}

// decodeUsage converts an OpenAI ResponseUsage to API SDK Usage
func decodeUsage(usage responses.ResponseUsage) api.Usage {
	totalTokens := usage.TotalTokens
	if totalTokens == 0 {
		totalTokens = usage.InputTokens + usage.OutputTokens
	}
	return api.Usage{
		InputTokens:       int(usage.InputTokens),
		OutputTokens:      int(usage.OutputTokens),
		TotalTokens:       int(totalTokens),
		ReasoningTokens:   int(usage.OutputTokensDetails.ReasoningTokens),
		CachedInputTokens: int(usage.InputTokensDetails.CachedTokens),
	}
}

// decodeFinishReason converts an OpenAI response status to an AI SDK FinishReason type.
func decodeFinishReason(incompleteReason string, hasToolCalls bool) api.FinishReason {
	// Determine finish reason based on incomplete details reason
	switch incompleteReason {
	case "max_output_tokens":
		return api.FinishReasonLength
	case "content_filter":
		return api.FinishReasonContentFilter
	case "": // Empty reason means normal completion
		if hasToolCalls {
			return api.FinishReasonToolCalls
		}
		return api.FinishReasonStop
	default:
		if hasToolCalls {
			return api.FinishReasonToolCalls
		}
		return api.FinishReasonUnknown
	}
}

// decodeProviderMetadata extracts OpenAI-specific metadata
func decodeProviderMetadata(msg *responses.Response) *api.ProviderMetadata {
	return api.NewProviderMetadata(map[string]any{
		"openai": &Metadata{
			ResponseID: msg.ID,
			Usage: Usage{
				InputTokens:           int(msg.Usage.InputTokens),
				OutputTokens:          int(msg.Usage.OutputTokens),
				InputCachedTokens:     int(msg.Usage.InputTokensDetails.CachedTokens),
				OutputReasoningTokens: int(msg.Usage.OutputTokensDetails.ReasoningTokens),
			},
		},
	})
}
