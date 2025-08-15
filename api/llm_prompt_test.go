package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemMessage_JSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
	}{
		{
			name: "basic",
			jsonStr: `{
				"role": "system",
				"content": "You are a helpful assistant."
			}`,
		},
		{
			name: "with_provider_metadata",
			jsonStr: `{
				"role": "system",
				"content": "You are a helpful assistant with special instructions.",
				"provider_metadata": {
					"openai": {
						"status": "completed"
					},
					"anthropic": {
						"cache_control": {"type": "ephemeral"}
					}
				}
			}`,
		},
		{
			name: "empty_content",
			jsonStr: `{
				"role": "system",
				"content": ""
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal JSON string into SystemMessage
			var systemMsg *SystemMessage
			err := json.Unmarshal([]byte(tt.jsonStr), &systemMsg)
			require.NoError(t, err, "Failed to unmarshal JSON")

			// Marshal SystemMessage back to JSON
			serializedJSON, err := json.Marshal(systemMsg)
			require.NoError(t, err, "Failed to marshal SystemMessage")

			// Compare original and re-serialized JSON
			assert.JSONEq(t, tt.jsonStr, string(serializedJSON), "JSON round-trip failed")
		})
	}
}

func TestUserMessage_JSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
	}{
		{
			name: "empty_content",
			jsonStr: `{
				"role": "user",
				"content": []
			}`,
		},
		{
			name: "with_provider_metadata",
			jsonStr: `{
				"role": "user",
				"content": [],
				"provider_metadata": {
					"openai": {
						"status": "completed"
					},
					"anthropic": {
						"cache_control": {"type": "ephemeral"}
					}
				}
			}`,
		},
		{
			name: "exhaustive_all_content_types",
			jsonStr: `{
				"role": "user",
				"content": [
					{
						"type": "text",
						"text": "Hello, this is a text block with unicode 🌟 and special chars <>&\"'",
						"provider_metadata": {
							"openai": {
								"token_count": 15
							}
						}
					},
					{
						"type": "image",
						"url": "https://example.com/image.jpg",
						"media_type": "image/jpeg",
						"provider_metadata": {
							"anthropic": {
								"vision_model": "claude-3-5-sonnet"
							}
						}
					},
					{
						"type": "image",
						"data": "aGVsbG8gd29ybGQ=",
						"media_type": "image/png"
					},
					{
						"type": "file",
						"filename": "document.pdf",
						"url": "https://example.com/doc.pdf",
						"media_type": "application/pdf",
						"provider_metadata": {
							"custom": {
								"file_size": 1024000,
								"pages": 42
							}
						}
					},
					{
						"type": "file",
						"filename": "data.csv",
						"data": "bmFtZSxhZ2UKSm9obiwzMA==",
						"media_type": "text/csv"
					},
					{
						"type": "reasoning",
						"text": "Let me think through this step by step:\n1. First I need to analyze the request\n2. Then I should consider the constraints\n3. Finally I'll provide my response",
						"signature": "sig_reasoning_12345",
						"provider_metadata": {
							"openai": {
								"reasoning_tokens": 150,
								"verified": true
							}
						}
					},
					{
						"type": "redacted-reasoning",
						"data": "redacted_reasoning_data_xyz789",
						"provider_metadata": {
							"anthropic": {
								"redaction_level": "high"
							}
						}
					},
					{
						"type": "tool-call",
						"tool_call_id": "call_abc123",
						"tool_name": "get_weather",
						"args": {
							"location": "San Francisco",
							"units": "metric",
							"include_forecast": true
						},
						"provider_metadata": {
							"openai": {
								"function_call_id": "fc_xyz789",
								"parallel_execution": true
							}
						}
					},
					{
						"type": "tool-result",
						"tool_call_id": "call_abc123",
						"tool_name": "get_weather",
						"result": {
							"temperature": 18,
							"condition": "sunny",
							"humidity": 65,
							"forecast": [
								{"day": "tomorrow", "temp": 20, "condition": "cloudy"}
							]
						},
						"provider_metadata": {
							"weather_api": {
								"api_version": "v2.1",
								"cache_hit": true,
								"response_time_ms": 120
							}
						}
					},
					{
						"type": "source",
						"id": "source_123",
						"url": "https://example.com/source",
						"title": "Important Reference Document",
						"provider_metadata": {
							"search_engine": {
								"relevance_score": 0.95,
								"crawl_date": "2024-01-15"
							}
						}
					}
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal JSON string into UserMessage
			var userMsg *UserMessage
			err := json.Unmarshal([]byte(tt.jsonStr), &userMsg)
			require.NoError(t, err, "Failed to unmarshal JSON")

			// Marshal UserMessage back to JSON
			serializedJSON, err := json.Marshal(userMsg)
			require.NoError(t, err, "Failed to marshal UserMessage")

			// Compare original and re-serialized JSON
			assert.JSONEq(t, tt.jsonStr, string(serializedJSON), "JSON round-trip failed")
		})
	}
}

func TestAssistantMessage_JSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
	}{
		{
			name: "empty_content",
			jsonStr: `{
				"role": "assistant",
				"content": []
			}`,
		},
		{
			name: "with_provider_metadata",
			jsonStr: `{
				"role": "assistant",
				"content": [],
				"provider_metadata": {
					"anthropic": {
						"model_version": "claude-3"
					}
				}
			}`,
		},
		{
			name: "exhaustive_all_content_types",
			jsonStr: `{
				"role": "assistant",
				"content": [
					{
						"type": "text",
						"text": "I'll help you with that! Here's my response with unicode ✨ and special chars <>&\"' properly handled.",
						"provider_metadata": {
							"anthropic": {
								"token_count": 18,
								"confidence": 0.95
							}
						}
					},
					{
						"type": "reasoning",
						"text": "Let me think through this problem step by step:\n1. I need to understand the user's request\n2. Consider the constraints and requirements\n3. Formulate a comprehensive response\n4. Ensure accuracy and completeness",
						"signature": "sig_assistant_reasoning_67890",
						"provider_metadata": {
							"openai": {
								"reasoning_tokens": 200,
								"reasoning_model": "o1-preview",
								"verified": true
							}
						}
					},
					{
						"type": "redacted-reasoning",
						"data": "redacted_assistant_reasoning_abc123",
						"provider_metadata": {
							"anthropic": {
								"redaction_level": "medium",
								"redaction_reason": "privacy"
							}
						}
					},
					{
						"type": "tool-call",
						"tool_call_id": "call_def456",
						"tool_name": "calculate_math",
						"args": {
							"expression": "2 + 2 * 3",
							"precision": 2,
							"show_steps": true
						},
						"provider_metadata": {
							"openai": {
								"function_call_id": "fc_math_789",
								"execution_priority": "high"
							}
						}
					},
					{
						"type": "tool-call",
						"tool_call_id": "call_ghi789",
						"tool_name": "search_knowledge",
						"args": {
							"query": "latest AI research 2024",
							"max_results": 5,
							"filter": {
								"date_range": "2024-01-01:2024-12-31",
								"domains": ["arxiv.org", "openai.com"]
							}
						},
						"provider_metadata": {
							"custom": {
								"search_engine": "semantic",
								"cache_enabled": false
							}
						}
					},
					{
						"type": "image",
						"url": "https://example.com/generated-chart.png",
						"media_type": "image/png",
						"provider_metadata": {
							"dalle": {
								"generation_id": "img_gen_123",
								"model": "dall-e-3",
								"style": "natural"
							}
						}
					},
					{
						"type": "file",
						"filename": "analysis_report.pdf",
						"data": "JVBERi0xLjQKJcOkw7zDtsOkwrw=",
						"media_type": "application/pdf",
						"provider_metadata": {
							"document_generator": {
								"pages": 3,
								"word_count": 1500,
								"format_version": "2.0"
							}
						}
					},
					{
						"type": "source",
						"id": "assistant_source_456",
						"url": "https://example.com/reference-doc",
						"title": "Supporting Documentation for Analysis",
						"provider_metadata": {
							"knowledge_base": {
								"relevance_score": 0.88,
								"last_updated": "2024-01-20",
								"citation_count": 42
							}
						}
					}
				],
				"provider_metadata": {
					"openai": {
						"status": "completed"
					},
					"anthropic": {
						"cache_control": {"type": "ephemeral"}
					}
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal JSON string into AssistantMessage
			var assistantMsg *AssistantMessage
			err := json.Unmarshal([]byte(tt.jsonStr), &assistantMsg)
			require.NoError(t, err, "Failed to unmarshal JSON")

			// Marshal AssistantMessage back to JSON
			serializedJSON, err := json.Marshal(assistantMsg)
			require.NoError(t, err, "Failed to marshal AssistantMessage")

			// Compare original and re-serialized JSON
			assert.JSONEq(t, tt.jsonStr, string(serializedJSON), "JSON round-trip failed")
		})
	}
}

func TestToolMessage_JSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
	}{
		{
			name: "basic_tool_result",
			jsonStr: `{
				"role": "tool",
				"content": [
					{
						"type": "tool-result",
						"tool_call_id": "call_123",
						"tool_name": "get_weather",
						"result": "It's sunny and 72°F"
					}
				]
			}`,
		},
		{
			name: "with_provider_metadata",
			jsonStr: `{
				"role": "tool",
				"content": [
					{
						"type": "tool-result",
						"tool_call_id": "call_123",
						"tool_name": "get_weather",
						"result": "It's sunny and 72°F"
					}
				],
				"provider_metadata": {
					"openai": {
						"execution_time_ms": 150
					}
				}
			}`,
		},
		{
			name: "empty_content",
			jsonStr: `{
				"role": "tool",
				"content": []
			}`,
		},
		{
			name: "exhaustive_all_tool_result_types",
			jsonStr: `{
				"role": "tool",
				"content": [
					{
						"type": "tool-result",
						"tool_call_id": "call_simple_001",
						"tool_name": "get_weather",
						"result": "It's sunny and 72°F in San Francisco",
						"provider_metadata": {
							"weather_api": {
								"api_version": "v2.1",
								"response_time_ms": 150,
								"cache_hit": true
							}
						}
					},
					{
						"type": "tool-result",
						"tool_call_id": "call_complex_002",
						"tool_name": "analyze_data",
						"result": {
							"summary": "Analysis complete",
							"total_records": 1000,
							"averages": {
								"score": 85.5,
								"confidence": 0.92
							},
							"categories": ["positive", "neutral", "negative"],
							"metadata": {
								"processing_time": "2.3s",
								"algorithm": "advanced_ml",
								"version": "1.2.0"
							}
						},
						"provider_metadata": {
							"analytics_engine": {
								"compute_units": 45,
								"memory_used_mb": 256,
								"optimization_level": "high"
							}
						}
					},
					{
						"type": "tool-result",
						"tool_call_id": "call_error_003",
						"tool_name": "fetch_user_data",
						"result": {
							"error_code": "AUTH_FAILED",
							"error_message": "Invalid API key provided",
							"details": {
								"timestamp": "2024-01-15T10:30:00Z",
								"request_id": "req_xyz789"
							}
						},
						"is_error": true,
						"provider_metadata": {
							"auth_service": {
								"validation_attempt": 3,
								"rate_limited": false,
								"suggested_action": "refresh_token"
							}
						}
					},
					{
						"type": "tool-result",
						"tool_call_id": "call_multimodal_004",
						"tool_name": "generate_report",
						"result": null,
						"content": [
							{
								"type": "text",
								"text": "Report Generation Summary:\n- Total pages: 5\n- Charts included: 3\n- Data sources: 2",
								"provider_metadata": {
									"report_generator": {
										"template": "executive_summary",
										"language": "en"
									}
								}
							},
							{
								"type": "image",
								"url": "https://example.com/generated-chart-1.png",
								"media_type": "image/png",
								"provider_metadata": {
									"chart_generator": {
										"chart_type": "bar",
										"data_points": 12,
										"style": "professional"
									}
								}
							},
							{
								"type": "file",
								"filename": "quarterly_report.pdf",
								"data": "JVBERi0xLjQKJcOkw7zDtsOkwrw=",
								"media_type": "application/pdf",
								"provider_metadata": {
									"pdf_generator": {
										"pages": 5,
										"file_size_bytes": 245760,
										"compression": "standard"
									}
								}
							},
							{
								"type": "source",
								"id": "data_source_001",
								"url": "https://example.com/quarterly-data",
								"title": "Q4 2024 Financial Data",
								"provider_metadata": {
									"data_warehouse": {
										"last_updated": "2024-01-10",
										"record_count": 5000,
										"confidence": 0.98
									}
								}
							}
						],
						"provider_metadata": {
							"report_service": {
								"generation_id": "rpt_gen_456",
								"template_version": "2.1",
								"processing_time_ms": 5500,
								"resources_used": {
									"cpu_seconds": 12.5,
									"memory_peak_mb": 512
								}
							}
						}
					},
					{
						"type": "tool-result",
						"tool_call_id": "call_array_result_005",
						"tool_name": "list_files",
						"result": [
							{
								"name": "document1.txt",
								"size": 1024,
								"modified": "2024-01-15T09:00:00Z"
							},
							{
								"name": "image.jpg",
								"size": 2048576,
								"modified": "2024-01-14T15:30:00Z"
							},
							{
								"name": "data.csv",
								"size": 5120,
								"modified": "2024-01-13T11:45:00Z"
							}
						],
						"provider_metadata": {
							"file_system": {
								"directory": "/user/documents",
								"total_files": 3,
								"scan_time_ms": 25
							}
						}
					},
					{
						"type": "tool-result",
						"tool_call_id": "call_numeric_006",
						"tool_name": "calculate_pi",
						"result": 3.141592653589793,
						"provider_metadata": {
							"math_engine": {
								"precision": 15,
								"algorithm": "chudnovsky",
								"computation_time_ms": 1
							}
						}
					}
				],
				"provider_metadata": {
					"openai": {
						"status": "completed"
					},
					"anthropic": {
						"cache_control": {"type": "ephemeral"}
					}
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal JSON string into ToolMessage
			var toolMsg *ToolMessage
			err := json.Unmarshal([]byte(tt.jsonStr), &toolMsg)
			require.NoError(t, err, "Failed to unmarshal JSON")

			// Marshal ToolMessage back to JSON
			serializedJSON, err := json.Marshal(toolMsg)
			require.NoError(t, err, "Failed to marshal ToolMessage")

			// Compare original and re-serialized JSON
			assert.JSONEq(t, tt.jsonStr, string(serializedJSON), "JSON round-trip failed")
		})
	}
}
