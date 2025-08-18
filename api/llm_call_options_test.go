package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallOptions_JSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
	}{
		{
			name: "basic_fields",
			jsonStr: `{
				"max_output_tokens": 1500,
				"temperature": 0.7,
				"stop_sequences": ["STOP", "END"],
				"top_p": 0.9,
				"top_k": 40,
				"presence_penalty": 0.1,
				"frequency_penalty": 0.2,
				"seed": 12345,
				"response_format": {
					"type": "json",
					"name": "response",
					"description": "A structured response"
				},
				"headers": {
					"X-Custom-Header": ["custom-value"],
					"Authorization": ["Bearer token"]
				}
			}`,
		},
		{
			name: "all_fields_comprehensive",
			jsonStr: `{
				"max_output_tokens": 4000,
				"temperature": 0.8,
				"stop_sequences": ["STOP", "END", "FINISH"],
				"top_p": 0.95,
				"top_k": 50,
				"presence_penalty": 0.3,
				"frequency_penalty": 0.4,
				"response_format": {
					"type": "json",
					"schema": {
						"type": "object",
						"properties": {
							"name": {"type": "string"},
							"age": {"type": "integer"}
						},
						"required": ["name"],
						"additionalProperties": false
					},
					"name": "user_info",
					"description": "User information structure"
				},
				"seed": 67890,
				"tools": [
					{
						"type": "function",
						"name": "get_weather",
						"description": "Get weather information",
						"input_schema": {
							"type": "object",
							"properties": {
								"location": {"type": "string"},
								"units": {"type": "string", "enum": ["metric", "imperial"]}
							},
							"required": ["location"],
							"additionalProperties": false
						}
					},
					{
						"type": "provider-defined",
						"id": "openai.file_search",
						"name": "file_search",
						"args": {
							"vector_store_ids": ["store1", "store2"],
							"max_num_results": 20
						}
					}
				],
				"tool_choice": {
					"type": "tool",
					"tool_name": "get_weather"
				},
				"headers": {
					"X-API-Key": ["test-key"],
					"User-Agent": ["test-agent"],
					"Content-Type": ["application/json"]
				},
				"provider_metadata": {
					"openai": {
						"custom_param": "custom_value",
						"extra_config": true
					},
					"anthropic": {
						"model_version": "claude-3"
					}
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal JSON string into CallOptions
			var callOptions CallOptions
			err := json.Unmarshal([]byte(tt.jsonStr), &callOptions)
			require.NoError(t, err, "Failed to unmarshal JSON")

			// Marshal CallOptions back to JSON
			serializedJSON, err := json.Marshal(callOptions)
			require.NoError(t, err, "Failed to marshal CallOptions")

			// Compare original and re-serialized JSON
			assert.JSONEq(t, tt.jsonStr, string(serializedJSON), "JSON round-trip failed")
		})
	}
}

// BenchmarkCallOptionsUnmarshal benchmarks the performance of UnmarshalJSON
func BenchmarkCallOptionsUnmarshal(b *testing.B) {
	// Complex test case with tools (from existing test)
	jsonData := []byte(`{
		"max_output_tokens": 4000,
		"temperature": 0.8,
		"stop_sequences": ["STOP", "END", "FINISH"],
		"top_p": 0.95,
		"top_k": 50,
		"presence_penalty": 0.3,
		"frequency_penalty": 0.4,
		"response_format": {
			"type": "json",
			"schema": {
				"type": "object",
				"properties": {
					"name": {"type": "string"},
					"age": {"type": "integer"}
				},
				"required": ["name"],
				"additionalProperties": {"not": {}}
			},
			"name": "user_info",
			"description": "User information structure"
		},
		"seed": 67890,
		"tools": [
			{
				"type": "function",
				"name": "get_weather",
				"description": "Get weather information",
				"input_schema": {
					"type": "object",
					"properties": {
						"location": {"type": "string"},
						"units": {"type": "string", "enum": ["metric", "imperial"]}
					},
					"required": ["location"],
					"additionalProperties": {"not": {}}
				}
			},
			{
				"type": "provider-defined",
				"id": "openai.file_search",
				"name": "file_search",
				"args": {
					"vector_store_ids": ["store1", "store2"],
					"max_num_results": 20
				}
			}
		],
		"tool_choice": {
			"type": "tool",
			"tool_name": "get_weather"
		},
		"headers": {
			"X-API-Key": ["test-key"],
			"User-Agent": ["test-agent"],
			"Content-Type": ["application/json"]
		},
		"provider_metadata": {
			"openai": {
				"custom_param": "custom_value",
				"extra_config": true
			},
			"anthropic": {
				"model_version": "claude-3"
			}
		}
	}`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var callOptions CallOptions
		err := json.Unmarshal(jsonData, &callOptions)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCallOptionsUnmarshalManyTools tests performance with many tools
func BenchmarkCallOptionsUnmarshalManyTools(b *testing.B) {
	// Generate JSON with many tools to test scaling
	toolsJSON := ""
	for i := 0; i < 50; i++ {
		if i > 0 {
			toolsJSON += ","
		}
		toolsJSON += `{
			"type": "function",
			"name": "tool_` + string(rune('A'+i%26)) + `",
			"description": "Tool description",
			"input_schema": {
				"type": "object",
				"properties": {
					"param": {"type": "string"}
				}
			}
		}`
	}

	jsonData := []byte(`{
		"max_output_tokens": 1000,
		"tools": [` + toolsJSON + `]
	}`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var callOptions CallOptions
		err := json.Unmarshal(jsonData, &callOptions)
		if err != nil {
			b.Fatal(err)
		}
	}
}
