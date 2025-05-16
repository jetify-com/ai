package openai

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/aitesting"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/openai/internal/codec"
	"go.jetify.com/pkg/httpmock"
)

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
	// - Only check text: expectedResp: api.Response{Text: "example"}
	// - Only check usage: expectedResp: api.Response{Usage: api.Usage{PromptTokens: 10}}
	// - Only check finish reason: expectedResp: api.Response{FinishReason: api.FinishReasonStop}
	// aitesting.ResponseContains only checks fields that are specified in expectedResp.
	tests := []struct {
		name         string
		modelID      string          // Optional: override default model
		options      api.CallOptions // Optional: custom call options
		prompt       []api.Message
		exchanges    []httpmock.Exchange
		wantErr      bool
		expectedResp api.Response // Expected response fields
		skip         bool         // Skip this test if true
	}{
		{
			name:      "response contains text",
			modelID:   "gpt-4o",
			prompt:    standardPrompt,
			exchanges: standardExchange,
			expectedResp: api.Response{
				Text: "answer text",
			},
		},
		{
			name:      "response contains usage information",
			modelID:   "gpt-4o",
			prompt:    standardPrompt,
			exchanges: standardExchange,
			expectedResp: api.Response{
				Usage: api.Usage{
					PromptTokens:     345,
					CompletionTokens: 538,
				},
			},
		},
		{
			name:      "response contains provider metadata",
			modelID:   "gpt-4o",
			prompt:    standardPrompt,
			exchanges: standardExchange,
			expectedResp: api.Response{
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
				Temperature: float64Ptr(0.5),
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
			expectedResp: api.Response{
				Text: "answer text",
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
				Temperature: float64Ptr(0.5),
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
			expectedResp: api.Response{
				Text: "answer text",
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
				Temperature: float64Ptr(0.5),
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
			expectedResp: api.Response{
				Text: "answer text",
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
					Schema: &jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"value": {Type: jsonschema.String},
						},
						Required:             []string{"value"},
						AdditionalProperties: false,
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
			expectedResp: api.Response{
				Text:     "answer text",
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
			expectedResp: api.Response{
				Text:     "answer text",
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
						ParallelToolCalls: boolPtr(false),
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
			expectedResp: api.Response{
				Text:     "answer text",
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
						Store: boolPtr(false),
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
			expectedResp: api.Response{
				Text:     "answer text",
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
			expectedResp: api.Response{
				Text:     "answer text",
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
			expectedResp: api.Response{
				Text:     "answer text",
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
			expectedResp: api.Response{
				Text:     "answer text",
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
			expectedResp: api.Response{
				Text:     "answer text",
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
			expectedResp: api.Response{
				Text:     "answer text",
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends object-tool format",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Mode: api.ObjectToolMode{
					Tool: api.FunctionTool{
						Name:        "response",
						Description: "A response",
						InputSchema: &jsonschema.Definition{
							Type: jsonschema.Object,
							Properties: map[string]jsonschema.Definition{
								"value": {Type: jsonschema.String},
							},
							Required:             []string{"value"},
							AdditionalProperties: false,
						},
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
			expectedResp: api.Response{
				Text:     "answer text",
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends object-json json_object format",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Mode: api.ObjectJSONMode{},
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
			expectedResp: api.Response{
				Text:     "answer text",
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends object-json json_schema format",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Mode: api.ObjectJSONMode{
					Name:        "response",
					Description: "A response",
					Schema: &jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"value": {Type: jsonschema.String},
						},
						Required:             []string{"value"},
						AdditionalProperties: false,
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
			expectedResp: api.Response{
				Text:     "answer text",
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends object-json json_schema format with strictSchemas false",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Mode: api.ObjectJSONMode{
					Name:        "response",
					Description: "A response",
					Schema: &jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"value": {Type: jsonschema.String},
						},
						Required:             []string{"value"},
						AdditionalProperties: false,
					},
				},
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"openai": &Metadata{
						StrictSchemas: boolPtr(false),
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
			expectedResp: api.Response{
				Text:     "answer text",
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends web_search tool",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Mode: api.RegularMode{
					Tools: []api.ToolDefinition{
						&codec.WebSearchTool{
							SearchContextSize: "high",
							UserLocation: &codec.WebSearchUserLocation{
								City: "San Francisco",
							},
						},
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
			expectedResp: api.Response{
				Text:     "answer text",
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "sends web_search tool as tool_choice",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Mode: api.RegularMode{
					ToolChoice: &api.ToolChoice{
						Type:     "tool",
						ToolName: "web_search_preview",
					},
					Tools: []api.ToolDefinition{
						&codec.WebSearchTool{
							SearchContextSize: "high",
							UserLocation: &codec.WebSearchUserLocation{
								City: "San Francisco",
							},
						},
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
			expectedResp: api.Response{
				Text:     "answer text",
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "warns about unsupported settings",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Mode:             api.RegularMode{},
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
			expectedResp: api.Response{
				Text: "answer text",
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

	standardTools := []api.ToolDefinition{
		&api.FunctionTool{
			Name: "weather",
			InputSchema: &jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"location": {Type: jsonschema.String},
				},
				Required:             []string{"location"},
				AdditionalProperties: false,
			},
		},
		&api.FunctionTool{
			Name: "cityAttractions",
			InputSchema: &jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"city": {Type: jsonschema.String},
				},
				Required:             []string{"city"},
				AdditionalProperties: false,
			},
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
		expectedResp api.Response // Expected response fields
		skip         bool         // Skip this test if true
	}{
		{
			name:    "should generate tool calls",
			modelID: "gpt-4o",
			prompt:  standardPrompt,
			options: api.CallOptions{
				Mode: api.RegularMode{
					Tools: standardTools,
				},
			},
			exchanges: standardExchange,
			expectedResp: api.Response{
				ToolCalls: []api.ToolCallBlock{
					{
						ToolCallID: "call_0NdsJqOS8N3J9l2p0p4WpYU9",
						ToolName:   "weather",
						Args:       json.RawMessage(`{"location":"San Francisco"}`),
					},
					{
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
				Mode: api.RegularMode{
					Tools: standardTools,
				},
			},
			exchanges: standardExchange,
			expectedResp: api.Response{
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
		expectedResp api.Response // Expected response fields
		skip         bool         // Skip this test if true
	}{
		{
			name:      "should generate text with web search",
			modelID:   "gpt-4o",
			prompt:    standardPrompt,
			exchanges: standardExchange,
			options: api.CallOptions{
				Mode: api.RegularMode{
					Tools: []api.ToolDefinition{
						&codec.WebSearchTool{
							SearchContextSize: "medium",
							UserLocation: &codec.WebSearchUserLocation{
								Country: "US",
							},
						},
					},
				},
			},
			expectedResp: api.Response{
				Text: expectedText,
			},
		},
		{
			name:      "should return sources from web search",
			modelID:   "gpt-4o",
			prompt:    standardPrompt,
			exchanges: standardExchange,
			options: api.CallOptions{
				Mode: api.RegularMode{
					Tools: []api.ToolDefinition{
						&codec.WebSearchTool{
							SearchContextSize: "medium",
							UserLocation: &codec.WebSearchUserLocation{
								Country: "US",
							},
						},
					},
				},
			},
			expectedResp: api.Response{
				Sources: []api.Source{
					{
						SourceType: "url",
						ID:         "source-0",
						URL:        "https://www.axios.com/local/san-francisco/2025/03/07/bruce-lee-statue-sf-chinatown?utm_source=chatgpt.com",
						Title:      "Bruce Lee statue to be installed in SF Chinatown",
					},
					{
						SourceType: "url",
						ID:         "source-1",
						URL:        "https://www.axios.com/local/san-francisco/2025/03/03/bay-area-office-leasing-activity?utm_source=chatgpt.com",
						Title:      "The Bay Area is set to make an office leasing comeback",
					},
					{
						SourceType: "url",
						ID:         "source-2",
						URL:        "https://www.axios.com/local/san-francisco/2025/03/03/where-to-see-spring-blooms-bay-area?utm_source=chatgpt.com",
						Title:      "Where to see spring blooms in the Bay Area",
					},
					{
						SourceType: "url",
						ID:         "source-3",
						URL:        "https://www.axios.com/local/san-francisco/2025/03/03/great-highway-park-opening-april-recall-campaign?utm_source=chatgpt.com",
						Title:      "Oceanfront Great Highway park set to open in April",
					},
					{
						SourceType: "url",
						ID:         "source-4",
						URL:        "https://www.axios.com/local/san-francisco/2025/03/03/climate-weather-spring-temperatures-warmer-sf?utm_source=chatgpt.com",
						Title:      "San Francisco's spring seasons are getting warmer",
					},
				},
			},
		},
	}

	runGenerateTests(t, tests)
}

// Helper function for float64 pointer
func float64Ptr(v float64) *float64 {
	return &v
}

// Helper function for boolean pointer
func boolPtr(v bool) *bool {
	return &v
}

func runGenerateTests(t *testing.T, tests []struct {
	name         string
	modelID      string
	options      api.CallOptions
	prompt       []api.Message
	exchanges    []httpmock.Exchange
	wantErr      bool
	expectedResp api.Response
	skip         bool
},
) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skipf("Skipping test: %s", tt.name)
			}

			server := httpmock.NewServer(t, tt.exchanges)
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
			modelID := tt.modelID

			// Create model with mocked client
			model := NewLanguageModel(modelID, WithClient(client))

			// Call Generate with the test's options (or empty if not specified)
			resp, err := model.Generate(t.Context(), tt.prompt, tt.options)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			// Use aitesting.ResponseContains to verify expected response fields
			aitesting.ResponseContains(t, tt.expectedResp, resp)
		})
	}
}
