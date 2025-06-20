package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponse_JSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
	}{
		{
			name: "empty_response",
			jsonStr: `{
				"content": [],
				"finish_reason": "stop",
				"usage": {
					"input_tokens": 0,
					"output_tokens": 0,
					"total_tokens": 0
				}
			}`,
		},
		{
			name: "basic_text_response",
			jsonStr: `{
				"content": [
					{
						"type": "text",
						"text": "Hello, world!"
					}
				],
				"finish_reason": "stop",
				"usage": {
					"input_tokens": 10,
					"output_tokens": 15,
					"total_tokens": 25
				}
			}`,
		},
		{
			name: "response_with_usage_details",
			jsonStr: `{
				"content": [
					{
						"type": "text",
						"text": "This is a reasoning response."
					}
				],
				"finish_reason": "stop",
				"usage": {
					"input_tokens": 50,
					"output_tokens": 75,
					"total_tokens": 150,
					"reasoning_tokens": 25,
					"cached_input_tokens": 10
				}
			}`,
		},
		{
			name: "response_with_provider_metadata",
			jsonStr: `{
				"content": [
					{
						"type": "text",
						"text": "Response with metadata"
					}
				],
				"finish_reason": "stop",
				"usage": {
					"input_tokens": 20,
					"output_tokens": 30,
					"total_tokens": 50
				},
				"provider_metadata": {
					"openai": {
						"model": "gpt-4",
						"response_id": "resp_123456"
					},
					"anthropic": {
						"model_version": "claude-3",
						"billing_tier": "premium"
					}
				}
			}`,
		},
		{
			name: "response_with_request_info",
			jsonStr: `{
				"content": [
					{
						"type": "text",
						"text": "Response with request info"
					}
				],
				"finish_reason": "length",
				"usage": {
					"input_tokens": 100,
					"output_tokens": 200,
					"total_tokens": 300
				},
				"request": {
					"body": "eyJtb2RlbCI6ImdwdC00In0="
				}
			}`,
		},
		{
			name: "response_with_response_info",
			jsonStr: `{
				"content": [
					{
						"type": "text",
						"text": "Response with response info"
					}
				],
				"finish_reason": "stop",
				"usage": {
					"input_tokens": 40,
					"output_tokens": 60,
					"total_tokens": 100
				},
				"response": {
					"id": "resp_abc123",
					"timestamp": "2024-01-15T10:30:00Z",
					"model_id": "gpt-4-turbo",
					"status": "200 OK",
					"status_code": 200,
					"body": "eyJyZXNwb25zZSI6InRlc3QifQ=="
				}
			}`,
		},
		{
			name: "response_with_headers",
			jsonStr: `{
				"content": [
					{
						"type": "text",
						"text": "Response with HTTP headers"
					}
				],
				"finish_reason": "stop",
				"usage": {
					"input_tokens": 30,
					"output_tokens": 45,
					"total_tokens": 75
				},
				"response": {
					"id": "resp_headers_456",
					"timestamp": "2024-01-15T11:45:00Z",
					"model_id": "gpt-4-turbo",
					"headers": {
						"Content-Type": ["application/json"],
						"X-Request-ID": ["req_123456789"],
						"X-RateLimit-Remaining": ["99"],
						"Cache-Control": ["no-cache", "no-store"]
					},
					"status": "200 OK",
					"status_code": 200,
					"body": "eyJyZXNwb25zZSI6ImhlYWRlcnMifQ=="
				}
			}`,
		},
		{
			name: "response_with_warnings",
			jsonStr: `{
				"content": [
					{
						"type": "text",
						"text": "Response with warnings"
					}
				],
				"finish_reason": "stop",
				"usage": {
					"input_tokens": 25,
					"output_tokens": 35,
					"total_tokens": 60
				},
				"warnings": [
					{
						"type": "unsupported-setting",
						"setting": "temperature",
						"details": "Temperature is not supported for reasoning models",
						"message": "The temperature setting was ignored"
					},
					{
						"type": "other",
						"message": "Model performed auto-truncation on input"
					}
				]
			}`,
		},
		{
			name: "response_with_tool_calls",
			jsonStr: `{
				"content": [
					{
						"type": "tool-call",
						"tool_call_id": "call_123",
						"tool_name": "get_weather",
						"args": {
							"location": "San Francisco",
							"units": "metric"
						}
					}
				],
				"finish_reason": "tool-calls",
				"usage": {
					"input_tokens": 80,
					"output_tokens": 40,
					"total_tokens": 120
				}
			}`,
		},
		{
			name: "exhaustive_response_all_fields",
			jsonStr: `{
				"content": [
					{
						"type": "text",
						"text": "Comprehensive response with all possible content types.",
						"provider_metadata": {
							"openai": {
								"token_count": 12
							}
						}
					},
					{
						"type": "reasoning",
						"text": "Let me think about this step by step:\n1. First, I need to understand the query\n2. Then provide a comprehensive answer",
						"signature": "sig_reasoning_abc123",
						"provider_metadata": {
							"openai": {
								"reasoning_tokens": 45,
								"verified": true
							}
						}
					},
					{
						"type": "tool-call",
						"tool_call_id": "call_weather_456",
						"tool_name": "get_weather",
						"args": {
							"location": "New York",
							"units": "imperial",
							"include_forecast": true
						},
						"provider_metadata": {
							"openai": {
								"function_call_id": "fc_789",
								"parallel_execution": false
							}
						}
					},
					{
						"type": "image",
						"url": "https://example.com/chart.png",
						"media_type": "image/png",
						"provider_metadata": {
							"dalle": {
								"generation_id": "img_123",
								"style": "photographic"
							}
						}
					},
					{
						"type": "file",
						"filename": "report.pdf",
						"data": "JVBERi0xLjQKJcOkw7zDtsOkwrw=",
						"media_type": "application/pdf",
						"provider_metadata": {
							"document_generator": {
								"pages": 2,
								"file_size_bytes": 12345
							}
						}
					}
				],
				"finish_reason": "stop",
				"usage": {
					"input_tokens": 150,
					"output_tokens": 250,
					"total_tokens": 425,
					"reasoning_tokens": 45,
					"cached_input_tokens": 25
				},
				"provider_metadata": {
					"openai": {
						"model": "gpt-4-turbo",
						"response_id": "resp_comprehensive_789",
						"usage": {
							"input_tokens": 150,
							"output_tokens": 250,
							"input_cached_tokens": 25,
							"output_reasoning_tokens": 45
						}
					},
					"anthropic": {
						"model_version": "claude-3-5-sonnet",
						"billing_tier": "pro",
						"cache_control": {
							"type": "ephemeral"
						}
					}
				},
				"request": {
					"body": "eyJtb2RlbCI6ImdwdC00LXR1cmJvIiwibWVzc2FnZXMiOlt7InJvbGUiOiJ1c2VyIiwiY29udGVudCI6IkhlbGxvIn1dfQ=="
				},
				"response": {
					"id": "resp_comprehensive_789",
					"timestamp": "2024-01-15T14:25:30Z",
					"model_id": "gpt-4-turbo-2024-04-09",
					"status": "200 OK",
					"status_code": 200,
					"body": "eyJpZCI6InJlc3BfY29tcHJlaGVuc2l2ZV83ODkiLCJvYmplY3QiOiJjaGF0LmNvbXBsZXRpb24ifQ=="
				},
				"warnings": [
					{
						"type": "unsupported-setting",
						"setting": "top_k",
						"details": "TopK is not supported by this model",
						"message": "The top_k parameter was ignored"
					},
					{
						"type": "other",
						"details": "This is a generic warning",
						"message": "Generic warning message"
					}
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal JSON string into Response
			var response Response
			err := json.Unmarshal([]byte(tt.jsonStr), &response)
			require.NoError(t, err, "Failed to unmarshal JSON")

			// Marshal Response back to JSON
			serializedJSON, err := json.Marshal(response)
			require.NoError(t, err, "Failed to marshal Response")

			// Compare original and re-serialized JSON
			assert.JSONEq(t, tt.jsonStr, string(serializedJSON), "JSON round-trip failed")
		})
	}
}
