package openrouter

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/internal/openrouter/client"
	"go.jetify.com/ai/provider/internal/openrouter/codec"
)

// OpenRouterChatLanguageModel implements the chat-based language model for OpenRouter.
type OpenRouterChatLanguageModel struct {
	provider *OpenRouterProvider
	modelID  string
	settings *client.ChatSettings
}

// NewOpenRouterChatLanguageModel creates a new OpenRouter chat language model.
func NewOpenRouterChatLanguageModel(
	provider *OpenRouterProvider,
	modelID string,
	settings *client.ChatSettings,
) *OpenRouterChatLanguageModel {
	return &OpenRouterChatLanguageModel{
		provider: provider,
		modelID:  modelID,
		settings: settings,
	}
}

// SpecificationVersion returns the specification version of the language model.
func (m *OpenRouterChatLanguageModel) SpecificationVersion() string {
	return "v1"
}

// ProviderName returns the name of the provider.
func (m *OpenRouterChatLanguageModel) ProviderName() string {
	return "openrouter.chat"
}

// ModelID returns the model identifier.
func (m *OpenRouterChatLanguageModel) ModelID() string {
	return m.modelID
}

// DefaultObjectGenerationMode returns the default mode for object generation.
func (m *OpenRouterChatLanguageModel) DefaultObjectGenerationMode() api.ObjectGenerationMode {
	return api.ObjectGenerationModeTool
}

// TODO: we need to double check that each of these private structs:
// - is not a duplicate of a struct we might have already defined
// - that all the data returned by openrouter is something we are able to
//   expose in our unified API

// chatResponse represents the OpenRouter chat completion response.
type chatResponse struct {
	Choices []struct {
		Message struct {
			Role      string            `json:"role"`
			Content   *string           `json:"content,omitempty"`
			Reasoning *string           `json:"reasoning,omitempty"`
			ToolCalls []client.ToolCall `json:"tool_calls,omitempty"`
		} `json:"message"`
		LogProbs     *client.LogProbs `json:"logprobs,omitempty"`
		FinishReason string           `json:"finish_reason,omitempty"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
}

// chatStreamResponse represents a streaming response chunk.
type chatStreamResponse struct {
	ID      string `json:"id,omitempty"`
	Choices []struct {
		Delta struct {
			Role      string          `json:"role,omitempty"`
			Content   *string         `json:"content,omitempty"`
			Reasoning *string         `json:"reasoning,omitempty"`
			ToolCalls []toolCallDelta `json:"tool_calls,omitempty"`
		} `json:"delta"`
		LogProbs     *client.LogProbs `json:"logprobs,omitempty"`
		FinishReason string           `json:"finish_reason,omitempty"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage,omitempty"`
}

// toolCallDelta represents a partial tool call in a streaming response.
type toolCallDelta struct {
	Index    int    `json:"index"`
	ID       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"`
	Function struct {
		Name      string `json:"name,omitempty"`
		Arguments string `json:"arguments,omitempty"`
	} `json:"function"`
}

// buildRequestBody builds the request body for the OpenRouter API.
//
//nolint:revive // TODO: Refactor to reduce cognitive complexity (currently 41 > max 30)
func (m *OpenRouterChatLanguageModel) buildRequestBody(
	prompt []api.Message,
	options api.CallOptions,
	stream bool,
) (map[string]any, error) {
	messages, err := codec.EncodePrompt(prompt)
	if err != nil {
		return nil, fmt.Errorf("convert prompt: %w", err)
	}

	// Base request body
	body := map[string]any{
		"model":    m.modelID,
		"messages": messages,
	}

	// Add stream options if streaming
	if stream {
		body["stream"] = true
		body["stream_options"] = map[string]any{
			"include_usage": true,
		}
	}

	// Add settings if present
	if m.settings != nil {
		if m.settings.LogitBias != nil {
			body["logit_bias"] = m.settings.LogitBias
		}
		if m.settings.Logprobs.Enabled {
			body["logprobs"] = true
			if m.settings.Logprobs.TopK > 0 {
				body["top_logprobs"] = m.settings.Logprobs.TopK
			}
		}
		if m.settings.User != nil {
			body["user"] = *m.settings.User
		}
		if m.settings.ParallelToolCalls != nil {
			body["parallel_tool_calls"] = *m.settings.ParallelToolCalls
		}
		if m.settings.IncludeReasoning != nil {
			body["include_reasoning"] = *m.settings.IncludeReasoning
		}
		if len(m.settings.Models) > 0 {
			body["models"] = m.settings.Models
		}
	}

	// Handle tools from top-level CallOptions
	if len(options.Tools) > 0 {
		tools := make([]map[string]any, 0, len(options.Tools))
		for _, tool := range options.Tools {
			if funcTool, ok := tool.(api.FunctionTool); ok {
				tools = append(tools, map[string]any{
					"type": "function",
					"function": map[string]any{
						"name":        funcTool.Name,
						"description": funcTool.Description,
						"parameters":  funcTool.InputSchema,
					},
				})
			}
		}
		body["tools"] = tools
	}

	// Handle tool choice from top-level CallOptions
	if options.ToolChoice != nil {
		switch options.ToolChoice.Type {
		case "auto", "none", "required":
			body["tool_choice"] = options.ToolChoice.Type
		case "tool":
			body["tool_choice"] = map[string]any{
				"type": "function",
				"function": map[string]any{
					"name": options.ToolChoice.ToolName,
				},
			}
		}
	}

	// Handle response format from top-level CallOptions
	if options.ResponseFormat != nil && options.ResponseFormat.Type == "json" {
		body["response_format"] = map[string]any{"type": "json_object"}
		if options.ResponseFormat.Schema != nil {
			body["response_format"].(map[string]any)["schema"] = options.ResponseFormat.Schema
		}
	}

	// Add call options
	if options.MaxOutputTokens > 0 {
		body["max_tokens"] = options.MaxOutputTokens
	}
	if options.Temperature != nil {
		body["temperature"] = options.Temperature
	}
	if options.TopP != 0 {
		body["top_p"] = options.TopP
	}
	if options.FrequencyPenalty != 0 {
		body["frequency_penalty"] = options.FrequencyPenalty
	}
	if options.PresencePenalty != 0 {
		body["presence_penalty"] = options.PresencePenalty
	}
	if options.Seed != 0 {
		body["seed"] = options.Seed
	}

	return body, nil
}

// DoGenerate implements the non-streaming generation method.
func (m *OpenRouterChatLanguageModel) DoGenerate(
	ctx context.Context,
	prompt []api.Message,
	opts api.CallOptions,
) (api.Response, error) {
	requestBody, err := m.buildRequestBody(prompt, opts, false)
	if err != nil {
		return api.Response{}, err
	}

	// Make request
	resp, err := m.provider.doJSONRequest(
		ctx,
		http.MethodPost,
		"/chat/completions",
		requestBody,
		opts.Headers,
	)
	if err != nil {
		return api.Response{}, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.Response{}, fmt.Errorf("read response body: %w", err)
	}

	// Handle non-200 responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return api.Response{}, OpenRouterFailedResponseHandler(resp, body, requestBody)
	}

	// Parse response
	var response chatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return api.Response{}, api.NewJSONParseError(string(body), err)
	}

	// Ensure we have at least one choice
	if len(response.Choices) == 0 {
		return api.Response{}, api.NewNoContentGeneratedError("no choices in response")
	}

	choice := response.Choices[0]

	// Build result
	result := api.Response{
		Text:         stringValue(choice.Message.Content),
		Reasoning:    convertReasoning(choice.Message.Reasoning),
		FinishReason: codec.DecodeFinishReason(choice.FinishReason),
	}

	// Add tool calls if present
	if len(choice.Message.ToolCalls) > 0 {
		result.ToolCalls = make([]api.ToolCallBlock, len(choice.Message.ToolCalls))
		for i, tc := range choice.Message.ToolCalls {
			result.ToolCalls[i] = api.ToolCallBlock{
				ToolCallID: tc.ID,
				ToolName:   tc.Function.Name,
				Args:       json.RawMessage(tc.Function.Arguments),
			}
		}
	}

	// Add logprobs if present
	// TODO: Add logprobs support
	// if choice.LogProbs != nil {
	// 	result.LogProbs = MapOpenRouterChatLogProbsOutput(choice.LogProbs)
	// }

	return result, nil
}

// DoStream implements the streaming generation method.
//
//nolint:revive // TODO: Refactor to reduce cognitive complexity (currently 66 > max 30)
func (m *OpenRouterChatLanguageModel) DoStream(
	ctx context.Context,
	prompt []api.Message,
	opts api.CallOptions,
) (api.StreamResponse, error) {
	requestBody, err := m.buildRequestBody(prompt, opts, true)
	if err != nil {
		return api.StreamResponse{}, err
	}

	// TODO: Consider replacing custom SSE parsing logic with a dedicated SSE library
	// for more robust handling of edge cases and reconnection logic.
	// Some options:
	// - github.com/r3labs/sse
	// - github.com/donovanhide/eventsource

	// Make request
	resp, err := m.provider.doJSONRequest(
		ctx,
		http.MethodPost,
		"/chat/completions",
		requestBody,
		opts.Headers,
	)
	if err != nil {
		return api.StreamResponse{}, err
	}

	// Create sequence for events
	events := func(yield func(api.StreamEvent) bool) {
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		var toolCalls []client.ToolCall
		// TODO: Add logprobs support
		// var logprobs []api.LogProb
		// TODO: Add usage support
		// var usage api.Usage

		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var chunk chatStreamResponse
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				if !yield(api.ErrorEvent{
					Err: api.NewJSONParseError(data, err),
				}) {
					return
				}
				return
			}

			// Handle error chunks
			if len(chunk.Choices) == 0 {
				continue
			}

			choice := chunk.Choices[0]
			delta := choice.Delta

			// Handle text delta
			if delta.Content != nil {
				if !yield(api.TextDeltaEvent{
					TextDelta: *delta.Content,
				}) {
					return
				}
			}

			// Handle reasoning delta
			if delta.Reasoning != nil {
				if !yield(api.ReasoningEvent{
					TextDelta: *delta.Reasoning,
				}) {
					return
				}
			}

			// Handle tool calls
			if len(delta.ToolCalls) > 0 {
				for _, tc := range delta.ToolCalls {
					if len(toolCalls) <= tc.Index {
						toolCalls = append(toolCalls, client.ToolCall{
							ID:   tc.ID,
							Type: tc.Type,
							Function: struct {
								Name      string `json:"name"`
								Arguments string `json:"arguments"`
							}{
								Name:      tc.Function.Name,
								Arguments: "",
							},
						})
					}

					toolCall := &toolCalls[tc.Index]

					if tc.Function.Arguments != "" {
						toolCall.Function.Arguments += tc.Function.Arguments
						if !yield(api.ToolCallDeltaEvent{
							ToolCallID: toolCall.ID,
							ToolName:   toolCall.Function.Name,
							ArgsDelta:  []byte(tc.Function.Arguments),
						}) {
							return
						}

						// TODO: Add tool call support
						// // If arguments form valid JSON, emit complete tool call
						// if json.Valid([]byte(toolCall.Function.Arguments)) {
						// 	stream <- api.StreamPart{
						// 		Type:       "tool-call",
						// 		ToolCallID: toolCall.ID,
						// 		ToolName:   toolCall.Function.Name,
						// 		Args:       toolCall.Function.Arguments,
						// 	}
						// }
					}
				}
			}

			// Handle logprobs
			// TODO: Add logprobs support
			// if choice.LogProbs != nil {
			//     if mapped := MapOpenRouterChatLogProbsOutput(choice.LogProbs); len(mapped) > 0 {
			//         logprobs = append(logprobs, mapped...)
			//     }
			// }

			// Handle finish
			if choice.FinishReason != "" {
				// TODO: Add usage support
				// if chunk.Usage != nil {
				// 	usage = api.Usage{
				// 		PromptTokens:     chunk.Usage.PromptTokens,
				// 		CompletionTokens: chunk.Usage.CompletionTokens,
				// 	}
				// }

				finishReason := codec.DecodeFinishReason(choice.FinishReason)
				if !yield(api.FinishEvent{
					FinishReason: finishReason,
					// TODO: Add logprobs support and usage support
					// LogProbs:     logprobs,
					// Usage:        &usage,
				}) {
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			yield(api.ErrorEvent{
				Err: fmt.Errorf("scanner error: %w", err),
			})
		}
	}

	return api.StreamResponse{Events: events}, nil
}

// stringValue returns an empty string if the pointer is nil, otherwise returns the string value
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// convertReasoning converts a string reasoning to a slice of api.Reasoning
func convertReasoning(s *string) []api.Reasoning {
	if s == nil || *s == "" {
		return nil
	}
	return []api.Reasoning{&api.ReasoningBlock{Text: *s}}
}
