package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/aitesting"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/openai/internal/codec"
	"go.jetify.com/pkg/httpmock"
	"go.jetify.com/pkg/pointer"
	"go.jetify.com/sse"
)

var standardTools = []api.ToolDefinition{
	&api.FunctionTool{
		Name: "weather",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"location": {Type: "string"},
			},
			Required:             []string{"location"},
			AdditionalProperties: api.FalseSchema(),
		},
	},
	&api.FunctionTool{
		Name: "cityAttractions",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"city": {Type: "string"},
			},
			Required:             []string{"city"},
			AdditionalProperties: api.FalseSchema(),
		},
	},
}

func TestGenerate(t *testing.T) {
	standardPrompt := []api.Message{
		&api.UserMessage{
			Content: api.ContentFromText("Hello"),
		},
	}

	// Standard OpenAI response body used across tests
	standardResponseBody := `{
		"id": "resp_67c97c0203188190a025beb4a75242bc",
		"object": "response",
		"created_at": 1741257730,
		"status": "completed",
		"error": null,
		"incomplete_details": null,
		"input": [],
		"instructions": null,
		"max_output_tokens": null,
		"model": "gpt-4o-2024-07-18",
		"output": [
			{
				"id": "msg_67c97c02656c81908e080dfdf4a03cd1",
				"type": "message",
				"status": "completed",
				"role": "assistant",
				"content": [
					{
						"type": "output_text",
						"text": "answer text",
						"annotations": []
					}
				]
			}
		],
		"parallel_tool_calls": true,
		"previous_response_id": null,
		"reasoning": {
			"effort": null,
			"summary": null
		},
		"store": true,
		"temperature": 1,
		"text": {
			"format": {
				"type": "text"
			}
		},
		"tool_choice": "auto",
		"tools": [],
		"top_p": 1,
		"truncation": "disabled",
		"usage": {
			"input_tokens": 345,
			"input_tokens_details": {
				"cached_tokens": 234
			},
			"output_tokens": 538,
			"output_tokens_details": {
				"reasoning_tokens": 123
			},
			"total_tokens": 883
		},
		"user": null,
		"metadata": {}
	}`

	standardExchange := []httpmock.Exchange{
		{
			Request: httpmock.Request{
				Method: http.MethodPost,
				Path:   "/responses",
			},
			Response: httpmock.Response{
				StatusCode: http.StatusOK,
				Body:       standardResponseBody,
			},
		},
	}

	// This test demonstrates various response field validations using aitesting.ResponseContains.
	// When creating new tests, you can omit fields you don't care about in expectedResp.
	// For example:
	// - Only check text: expectedResp: &api.Response{Text: "example"}
	// - Only check usage: expectedResp: &api.Response{Usage: api.Usage{PromptTokens: 10}}
	// - Only check finish reason: expectedResp: &api.Response{FinishReason: api.FinishReasonStop}
	// aitesting.ResponseContains only checks fields that are specified in expectedResp.
	tests := []struct {
		name         string
		modelID      string          // Optional: override default model
		options      api.CallOptions // Optional: custom call options
		prompt       []api.Message
		exchanges    []httpmock.Exchange
		wantErr      bool
		expectedResp *api.Response // Expected response fields
		skip         bool          // Skip this test if true
	}{
		{
			name:      "response contains text",
			modelID:   "gpt-4o",
			prompt:    standardPrompt,
			exchanges: standardExchange,
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
			},
		},
		{
			name:      "response contains usage information",
			modelID:   "gpt-4o",
			prompt:    standardPrompt,
			exchanges: standardExchange,
			expectedResp: &api.Response{
				Usage: api.Usage{
					InputTokens:       345,
					OutputTokens:      538,
					TotalTokens:       883,
					ReasoningTokens:   123,
					CachedInputTokens: 234,
				},
			},
		},
		{
			name:      "response contains provider metadata",
			modelID:   "gpt-4o",
			prompt:    standardPrompt,
			exchanges: standardExchange,
			expectedResp: &api.Response{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"openai": &Metadata{
						ResponseID: "resp_67c97c0203188190a025beb4a75242bc",
						Usage: codec.Usage{
							InputTokens:           345,
							OutputTokens:          538,
							InputCachedTokens:     234,
							OutputReasoningTokens: 123,
						},
					},
				}),
			},
		},
		{
			name:    "provider sends correct model id, settings, and input",
			modelID: "gpt-4o",
			prompt: []api.Message{
				&api.SystemMessage{
					Content: "You are a helpful assistant.",
				},
				&api.UserMessage{
					Content: api.ContentFromText("Hello"),
				},
			},
			options: api.CallOptions{
				// TODO: We should have a better way for setting temperature.
				// The field should either be a float, or if it's a pointer, we need to include
				// helper functions as part of the API.
				Temperature: pointer.Ptr(0.5),
				TopP:        0.3,
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"temperature": 0.5,
							"top_p": 0.3,
							"input": [
								{
									"role": "system",
									"content": "You are a helpful assistant."
								},
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							]
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				// We don't expect any warnings
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "removes unsupported settings for o1 models",
			modelID: "o1-mini",
			prompt: []api.Message{
				&api.SystemMessage{
					Content: "You are a helpful assistant.",
				},
				&api.UserMessage{
					Content: api.ContentFromText("Hello"),
				},
			},
			options: api.CallOptions{
				Temperature: pointer.Ptr(0.5),
				TopP:        0.3,
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "o1-mini",
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							]
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{
					{
						Type:    "unsupported-setting",
						Setting: "Temperature",
						Details: "Temperature is not supported for reasoning models",
					},
					{
						Type:    "unsupported-setting",
						Setting: "TopP",
						Details: "TopP is not supported for reasoning models",
					},
					{
						Type:    "other",
						Message: "system messages are removed for this model",
					},
				},
			},
		},
		{
			name:    "converts system messages and removes unsupported settings for o3 models",
			modelID: "o3",
			prompt: []api.Message{
				&api.SystemMessage{
					Content: "You are a helpful assistant.",
				},
				&api.UserMessage{
					Content: api.ContentFromText("Hello"),
				},
			},
			options: api.CallOptions{
				Temperature: pointer.Ptr(0.5),
				TopP:        0.3,
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "o3",
							"input": [
								{
									"role": "developer",
									"content": "You are a helpful assistant."
								},
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							]
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{
					{
						Type:    "unsupported-setting",
						Setting: "Temperature",
						Details: "Temperature is not supported for reasoning models",
					},
					{
						Type:    "unsupported-setting",
						Setting: "TopP",
						Details: "TopP is not supported for reasoning models",
					},
				},
			},
		},
		{
			name:    "sends response format json schema",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ResponseFormat: &api.ResponseFormat{
					Type:        "json",
					Name:        "response",
					Description: "A response",
					// TODO: I keep going back and forth on whether Schema should by typed, or simply a map[string]any.
					// FWIW, when the LLM generates the schema, it first defaults to a map[string]any.
					Schema: &jsonschema.Schema{
						Type: "object",
						Properties: map[string]*jsonschema.Schema{
							"value": {Type: "string"},
						},
						Required:             []string{"value"},
						AdditionalProperties: api.FalseSchema(),
					},
				},
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"text": {
								"format": {
									"type": "json_schema",
									"strict": true,
									"name": "response",
									"description": "A response",
									"schema": {
										"type": "object",
										"properties": {
											"value": {
												"type": "string"
											}
										},
										"required": ["value"],
										"additionalProperties": false
									}
								}
							},
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							]
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends response format json object",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ResponseFormat: &api.ResponseFormat{
					Type: "json",
				},
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"text": {
								"format": {
									"type": "json_object"
								}
							},
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							]
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends parallelToolCalls provider option",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"openai": &Metadata{
						ParallelToolCalls: pointer.Ptr(false),
					},
				}),
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							],
							"parallel_tool_calls": false
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends store provider option",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"openai": &Metadata{
						Store: pointer.Ptr(false),
					},
				}),
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							],
							"store": false
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends user provider option",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"openai": &Metadata{
						User: "test-user",
					},
				}),
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							],
							"user": "test-user"
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends previous response id provider option",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"openai": &Metadata{
						PreviousResponseID: "resp_123",
					},
				}),
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							],
							"previous_response_id": "resp_123"
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends reasoningEffort provider option",
			modelID: "o3",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"openai": &Metadata{
						ReasoningEffort: "low",
					},
				}),
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "o3",
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							],
							"reasoning": {
								"effort": "low"
							}
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends metadata provider option with user_123",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"openai": &Metadata{
						User: "user_123",
					},
				}),
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							],
							"user": "user_123"
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends instructions provider option",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"openai": &Metadata{
						Instructions: "You are a friendly assistant.",
					},
				}),
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							],
							"instructions": "You are a friendly assistant."
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends object-tool format",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Tools: []api.ToolDefinition{
					&api.FunctionTool{
						Name:        "response",
						Description: "A response",
						InputSchema: &jsonschema.Schema{
							Type: "object",
							Properties: map[string]*jsonschema.Schema{
								"value": {Type: "string"},
							},
							Required:             []string{"value"},
							AdditionalProperties: api.FalseSchema(),
						},
					},
				},
				ToolChoice: &api.ToolChoice{
					Type:     "tool",
					ToolName: "response",
				},
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"tool_choice": {
								"type": "function",
								"name": "response"
							},
							"tools": [
								{
									"type": "function",
									"strict": true,
									"name": "response",
									"description": "A response",
									"parameters": {
										"type": "object",
										"properties": {
											"value": {
												"type": "string"
											}
										},
										"required": ["value"],
										"additionalProperties": false
									}
								}
							],
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							]
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends object-json json_object format",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ResponseFormat: &api.ResponseFormat{
					Type: "json",
				},
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"text": {
								"format": {
									"type": "json_object"
								}
							},
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							]
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends object-json json_schema format",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ResponseFormat: &api.ResponseFormat{
					Type:        "json",
					Name:        "response",
					Description: "A response",
					Schema: &jsonschema.Schema{
						Type: "object",
						Properties: map[string]*jsonschema.Schema{
							"value": {Type: "string"},
						},
						Required:             []string{"value"},
						AdditionalProperties: api.FalseSchema(),
					},
				},
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"text": {
								"format": {
									"type": "json_schema",
									"strict": true,
									"name": "response",
									"description": "A response",
									"schema": {
										"type": "object",
										"properties": {
											"value": {
												"type": "string"
											}
										},
										"required": ["value"],
										"additionalProperties": false
									}
								}
							},
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							]
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends object-json json_schema format with strictSchemas false",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ResponseFormat: &api.ResponseFormat{
					Type:        "json",
					Name:        "response",
					Description: "A response",
					Schema: &jsonschema.Schema{
						Type: "object",
						Properties: map[string]*jsonschema.Schema{
							"value": {Type: "string"},
						},
						Required:             []string{"value"},
						AdditionalProperties: api.FalseSchema(),
					},
				},
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"openai": &Metadata{
						StrictSchemas: pointer.Ptr(false),
					},
				}),
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"text": {
								"format": {
									"type": "json_schema",
									"strict": false,
									"name": "response",
									"description": "A response",
									"schema": {
										"type": "object",
										"properties": {
											"value": {
												"type": "string"
											}
										},
										"required": ["value"],
										"additionalProperties": false
									}
								}
							},
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							]
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends web_search tool",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Tools: []api.ToolDefinition{
					WebSearchTool(
						WithSearchContextSize("high"),
						WithUserLocation(&WebSearchUserLocation{
							City: "San Francisco",
						}),
					),
				},
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"tools": [
								{
									"type": "web_search_preview",
									"search_context_size": "high",
									"user_location": {
										"city": "San Francisco",
										"type": "approximate"
									}
								}
							],
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							]
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends web_search tool as tool_choice",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ToolChoice: &api.ToolChoice{
					Type:     "tool",
					ToolName: "web_search_preview",
				},
				Tools: []api.ToolDefinition{
					WebSearchTool(
						WithSearchContextSize("high"),
						WithUserLocation(&WebSearchUserLocation{
							City: "San Francisco",
						}),
					),
				},
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"tool_choice": {
								"type": "web_search_preview"
							},
							"tools": [
								{
									"type": "web_search_preview",
									"search_context_size": "high",
									"user_location": {
										"city": "San Francisco",
										"type": "approximate"
									}
								}
							],
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							]
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "warns about unsupported settings",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				StopSequences:    []string{"\n\n"},
				TopK:             1,
				PresencePenalty:  0.1,
				FrequencyPenalty: 0.1,
				Seed:             42,
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "gpt-4o",
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							]
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "answer text"},
				},
				Warnings: []api.CallWarning{
					{Type: "unsupported-setting", Setting: "FrequencyPenalty"},
					{Type: "unsupported-setting", Setting: "PresencePenalty"},
					{Type: "unsupported-setting", Setting: "TopK"},
					{Type: "unsupported-setting", Setting: "Seed"},
					{Type: "unsupported-setting", Setting: "StopSequences"},
				},
			},
		},
	}

	runGenerateTests(t, tests)
}

func TestGenerate_ToolCalls(t *testing.T) {
	standardPrompt := []api.Message{
		&api.UserMessage{
			Content: api.ContentFromText("Hello"),
		},
	}

	// Standard OpenAI response body used across tests
	standardResponseBody := `{
		"id": "resp_67c97c0203188190a025beb4a75242bc",
		"object": "response",
		"created_at": 1741257730,
		"status": "completed",
		"error": null,
		"incomplete_details": null,
		"input": [],
		"instructions": null,
		"max_output_tokens": null,
		"model": "gpt-4o-2024-07-18",
		"output": [
			{
				"type": "function_call",
				"id": "fc_67caf7f4c1ec8190b27edfb5580cfd31",
				"call_id": "call_0NdsJqOS8N3J9l2p0p4WpYU9",
				"name": "weather",
				"arguments": "{\"location\":\"San Francisco\"}",
				"status": "completed"
			},
			{
				"type": "function_call",
				"id": "fc_67caf7f5071c81908209c2909c77af05",
				"call_id": "call_gexo0HtjUfmAIW4gjNOgyrcr",
				"name": "cityAttractions",
				"arguments": "{\"city\":\"San Francisco\"}",
				"status": "completed"
			}
		],
		"parallel_tool_calls": true,
		"previous_response_id": null,
		"reasoning": {
			"effort": null,
			"summary": null
		},
		"store": true,
		"temperature": 1,
		"text": {
			"format": {
				"type": "text"
			}
		},
		"tool_choice": "auto",
		"tools": [
			{
				"type": "function",
				"description": "Get the weather in a location",
				"name": "weather",
				"parameters": {
					"type": "object",
					"properties": {
						"location": {
							"type": "string",
							"description": "The location to get the weather for"
						}
					},
					"required": ["location"],
					"additionalProperties": false
				},
				"strict": true
			},
			{
				"type": "function",
				"description": null,
				"name": "cityAttractions",
				"parameters": {
					"type": "object",
					"properties": {
						"city": {
							"type": "string"
						}
					},
					"required": ["city"],
					"additionalProperties": false
				},
				"strict": true
			}
		],
		"top_p": 1,
		"truncation": "disabled",
		"usage": {
			"input_tokens": 34,
			"output_tokens": 538,
			"output_tokens_details": {
				"reasoning_tokens": 0
			},
			"total_tokens": 572
		},
		"user": null,
		"metadata": {}
	}`

	standardExchange := []httpmock.Exchange{
		{
			Request: httpmock.Request{
				Method: http.MethodPost,
				Path:   "/responses",
			},
			Response: httpmock.Response{
				StatusCode: http.StatusOK,
				Body:       standardResponseBody,
			},
		},
	}

	// Tool-specific test cases
	tests := []struct {
		name         string
		modelID      string          // Optional: override default model
		options      api.CallOptions // Optional: custom call options
		prompt       []api.Message
		exchanges    []httpmock.Exchange
		wantErr      bool
		expectedResp *api.Response // Expected response fields
		skip         bool          // Skip this test if true
	}{
		{
			name:    "should generate tool calls",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Tools: standardTools,
			},
			exchanges: standardExchange,
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "call_0NdsJqOS8N3J9l2p0p4WpYU9",
						ToolName:   "weather",
						Args:       json.RawMessage(`{"location":"San Francisco"}`),
					},
					&api.ToolCallBlock{
						ToolCallID: "call_gexo0HtjUfmAIW4gjNOgyrcr",
						ToolName:   "cityAttractions",
						Args:       json.RawMessage(`{"city":"San Francisco"}`),
					},
				},
			},
		},
		{
			name:    "should have tool-calls finish reason",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Tools: standardTools,
			},
			exchanges: standardExchange,
			expectedResp: &api.Response{
				FinishReason: api.FinishReasonToolCalls,
			},
		},
	}

	runGenerateTests(t, tests)
}

func TestGenerate_WebSearch(t *testing.T) {
	standardPrompt := []api.Message{
		&api.UserMessage{
			Content: api.ContentFromText("What happened in San Francisco last week?"),
		},
	}

	// Get the expected text - defined early to use in the response
	expectedText := "Last week in San Francisco, several notable events and developments took place:\n\n" +
		"**Bruce Lee Statue in Chinatown**\n\n" +
		"The Chinese Historical Society of America Museum announced plans to install a Bruce Lee statue in Chinatown. This initiative, supported by the Rose Pak Community Fund, the Bruce Lee Foundation, and Stand With Asians, aims to honor Lee's contributions to film and martial arts. Artist Arnie Kim has been commissioned for the project, with a fundraising goal of $150,000. ([axios.com](https://www.axios.com/local/san-francisco/2025/03/07/bruce-lee-statue-sf-chinatown?utm_source=chatgpt.com))\n\n" +
		"**Office Leasing Revival**\n\n" +
		"The Bay Area experienced a resurgence in office leasing, securing 11 of the largest U.S. office leases in 2024. This trend, driven by the tech industry's growth and advancements in generative AI, suggests a potential boost to downtown recovery through increased foot traffic. ([axios.com](https://www.axios.com/local/san-francisco/2025/03/03/bay-area-office-leasing-activity?utm_source=chatgpt.com))\n\n" +
		"**Spring Blooms in the Bay Area**\n\n" +
		"With the arrival of spring, several locations in the Bay Area are showcasing vibrant blooms. Notable spots include the Conservatory of Flowers, Japanese Tea Garden, Queen Wilhelmina Tulip Garden, and the San Francisco Botanical Garden, each offering unique floral displays. ([axios.com](https://www.axios.com/local/san-francisco/2025/03/03/where-to-see-spring-blooms-bay-area?utm_source=chatgpt.com))\n\n" +
		"**Oceanfront Great Highway Park**\n\n" +
		"San Francisco's long-awaited Oceanfront Great Highway park is set to open on April 12. This 43-acre, car-free park will span a two-mile stretch of the Great Highway from Lincoln Way to Sloat Boulevard, marking the largest pedestrianization project in California's history. The park follows voter approval of Proposition K, which permanently bans cars on part of the highway. ([axios.com](https://www.axios.com/local/san-francisco/2025/03/03/great-highway-park-opening-april-recall-campaign?utm_source=chatgpt.com))\n\n" +
		"**Warmer Spring Seasons**\n\n" +
		"An analysis by Climate Central revealed that San Francisco, along with most U.S. cities, is experiencing increasingly warmer spring seasons. Over a 55-year period from 1970 to 2024, the national average temperature during March through May rose by 2.4°F. This warming trend poses various risks, including early snowmelt and increased wildfire threats. ([axios.com](https://www.axios.com/local/san-francisco/2025/03/03/climate-weather-spring-temperatures-warmer-sf?utm_source=chatgpt.com))\n\n\n" +
		"# Key San Francisco Developments Last Week:\n" +
		"- [Bruce Lee statue to be installed in SF Chinatown](https://www.axios.com/local/san-francisco/2025/03/07/bruce-lee-statue-sf-chinatown?utm_source=chatgpt.com)\n" +
		"- [The Bay Area is set to make an office leasing comeback](https://www.axios.com/local/san-francisco/2025/03/03/bay-area-office-leasing-activity?utm_source=chatgpt.com)\n" +
		"- [Oceanfront Great Highway park set to open in April](https://www.axios.com/local/san-francisco/2025/03/03/great-highway-park-opening-april-recall-campaign?utm_source=chatgpt.com)"

	// Standard OpenAI response body using map[string]any instead of a string
	standardResponseBody := map[string]any{
		"id":                 "resp_67cf2b2f6bd081909be2c8054ddef0eb",
		"object":             "response",
		"created_at":         1741630255,
		"status":             "completed",
		"error":              nil,
		"incomplete_details": nil,
		"instructions":       nil,
		"max_output_tokens":  nil,
		"model":              "gpt-4o-2024-07-18",
		"output": []map[string]any{
			{
				"type":   "web_search_call",
				"id":     "ws_67cf2b3051e88190b006770db6fdb13d",
				"status": "completed",
			},
			{
				"type":   "message",
				"id":     "msg_67cf2b35467481908f24412e4fd40d66",
				"status": "completed",
				"role":   "assistant",
				"content": []map[string]any{
					{
						"type": "output_text",
						"text": expectedText,
						"annotations": []map[string]any{
							{
								"type":        "url_citation",
								"start_index": 486,
								"end_index":   606,
								"url":         "https://www.axios.com/local/san-francisco/2025/03/07/bruce-lee-statue-sf-chinatown?utm_source=chatgpt.com",
								"title":       "Bruce Lee statue to be installed in SF Chinatown",
							},
							{
								"type":        "url_citation",
								"start_index": 912,
								"end_index":   1035,
								"url":         "https://www.axios.com/local/san-francisco/2025/03/03/bay-area-office-leasing-activity?utm_source=chatgpt.com",
								"title":       "The Bay Area is set to make an office leasing comeback",
							},
							{
								"type":        "url_citation",
								"start_index": 1346,
								"end_index":   1472,
								"url":         "https://www.axios.com/local/san-francisco/2025/03/03/where-to-see-spring-blooms-bay-area?utm_source=chatgpt.com",
								"title":       "Where to see spring blooms in the Bay Area",
							},
							{
								"type":        "url_citation",
								"start_index": 1884,
								"end_index":   2023,
								"url":         "https://www.axios.com/local/san-francisco/2025/03/03/great-highway-park-opening-april-recall-campaign?utm_source=chatgpt.com",
								"title":       "Oceanfront Great Highway park set to open in April",
							},
							{
								"type":        "url_citation",
								"start_index": 2404,
								"end_index":   2540,
								"url":         "https://www.axios.com/local/san-francisco/2025/03/03/climate-weather-spring-temperatures-warmer-sf?utm_source=chatgpt.com",
								"title":       "San Francisco's spring seasons are getting warmer",
							},
						},
					},
				},
			},
		},
		"parallel_tool_calls":  true,
		"previous_response_id": nil,
		"reasoning": map[string]any{
			"effort":  nil,
			"summary": nil,
		},
		"store":       true,
		"temperature": 0,
		"text": map[string]any{
			"format": map[string]any{
				"type": "text",
			},
		},
		"tool_choice": "auto",
		"tools": []map[string]any{
			{
				"type":                "web_search_preview",
				"search_context_size": "medium",
				"user_location": map[string]any{
					"type":     "approximate",
					"city":     nil,
					"country":  "US",
					"region":   nil,
					"timezone": nil,
				},
			},
		},
		"top_p":      1,
		"truncation": "disabled",
		"usage": map[string]any{
			"input_tokens": 327,
			"input_tokens_details": map[string]any{
				"cached_tokens": 0,
			},
			"output_tokens": 770,
			"output_tokens_details": map[string]any{
				"reasoning_tokens": 0,
			},
			"total_tokens": 1097,
		},
		"user":     nil,
		"metadata": map[string]any{},
	}

	standardExchange := []httpmock.Exchange{
		{
			Request: httpmock.Request{
				Method: http.MethodPost,
				Path:   "/responses",
			},
			Response: httpmock.Response{
				StatusCode: http.StatusOK,
				Body:       standardResponseBody,
			},
		},
	}

	// No need to extract text from the response body now, just use the expectedText variable
	// that we defined earlier

	// Tool-specific test cases
	tests := []struct {
		name         string
		modelID      string          // Optional: override default model
		options      api.CallOptions // Optional: custom call options
		prompt       []api.Message
		exchanges    []httpmock.Exchange
		wantErr      bool
		expectedResp *api.Response // Expected response fields
		skip         bool          // Skip this test if true
	}{
		{
			name:      "should generate text with web search",
			modelID:   "gpt-4o",
			prompt:    standardPrompt,
			exchanges: standardExchange,
			options: api.CallOptions{
				Tools: []api.ToolDefinition{
					WebSearchTool(
						WithSearchContextSize("medium"),
						WithUserLocation(&WebSearchUserLocation{
							Country: "US",
						}),
					),
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "ws_67cf2b3051e88190b006770db6fdb13d",
						ToolName:   "openai.web_search_preview",
					},
					&api.TextBlock{Text: expectedText},
				},
			},
		},
		{
			name:      "should return sources from web search",
			modelID:   "gpt-4o",
			prompt:    standardPrompt,
			exchanges: standardExchange,
			options: api.CallOptions{
				Tools: []api.ToolDefinition{
					WebSearchTool(
						WithSearchContextSize("medium"),
						WithUserLocation(&WebSearchUserLocation{
							Country: "US",
						}),
					),
				},
			},
			expectedResp: &api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "ws_67cf2b3051e88190b006770db6fdb13d",
						ToolName:   "openai.web_search_preview",
					},
					&api.TextBlock{Text: expectedText},
					&api.SourceBlock{
						ID:    "source-0",
						URL:   "https://www.axios.com/local/san-francisco/2025/03/07/bruce-lee-statue-sf-chinatown?utm_source=chatgpt.com",
						Title: "Bruce Lee statue to be installed in SF Chinatown",
					},
					&api.SourceBlock{
						ID:    "source-1",
						URL:   "https://www.axios.com/local/san-francisco/2025/03/03/bay-area-office-leasing-activity?utm_source=chatgpt.com",
						Title: "The Bay Area is set to make an office leasing comeback",
					},
					&api.SourceBlock{
						ID:    "source-2",
						URL:   "https://www.axios.com/local/san-francisco/2025/03/03/where-to-see-spring-blooms-bay-area?utm_source=chatgpt.com",
						Title: "Where to see spring blooms in the Bay Area",
					},
					&api.SourceBlock{
						ID:    "source-3",
						URL:   "https://www.axios.com/local/san-francisco/2025/03/03/great-highway-park-opening-april-recall-campaign?utm_source=chatgpt.com",
						Title: "Oceanfront Great Highway park set to open in April",
					},
					&api.SourceBlock{
						ID:    "source-4",
						URL:   "https://www.axios.com/local/san-francisco/2025/03/03/climate-weather-spring-temperatures-warmer-sf?utm_source=chatgpt.com",
						Title: "San Francisco's spring seasons are getting warmer",
					},
				},
			},
		},
	}

	runGenerateTests(t, tests)
}

// eventsToString converts a slice of SSE events to a string format expected by the mock server
func eventsToString(events []sse.Event) string {
	var buf bytes.Buffer
	enc := sse.NewEncoder(&buf)
	for _, event := range events {
		if err := enc.EncodeEvent(&event); err != nil {
			panic(fmt.Sprintf("failed to encode event: %v", err))
		}
	}
	// Add the [DONE] marker
	buf.WriteString("data: [DONE]\n\n")
	return buf.String()
}

func TestStream(t *testing.T) {
	standardPrompt := []api.Message{
		&api.UserMessage{
			Content: api.ContentFromText("Hello"),
		},
	}

	tests := []struct {
		name           string
		modelID        string
		options        api.CallOptions
		prompt         []api.Message
		exchanges      []httpmock.Exchange
		wantErr        bool
		expectedEvents []api.StreamEvent
		skip           bool
	}{
		{
			name:    "should stream text deltas",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Headers: map[string]string{
							"Content-Type": "text/event-stream",
						},
						Body: eventsToString([]sse.Event{
							{
								Data: map[string]any{
									"type": "response.created",
									"response": map[string]any{
										"id":         "resp_67c9a81b6a048190a9ee441c5755a4e8",
										"object":     "response",
										"created_at": 1741269019,
										"status":     "in_progress",
										"model":      "gpt-4o-2024-07-18",
									},
								},
							},
							{
								Data: map[string]any{
									"type": "response.in_progress",
									"response": map[string]any{
										"id":         "resp_67c9a81b6a048190a9ee441c5755a4e8",
										"object":     "response",
										"created_at": 1741269019,
										"status":     "in_progress",
										"model":      "gpt-4o-2024-07-18",
									},
								},
							},
							{
								Data: map[string]any{
									"type":         "response.output_item.added",
									"output_index": 0,
									"item": map[string]any{
										"id":      "msg_67c9a81dea8c8190b79651a2b3adf91e",
										"type":    "message",
										"status":  "in_progress",
										"role":    "assistant",
										"content": []any{},
									},
								},
							},
							{
								Data: map[string]any{
									"type":          "response.content_part.added",
									"item_id":       "msg_67c9a81dea8c8190b79651a2b3adf91e",
									"output_index":  0,
									"content_index": 0,
									"part": map[string]any{
										"type":        "output_text",
										"text":        "",
										"annotations": []any{},
									},
								},
							},
							{
								Data: map[string]any{
									"type":          "response.output_text.delta",
									"item_id":       "msg_67c9a81dea8c8190b79651a2b3adf91e",
									"output_index":  0,
									"content_index": 0,
									"delta":         "Hello,",
								},
							},
							{
								Data: map[string]any{
									"type":          "response.output_text.delta",
									"item_id":       "msg_67c9a81dea8c8190b79651a2b3adf91e",
									"output_index":  0,
									"content_index": 0,
									"delta":         " World!",
								},
							},
							{
								Data: map[string]any{
									"type":          "response.output_text.done",
									"item_id":       "msg_67c9a8787f4c8190b49c858d4c1cf20c",
									"output_index":  0,
									"content_index": 0,
									"text":          "Hello, World!",
								},
							},
							{
								Data: map[string]any{
									"type":          "response.content_part.done",
									"item_id":       "msg_67c9a8787f4c8190b49c858d4c1cf20c",
									"output_index":  0,
									"content_index": 0,
									"part": map[string]any{
										"type":        "output_text",
										"text":        "Hello, World!",
										"annotations": []any{},
									},
								},
							},
							{
								Data: map[string]any{
									"type":         "response.output_item.done",
									"output_index": 0,
									"item": map[string]any{
										"id":     "msg_67c9a8787f4c8190b49c858d4c1cf20c",
										"type":   "message",
										"status": "completed",
										"role":   "assistant",
										"content": []any{
											map[string]any{
												"type":        "output_text",
												"text":        "Hello, World!",
												"annotations": []any{},
											},
										},
									},
								},
							},
							{
								Data: map[string]any{
									"type": "response.completed",
									"response": map[string]any{
										"id":         "resp_67c9a878139c8190aa2e3105411b408b",
										"object":     "response",
										"created_at": 1741269112,
										"status":     "completed",
										"model":      "gpt-4o-2024-07-18",
										"output": []any{
											map[string]any{
												"id":     "msg_67c9a8787f4c8190b49c858d4c1cf20c",
												"type":   "message",
												"status": "completed",
												"role":   "assistant",
												"content": []any{
													map[string]any{
														"type":        "output_text",
														"text":        "Hello, World!",
														"annotations": []any{},
													},
												},
											},
										},
										"usage": map[string]any{
											"input_tokens": 543,
											"input_tokens_details": map[string]any{
												"cached_tokens": 234,
											},
											"output_tokens": 478,
											"output_tokens_details": map[string]any{
												"reasoning_tokens": 123,
											},
											"total_tokens": 512,
										},
									},
								},
							},
						}),
					},
				},
			},
			expectedEvents: []api.StreamEvent{
				&api.ResponseMetadataEvent{
					ID:        "resp_67c9a81b6a048190a9ee441c5755a4e8",
					ModelID:   "gpt-4o-2024-07-18",
					Timestamp: time.Date(2025, 3, 6, 13, 50, 19, 0, time.UTC),
				},
				&api.TextDeltaEvent{
					TextDelta: "Hello,",
				},
				&api.TextDeltaEvent{
					TextDelta: " World!",
				},
				&api.FinishEvent{
					FinishReason: api.FinishReasonStop,
					Usage: api.Usage{
						InputTokens:       543,
						OutputTokens:      478,
						TotalTokens:       512,
						ReasoningTokens:   123,
						CachedInputTokens: 234,
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{
						"openai": &Metadata{
							ResponseID: "resp_67c9a81b6a048190a9ee441c5755a4e8",
							Usage: codec.Usage{
								InputTokens:           543,
								OutputTokens:          478,
								InputCachedTokens:     234,
								OutputReasoningTokens: 123,
							},
						},
					}),
				},
			},
		},
		{
			name:    "should send finish reason for incomplete response",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Headers: map[string]string{
							"Content-Type": "text/event-stream",
						},
						Body: eventsToString([]sse.Event{
							{
								Data: map[string]any{
									"type": "response.created",
									"response": map[string]any{
										"id":         "resp_67c9a81b6a048190a9ee441c5755a4e8",
										"object":     "response",
										"created_at": 1741269019,
										"status":     "in_progress",
										"model":      "gpt-4o-2024-07-18",
									},
								},
							},
							{
								Data: map[string]any{
									"type": "response.in_progress",
									"response": map[string]any{
										"id":         "resp_67c9a81b6a048190a9ee441c5755a4e8",
										"object":     "response",
										"created_at": 1741269019,
										"status":     "in_progress",
										"model":      "gpt-4o-2024-07-18",
									},
								},
							},
							{
								Data: map[string]any{
									"type":         "response.output_item.added",
									"output_index": 0,
									"item": map[string]any{
										"id":      "msg_67c9a81dea8c8190b79651a2b3adf91e",
										"type":    "message",
										"status":  "in_progress",
										"role":    "assistant",
										"content": []any{},
									},
								},
							},
							{
								Data: map[string]any{
									"type":          "response.content_part.added",
									"item_id":       "msg_67c9a81dea8c8190b79651a2b3adf91e",
									"output_index":  0,
									"content_index": 0,
									"part": map[string]any{
										"type":        "output_text",
										"text":        "",
										"annotations": []any{},
									},
								},
							},
							{
								Data: map[string]any{
									"type":          "response.output_text.delta",
									"item_id":       "msg_67c9a81dea8c8190b79651a2b3adf91e",
									"output_index":  0,
									"content_index": 0,
									"delta":         "Hello,",
								},
							},
							{
								Data: map[string]any{
									"type":          "response.output_text.done",
									"item_id":       "msg_67c9a8787f4c8190b49c858d4c1cf20c",
									"output_index":  0,
									"content_index": 0,
									"text":          "Hello,!",
								},
							},
							{
								Data: map[string]any{
									"type":          "response.content_part.done",
									"item_id":       "msg_67c9a8787f4c8190b49c858d4c1cf20c",
									"output_index":  0,
									"content_index": 0,
									"part": map[string]any{
										"type":        "output_text",
										"text":        "Hello,",
										"annotations": []any{},
									},
								},
							},
							{
								Data: map[string]any{
									"type":         "response.output_item.done",
									"output_index": 0,
									"item": map[string]any{
										"id":     "msg_67c9a8787f4c8190b49c858d4c1cf20c",
										"type":   "message",
										"status": "incomplete",
										"role":   "assistant",
										"content": []any{
											map[string]any{
												"type":        "output_text",
												"text":        "Hello,",
												"annotations": []any{},
											},
										},
									},
								},
							},
							{
								Data: map[string]any{
									"type": "response.incomplete",
									"response": map[string]any{
										"id":         "resp_67cadb40a0708190ac2763c0b6960f6f",
										"object":     "response",
										"created_at": 1741347648,
										"status":     "incomplete",
										"incomplete_details": map[string]any{
											"reason": "max_output_tokens",
										},
										"model": "gpt-4o-2024-07-18",
										"output": []any{
											map[string]any{
												"type":   "message",
												"id":     "msg_67cadb410ccc81909fe1d8f427b9cf02",
												"status": "incomplete",
												"role":   "assistant",
												"content": []any{
													map[string]any{
														"type":        "output_text",
														"text":        "Hello,",
														"annotations": []any{},
													},
												},
											},
										},
										"max_output_tokens": 100,
									},
								},
							},
						}),
					},
				},
			},
			expectedEvents: []api.StreamEvent{
				&api.ResponseMetadataEvent{
					ID:        "resp_67c9a81b6a048190a9ee441c5755a4e8",
					ModelID:   "gpt-4o-2024-07-18",
					Timestamp: time.Date(2025, 3, 6, 13, 50, 19, 0, time.UTC),
				},
				&api.TextDeltaEvent{
					TextDelta: "Hello,",
				},
				&api.FinishEvent{
					FinishReason: api.FinishReasonLength,
					Usage: api.Usage{
						InputTokens:  0,
						OutputTokens: 0,
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{
						"openai": &Metadata{
							ResponseID: "resp_67c9a81b6a048190a9ee441c5755a4e8",
							Usage:      codec.Usage{},
						},
					}),
				},
			},
		},
		{
			name:    "should stream tool calls",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Tools: standardTools,
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Headers: map[string]string{
							"Content-Type": "text/event-stream",
						},
						Body: eventsToString([]sse.Event{
							{
								Data: map[string]any{
									"type": "response.created",
									"response": map[string]any{
										"id":         "resp_67cb13a755c08190acbe3839a49632fc",
										"object":     "response",
										"created_at": 1741362087,
										"status":     "in_progress",
										"model":      "gpt-4o-2024-07-18",
										"tools": []any{
											map[string]any{
												"type":        "function",
												"description": "Get the current location.",
												"name":        "currentLocation",
												"parameters": map[string]any{
													"type":                 "object",
													"properties":           map[string]any{},
													"additionalProperties": false,
												},
												"strict": true,
											},
											map[string]any{
												"type":        "function",
												"description": "Get the weather in a location",
												"name":        "weather",
												"parameters": map[string]any{
													"type": "object",
													"properties": map[string]any{
														"location": map[string]any{
															"type":        "string",
															"description": "The location to get the weather for",
														},
													},
													"required":             []string{"location"},
													"additionalProperties": false,
												},
												"strict": true,
											},
										},
									},
								},
							},
							{
								Data: map[string]any{
									"type": "response.in_progress",
									"response": map[string]any{
										"id":         "resp_67cb13a755c08190acbe3839a49632fc",
										"object":     "response",
										"created_at": 1741362087,
										"status":     "in_progress",
										"model":      "gpt-4o-2024-07-18",
										"tools": []any{
											map[string]any{
												"type":        "function",
												"description": "Get the current location.",
												"name":        "currentLocation",
												"parameters": map[string]any{
													"type":                 "object",
													"properties":           map[string]any{},
													"additionalProperties": false,
												},
												"strict": true,
											},
											map[string]any{
												"type":        "function",
												"description": "Get the weather in a location",
												"name":        "weather",
												"parameters": map[string]any{
													"type": "object",
													"properties": map[string]any{
														"location": map[string]any{
															"type":        "string",
															"description": "The location to get the weather for",
														},
													},
													"required":             []string{"location"},
													"additionalProperties": false,
												},
												"strict": true,
											},
										},
									},
								},
							},
							{
								Data: map[string]any{
									"type":         "response.output_item.added",
									"output_index": 0,
									"item": map[string]any{
										"type":      "function_call",
										"id":        "fc_67cb13a838088190be08eb3927c87501",
										"call_id":   "call_6KxSghkb4MVnunFH2TxPErLP",
										"name":      "currentLocation",
										"arguments": "",
										"status":    "completed",
									},
								},
							},
							{
								Data: map[string]any{
									"type":         "response.function_call_arguments.delta",
									"item_id":      "fc_67cb13a838088190be08eb3927c87501",
									"output_index": 0,
									"delta":        "{}",
								},
							},
							{
								Data: map[string]any{
									"type":         "response.function_call_arguments.done",
									"item_id":      "fc_67cb13a838088190be08eb3927c87501",
									"output_index": 0,
									"arguments":    "{}",
								},
							},
							{
								Data: map[string]any{
									"type":         "response.output_item.done",
									"output_index": 0,
									"item": map[string]any{
										"type":      "function_call",
										"id":        "fc_67cb13a838088190be08eb3927c87501",
										"call_id":   "call_pgjcAI4ZegMkP6bsAV7sfrJA",
										"name":      "currentLocation",
										"arguments": "{}",
										"status":    "completed",
									},
								},
							},
							{
								Data: map[string]any{
									"type":         "response.output_item.added",
									"output_index": 1,
									"item": map[string]any{
										"type":      "function_call",
										"id":        "fc_67cb13a858f081908a600343fa040f47",
										"call_id":   "call_Dg6WUmFHNeR5JxX1s53s1G4b",
										"name":      "weather",
										"arguments": "",
										"status":    "in_progress",
									},
								},
							},
							{
								Data: map[string]any{
									"type":         "response.function_call_arguments.delta",
									"item_id":      "fc_67cb13a858f081908a600343fa040f47",
									"output_index": 1,
									"delta":        "{",
								},
							},
							{
								Data: map[string]any{
									"type":         "response.function_call_arguments.delta",
									"item_id":      "fc_67cb13a858f081908a600343fa040f47",
									"output_index": 1,
									"delta":        "\"location\"",
								},
							},
							{
								Data: map[string]any{
									"type":         "response.function_call_arguments.delta",
									"item_id":      "fc_67cb13a858f081908a600343fa040f47",
									"output_index": 1,
									"delta":        "\":\"",
								},
							},
							{
								Data: map[string]any{
									"type":         "response.function_call_arguments.delta",
									"item_id":      "fc_67cb13a858f081908a600343fa040f47",
									"output_index": 1,
									"delta":        "\"Rome\"",
								},
							},
							{
								Data: map[string]any{
									"type":         "response.function_call_arguments.delta",
									"item_id":      "fc_67cb13a858f081908a600343fa040f47",
									"output_index": 1,
									"delta":        "\"}\"",
								},
							},
							{
								Data: map[string]any{
									"type":         "response.function_call_arguments.done",
									"item_id":      "fc_67cb13a858f081908a600343fa040f47",
									"output_index": 1,
									"arguments":    "{\"location\":\"Rome\"}",
								},
							},
							{
								Data: map[string]any{
									"type":         "response.output_item.done",
									"output_index": 1,
									"item": map[string]any{
										"type":      "function_call",
										"id":        "fc_67cb13a858f081908a600343fa040f47",
										"call_id":   "call_X2PAkDJInno9VVnNkDrfhboW",
										"name":      "weather",
										"arguments": "{\"location\":\"Rome\"}",
										"status":    "completed",
									},
								},
							},
							{
								Data: map[string]any{
									"type": "response.completed",
									"response": map[string]any{
										"id":         "resp_67cb13a755c08190acbe3839a49632fc",
										"object":     "response",
										"created_at": 1741362087,
										"status":     "completed",
										"model":      "gpt-4o-2024-07-18",
										"output": []any{
											map[string]any{
												"type":      "function_call",
												"id":        "fc_67cb13a838088190be08eb3927c87501",
												"call_id":   "call_KsVqaVAf3alAtCCkQe4itE7W",
												"name":      "currentLocation",
												"arguments": "{}",
												"status":    "completed",
											},
											map[string]any{
												"type":      "function_call",
												"id":        "fc_67cb13a858f081908a600343fa040f47",
												"call_id":   "call_X2PAkDJInno9VVnNkDrfhboW",
												"name":      "weather",
												"arguments": "{\"location\":\"Rome\"}",
												"status":    "completed",
											},
										},
									},
								},
							},
						}),
					},
				},
			},
			expectedEvents: []api.StreamEvent{
				&api.ResponseMetadataEvent{
					ID:        "resp_67cb13a755c08190acbe3839a49632fc",
					ModelID:   "gpt-4o-2024-07-18",
					Timestamp: time.Date(2025, 3, 7, 15, 41, 27, 0, time.UTC),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_6KxSghkb4MVnunFH2TxPErLP",
					ToolName:   "currentLocation",
					ArgsDelta:  []byte(""),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_6KxSghkb4MVnunFH2TxPErLP",
					ToolName:   "currentLocation",
					ArgsDelta:  []byte("{}"),
				},
				&api.ToolCallEvent{
					ToolCallID: "call_pgjcAI4ZegMkP6bsAV7sfrJA",
					ToolName:   "currentLocation",
					Args:       json.RawMessage(`{}`),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_Dg6WUmFHNeR5JxX1s53s1G4b",
					ToolName:   "weather",
					ArgsDelta:  []byte(""),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_Dg6WUmFHNeR5JxX1s53s1G4b",
					ToolName:   "weather",
					ArgsDelta:  []byte("{"),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_Dg6WUmFHNeR5JxX1s53s1G4b",
					ToolName:   "weather",
					ArgsDelta:  []byte("\"location\""),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_Dg6WUmFHNeR5JxX1s53s1G4b",
					ToolName:   "weather",
					ArgsDelta:  []byte("\":\""),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_Dg6WUmFHNeR5JxX1s53s1G4b",
					ToolName:   "weather",
					ArgsDelta:  []byte("\"Rome\""),
				},
				&api.ToolCallDeltaEvent{
					ToolCallID: "call_Dg6WUmFHNeR5JxX1s53s1G4b",
					ToolName:   "weather",
					ArgsDelta:  []byte("\"}\""),
				},
				&api.ToolCallEvent{
					ToolCallID: "call_X2PAkDJInno9VVnNkDrfhboW",
					ToolName:   "weather",
					Args:       json.RawMessage(`{"location":"Rome"}`),
				},
				&api.FinishEvent{
					FinishReason: api.FinishReasonToolCalls,
					Usage: api.Usage{
						InputTokens:  0,
						OutputTokens: 0,
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{
						"openai": &Metadata{
							ResponseID: "resp_67cb13a755c08190acbe3839a49632fc",
							Usage:      codec.Usage{},
						},
					}),
				},
			},
		},
		{
			name:    "should stream sources",
			modelID: "gpt-4o-mini",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Tools: []api.ToolDefinition{
					WebSearchTool(
						WithSearchContextSize("medium"),
						WithUserLocation(&WebSearchUserLocation{
							Country: "US",
						}),
					),
				},
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Headers: map[string]string{
							"Content-Type": "text/event-stream",
						},
						Body: eventsToString([]sse.Event{
							{
								Data: map[string]any{
									"type": "response.created",
									"response": map[string]any{
										"id":         "resp_67cf3390786881908b27489d7e8cfb6b",
										"object":     "response",
										"created_at": 1741632400,
										"status":     "in_progress",
										"model":      "gpt-4o-mini-2024-07-18",
										"tools": []any{
											map[string]any{
												"type":                "web_search_preview",
												"search_context_size": "medium",
												"user_location": map[string]any{
													"type":     "approximate",
													"city":     nil,
													"country":  "US",
													"region":   nil,
													"timezone": nil,
												},
											},
										},
									},
								},
							},
							{
								Data: map[string]any{
									"type": "response.in_progress",
									"response": map[string]any{
										"id":         "resp_67cf3390786881908b27489d7e8cfb6b",
										"object":     "response",
										"created_at": 1741632400,
										"status":     "in_progress",
										"model":      "gpt-4o-mini-2024-07-18",
										"tools": []any{
											map[string]any{
												"type":                "web_search_preview",
												"search_context_size": "medium",
												"user_location": map[string]any{
													"type":     "approximate",
													"city":     nil,
													"country":  "US",
													"region":   nil,
													"timezone": nil,
												},
											},
										},
									},
								},
							},
							{
								Data: map[string]any{
									"type":         "response.output_item.added",
									"output_index": 0,
									"item": map[string]any{
										"type":   "web_search_call",
										"id":     "ws_67cf3390e9608190869b5d45698a7067",
										"status": "in_progress",
									},
								},
							},
							{
								Data: map[string]any{
									"type":         "response.web_search_call.in_progress",
									"output_index": 0,
									"item_id":      "ws_67cf3390e9608190869b5d45698a7067",
								},
							},
							{
								Data: map[string]any{
									"type":         "response.web_search_call.searching",
									"output_index": 0,
									"item_id":      "ws_67cf3390e9608190869b5d45698a7067",
								},
							},
							{
								Data: map[string]any{
									"type":         "response.web_search_call.completed",
									"output_index": 0,
									"item_id":      "ws_67cf3390e9608190869b5d45698a7067",
								},
							},
							{
								Data: map[string]any{
									"type":         "response.output_item.done",
									"output_index": 0,
									"item": map[string]any{
										"type":   "web_search_call",
										"id":     "ws_67cf3390e9608190869b5d45698a7067",
										"status": "completed",
									},
								},
							},
							{
								Data: map[string]any{
									"type":         "response.output_item.added",
									"output_index": 1,
									"item": map[string]any{
										"type":    "message",
										"id":      "msg_67cf33924ea88190b8c12bf68c1f6416",
										"status":  "in_progress",
										"role":    "assistant",
										"content": []any{},
									},
								},
							},
							{
								Data: map[string]any{
									"type":          "response.content_part.added",
									"item_id":       "msg_67cf33924ea88190b8c12bf68c1f6416",
									"output_index":  1,
									"content_index": 0,
									"part": map[string]any{
										"type":        "output_text",
										"text":        "",
										"annotations": []any{},
									},
								},
							},
							{
								Data: map[string]any{
									"type":          "response.output_text.delta",
									"item_id":       "msg_67cf33924ea88190b8c12bf68c1f6416",
									"output_index":  1,
									"content_index": 0,
									"delta":         "Last week",
								},
							},
							{
								Data: map[string]any{
									"type":          "response.output_text.delta",
									"item_id":       "msg_67cf33924ea88190b8c12bf68c1f6416",
									"output_index":  1,
									"content_index": 0,
									"delta":         " in San Francisco",
								},
							},
							{
								Data: map[string]any{
									"type":             "response.output_text.annotation.added",
									"item_id":          "msg_67cf33924ea88190b8c12bf68c1f6416",
									"output_index":     1,
									"content_index":    0,
									"annotation_index": 0,
									"annotation": map[string]any{
										"type":        "url_citation",
										"start_index": 383,
										"end_index":   493,
										"url":         "https://www.sftourismtips.com/san-francisco-events-in-march.html?utm_source=chatgpt.com",
										"title":       "San Francisco Events in March 2025: Festivals, Theater & Easter",
									},
								},
							},
							{
								Data: map[string]any{
									"type":          "response.output_text.delta",
									"item_id":       "msg_67cf33924ea88190b8c12bf68c1f6416",
									"output_index":  1,
									"content_index": 0,
									"delta":         " a themed party",
								},
							},
							{
								Data: map[string]any{
									"type":          "response.output_text.delta",
									"item_id":       "msg_67cf33924ea88190b8c12bf68c1f6416",
									"output_index":  1,
									"content_index": 0,
									"delta":         "([axios.com](https://www.axios.com/local/san-francisco/2025/03/06/sf-events-march-what-to-do-giants-fanfest?utm_source=chatgpt.com))",
								},
							},
							{
								Data: map[string]any{
									"type":             "response.output_text.annotation.added",
									"item_id":          "msg_67cf33924ea88190b8c12bf68c1f6416",
									"output_index":     1,
									"content_index":    0,
									"annotation_index": 1,
									"annotation": map[string]any{
										"type":        "url_citation",
										"start_index": 630,
										"end_index":   762,
										"url":         "https://www.axios.com/local/san-francisco/2025/03/06/sf-events-march-what-to-do-giants-fanfest?utm_source=chatgpt.com",
										"title":       "SF weekend events: Giants FanFest, crab crawl and more",
									},
								},
							},
							{
								Data: map[string]any{
									"type":          "response.output_text.delta",
									"item_id":       "msg_67cf33924ea88190b8c12bf68c1f6416",
									"output_index":  1,
									"content_index": 0,
									"delta":         ".",
								},
							},
							{
								Data: map[string]any{
									"type":          "response.output_text.done",
									"item_id":       "msg_67cf33924ea88190b8c12bf68c1f6416",
									"output_index":  1,
									"content_index": 0,
									"text":          "Last week in San Francisco a themed party...",
								},
							},
							{
								Data: map[string]any{
									"type":          "response.content_part.done",
									"item_id":       "msg_67cf33924ea88190b8c12bf68c1f6416",
									"output_index":  1,
									"content_index": 0,
									"part": map[string]any{
										"type": "output_text",
										"text": "Last week in San Francisco a themed party...",
										"annotations": []any{
											map[string]any{
												"type":        "url_citation",
												"start_index": 383,
												"end_index":   493,
												"url":         "https://www.sftourismtips.com/san-francisco-events-in-march.html?utm_source=chatgpt.com",
												"title":       "San Francisco Events in March 2025: Festivals, Theater & Easter",
											},
											map[string]any{
												"type":        "url_citation",
												"start_index": 630,
												"end_index":   762,
												"url":         "https://www.axios.com/local/san-francisco/2025/03/06/sf-events-march-what-to-do-giants-fanfest?utm_source=chatgpt.com",
												"title":       "SF weekend events: Giants FanFest, crab crawl and more",
											},
										},
									},
								},
							},
							{
								Data: map[string]any{
									"type":         "response.output_item.done",
									"output_index": 1,
									"item": map[string]any{
										"type":   "message",
										"id":     "msg_67cf33924ea88190b8c12bf68c1f6416",
										"status": "completed",
										"role":   "assistant",
										"content": []any{
											map[string]any{
												"type": "output_text",
												"text": "Last week in San Francisco a themed party...",
												"annotations": []any{
													map[string]any{
														"type":        "url_citation",
														"start_index": 383,
														"end_index":   493,
														"url":         "https://www.sftourismtips.com/san-francisco-events-in-march.html?utm_source=chatgpt.com",
														"title":       "San Francisco Events in March 2025: Festivals, Theater & Easter",
													},
													map[string]any{
														"type":        "url_citation",
														"start_index": 630,
														"end_index":   762,
														"url":         "https://www.axios.com/local/san-francisco/2025/03/06/sf-events-march-what-to-do-giants-fanfest?utm_source=chatgpt.com",
														"title":       "SF weekend events: Giants FanFest, crab crawl and more",
													},
												},
											},
										},
									},
								},
							},
							{
								Data: map[string]any{
									"type": "response.completed",
									"response": map[string]any{
										"id":         "resp_67cf3390786881908b27489d7e8cfb6b",
										"object":     "response",
										"created_at": 1741632400,
										"status":     "completed",
										"model":      "gpt-4o-mini-2024-07-18",
										"output": []any{
											map[string]any{
												"type":   "web_search_call",
												"id":     "ws_67cf3390e9608190869b5d45698a7067",
												"status": "completed",
											},
											map[string]any{
												"type":   "message",
												"id":     "msg_67cf33924ea88190b8c12bf68c1f6416",
												"status": "completed",
												"role":   "assistant",
												"content": []any{
													map[string]any{
														"type": "output_text",
														"text": "Last week in San Francisco a themed party...",
														"annotations": []any{
															map[string]any{
																"type":        "url_citation",
																"start_index": 383,
																"end_index":   493,
																"url":         "https://www.sftourismtips.com/san-francisco-events-in-march.html?utm_source=chatgpt.com",
																"title":       "San Francisco Events in March 2025: Festivals, Theater & Easter",
															},
															map[string]any{
																"type":        "url_citation",
																"start_index": 630,
																"end_index":   762,
																"url":         "https://www.axios.com/local/san-francisco/2025/03/06/sf-events-march-what-to-do-giants-fanfest?utm_source=chatgpt.com",
																"title":       "SF weekend events: Giants FanFest, crab crawl and more",
															},
														},
													},
												},
											},
										},
										"usage": map[string]any{
											"input_tokens": 327,
											"input_tokens_details": map[string]any{
												"cached_tokens": 0,
											},
											"output_tokens": 834,
											"output_tokens_details": map[string]any{
												"reasoning_tokens": 0,
											},
											"total_tokens": 1161,
										},
									},
								},
							},
						}),
					},
				},
			},
			expectedEvents: []api.StreamEvent{
				&api.ResponseMetadataEvent{
					ID:        "resp_67cf3390786881908b27489d7e8cfb6b",
					ModelID:   "gpt-4o-mini-2024-07-18",
					Timestamp: time.Date(2025, 3, 10, 18, 46, 40, 0, time.UTC),
				},
				&api.TextDeltaEvent{
					TextDelta: "Last week",
				},
				&api.TextDeltaEvent{
					TextDelta: " in San Francisco",
				},
				&api.SourceEvent{
					Source: api.Source{
						ID:         "source-0",
						SourceType: "url",
						Title:      "San Francisco Events in March 2025: Festivals, Theater & Easter",
						URL:        "https://www.sftourismtips.com/san-francisco-events-in-march.html?utm_source=chatgpt.com",
					},
				},
				&api.TextDeltaEvent{
					TextDelta: " a themed party",
				},
				&api.TextDeltaEvent{
					TextDelta: "([axios.com](https://www.axios.com/local/san-francisco/2025/03/06/sf-events-march-what-to-do-giants-fanfest?utm_source=chatgpt.com))",
				},
				&api.SourceEvent{
					Source: api.Source{
						ID:         "source-1",
						SourceType: "url",
						Title:      "SF weekend events: Giants FanFest, crab crawl and more",
						URL:        "https://www.axios.com/local/san-francisco/2025/03/06/sf-events-march-what-to-do-giants-fanfest?utm_source=chatgpt.com",
					},
				},
				&api.TextDeltaEvent{
					TextDelta: ".",
				},
				&api.FinishEvent{
					FinishReason: api.FinishReasonStop,
					Usage: api.Usage{
						InputTokens:       327,
						OutputTokens:      834,
						TotalTokens:       1161,
						ReasoningTokens:   0,
						CachedInputTokens: 0,
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{
						"openai": &Metadata{
							ResponseID: "resp_67cf3390786881908b27489d7e8cfb6b",
							Usage: codec.Usage{
								InputTokens:           327,
								OutputTokens:          834,
								InputCachedTokens:     0,
								OutputReasoningTokens: 0,
							},
						},
					}),
				},
			},
		},
		{
			name:    "should stream reasoning summary",
			modelID: "o3-mini-2025-01-31",
			prompt:  standardPrompt,
			options: api.CallOptions{
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"openai": &Metadata{
						ReasoningEffort:  "low",
						ReasoningSummary: "auto",
					},
				}),
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/responses",
						Body: `{
							"model": "o3-mini-2025-01-31",
							"input": [
								{
									"role": "user",
									"content": [
										{
											"type": "input_text",
											"text": "Hello"
										}
									]
								}
							],
							"reasoning": {
								"effort": "low",
								"summary": "auto"
							},
							"stream": true
						}`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Headers: map[string]string{
							"Content-Type": "text/event-stream",
						},
						Body: eventsToString([]sse.Event{
							{
								Data: map[string]any{
									"type": "response.created",
									"response": map[string]any{
										"id":         "resp_67c9a81b6a048190a9ee441c5755a4e8",
										"object":     "response",
										"created_at": 1741269019,
										"status":     "in_progress",
										"model":      "o3-mini-2025-01-31",
										"reasoning":  map[string]any{"effort": "low", "summary": "auto"},
									},
								},
							},
							{Data: map[string]any{"type": "response.reasoning_summary_text.delta", "item_id": "rs_68082c0556348191af675cee0453109b", "output_index": 0, "summary_index": 0, "delta": "**Exploring burrito origins**\n\nThe user is"}},
							{Data: map[string]any{"type": "response.reasoning_summary_text.delta", "item_id": "rs_68082c0556348191af675cee0453109b", "output_index": 0, "summary_index": 0, "delta": " curious about the debate regarding Taqueria La Cumbre and El Farolito."}},
							{Data: map[string]any{"type": "response.reasoning_summary_text.done", "item_id": "rs_68082c0556348191af675cee0453109b", "output_index": 0, "summary_index": 0, "text": "**Exploring burrito origins**\n\nThe user is curious about the debate regarding Taqueria La Cumbre and El Farolito."}},
							{Data: map[string]any{"type": "response.reasoning_summary_text.delta", "item_id": "rs_68082c0556348191af675cee0453109b", "output_index": 0, "summary_index": 1, "delta": "**Investigating burrito origins**\n\nThere's a fascinating debate about who created the Mission burrito."}},
							{Data: map[string]any{"type": "response.reasoning_summary_part.done", "item_id": "rs_68082c0556348191af675cee0453109b", "output_index": 0, "summary_index": 1, "part": map[string]any{"type": "summary_text", "text": "**Investigating burrito origins**\n\nThere's a fascinating debate about who created the Mission burrito."}}},
							{Data: map[string]any{"type": "response.output_item.added", "output_index": 1, "item": map[string]any{"id": "msg_67c9a81dea8c8190b79651a2b3adf91e", "type": "message", "status": "in_progress", "role": "assistant", "content": []any{}}}},
							{Data: map[string]any{"type": "response.content_part.added", "item_id": "msg_67c9a81dea8c8190b79651a2b3adf91e", "output_index": 1, "content_index": 0, "part": map[string]any{"type": "output_text", "text": "", "annotations": []any{}}}},
							{Data: map[string]any{"type": "response.output_text.delta", "item_id": "msg_67c9a81dea8c8190b79651a2b3adf91e", "output_index": 1, "content_index": 0, "delta": "Taqueria La Cumbre"}},
							{
								Data: map[string]any{
									"type": "response.completed",
									"response": map[string]any{
										"id":         "resp_67c9a81b6a048190a9ee441c5755a4e8",
										"object":     "response",
										"created_at": 1741269019,
										"status":     "completed",
										"model":      "o3-mini-2025-01-31",
										"output": []any{
											map[string]any{
												"id": "rs_68082c0556348191af675cee0453109b", "type": "reasoning",
												"summary": []any{
													map[string]any{"type": "summary_text", "text": "**Exploring burrito origins**\n\nThe user is curious about the debate regarding Taqueria La Cumbre and El Farolito."},
													map[string]any{"type": "summary_text", "text": "**Investigating burrito origins**\n\nThere's a fascinating debate about who created the Mission burrito."},
												},
											},
											map[string]any{
												"id": "msg_67c9a81dea8c8190b79651a2b3adf91e", "type": "message", "status": "completed", "role": "assistant",
												"content": []any{map[string]any{"type": "output_text", "text": "Taqueria La Cumbre", "annotations": []any{}}},
											},
										},
										"reasoning": map[string]any{"effort": "low", "summary": "auto"},
										"usage": map[string]any{
											"input_tokens":          543,
											"input_tokens_details":  map[string]any{"cached_tokens": 234},
											"output_tokens":         478,
											"output_tokens_details": map[string]any{"reasoning_tokens": 350},
											"total_tokens":          1021,
										},
									},
								},
							},
						}),
					},
				},
			},
			expectedEvents: []api.StreamEvent{
				&api.ResponseMetadataEvent{
					ID:        "resp_67c9a81b6a048190a9ee441c5755a4e8",
					ModelID:   "o3-mini-2025-01-31",
					Timestamp: time.Date(2025, 3, 6, 13, 50, 19, 0, time.UTC),
				},
				&api.ReasoningEvent{
					TextDelta: "**Exploring burrito origins**\n\nThe user is",
				},
				&api.ReasoningEvent{
					TextDelta: " curious about the debate regarding Taqueria La Cumbre and El Farolito.",
				},
				&api.ReasoningEvent{
					TextDelta: "**Investigating burrito origins**\n\nThere's a fascinating debate about who created the Mission burrito.",
				},
				&api.TextDeltaEvent{
					TextDelta: "Taqueria La Cumbre",
				},
				&api.FinishEvent{
					FinishReason: api.FinishReasonStop,
					Usage: api.Usage{
						InputTokens:       543,
						OutputTokens:      478,
						TotalTokens:       1021,
						ReasoningTokens:   350,
						CachedInputTokens: 234,
					},
					ProviderMetadata: api.NewProviderMetadata(map[string]any{
						"openai": &Metadata{
							ResponseID: "resp_67c9a81b6a048190a9ee441c5755a4e8",
							Usage: codec.Usage{
								InputTokens:           543,
								OutputTokens:          478,
								InputCachedTokens:     234,
								OutputReasoningTokens: 350,
							},
						},
					}),
				},
			},
		},
	}

	runStreamTests(t, tests)
}

func runGenerateTests(t *testing.T, tests []struct {
	name         string
	modelID      string
	options      api.CallOptions
	prompt       []api.Message
	exchanges    []httpmock.Exchange
	wantErr      bool
	expectedResp *api.Response
	skip         bool
},
) {
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.skip {
				t.Skipf("Skipping test: %s", testCase.name)
			}

			server := httpmock.NewServer(t, testCase.exchanges)
			defer server.Close()

			// Set up client options for the OpenAI client
			clientOptions := []option.RequestOption{
				option.WithBaseURL(server.BaseURL()),
				option.WithAPIKey("test-key"),
				option.WithMaxRetries(0), // Disable retries
			}

			// Create client with options
			client := openai.NewClient(clientOptions...)

			// Use custom model ID
			modelID := testCase.modelID

			// Create model with mocked client
			model := NewLanguageModel(modelID, WithClient(client))

			// Call Generate with the test's options (or empty if not specified)
			resp, err := model.Generate(t.Context(), testCase.prompt, testCase.options)

			if testCase.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			// Use aitesting.ResponseContains to verify expected response fields
			aitesting.ResponseContains(t, testCase.expectedResp, resp)
		})
	}
}

func runStreamTests(t *testing.T, tests []struct {
	name           string
	modelID        string
	options        api.CallOptions
	prompt         []api.Message
	exchanges      []httpmock.Exchange
	wantErr        bool
	expectedEvents []api.StreamEvent
	skip           bool
},
) {
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.skip {
				t.Skipf("Skipping test: %s", testCase.name)
			}

			server := httpmock.NewServer(t, testCase.exchanges)
			defer server.Close()

			// Set up client options for the OpenAI client
			clientOptions := []option.RequestOption{
				option.WithBaseURL(server.BaseURL()),
				option.WithAPIKey("test-key"),
				option.WithMaxRetries(0), // Disable retries
			}

			// Create client with options
			client := openai.NewClient(clientOptions...)

			// Use custom model ID
			modelID := testCase.modelID

			// Create model with mocked client
			model := NewLanguageModel(modelID, WithClient(client))

			// Call Stream with the test's options (or empty if not specified)
			resp, err := model.Stream(t.Context(), testCase.prompt, testCase.options)

			if testCase.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			// Collect all events from the stream
			var gotEvents []api.StreamEvent
			for event := range resp.Stream {
				gotEvents = append(gotEvents, event)
			}

			// Compare events using deep equality
			require.Equal(t, testCase.expectedEvents, gotEvents)
		})
	}
}
