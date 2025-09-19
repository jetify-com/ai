package codec

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v2/responses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/aitesting"
	"go.jetify.com/ai/api"
)

func TestDecodeResponse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *api.Response
	}{
		{
			name: "simple message",
			input: `{
				"output": [
					{
						"type": "message",
						"content": [
							{
								"type": "output_text",
								"text": "Hello world"
							}
						]
					}
				],
				"usage": {
					"input_tokens": 100,
					"output_tokens": 50
				}
			}`,
			want: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Hello world"},
				},
				Usage: api.Usage{
					InputTokens:  100,
					OutputTokens: 50,
					TotalTokens:  150,
				},
				FinishReason: api.FinishReasonStop,
			},
		},
		{
			name: "message with sources",
			input: `{
				"output": [
					{
						"type": "message",
						"content": [
							{
								"type": "output_text",
								"text": "Info from [example](https://example.com)",
								"annotations": [
									{
										"type": "url_citation",
										"text": "example",
										"url": "https://example.com",
										"title": "Example Site"
									}
								]
							}
						]
					}
				],
				"usage": {
					"input_tokens": 150,
					"output_tokens": 75
				}
			}`,
			want: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Info from [example](https://example.com)"},
					&api.SourceBlock{
						ID:    "source-0",
						URL:   "https://example.com",
						Title: "Example Site",
					},
				},
				Usage: api.Usage{
					InputTokens:  150,
					OutputTokens: 75,
					TotalTokens:  225,
				},
				FinishReason: api.FinishReasonStop,
			},
		},
		{
			name: "response with tool calls",
			input: `{
				"output": [
					{
						"type": "function_call",
						"call_id": "call1",
						"name": "get_weather",
						"arguments": "{\"location\":\"New York\"}"
					}
				],
				"usage": {
					"input_tokens": 200,
					"output_tokens": 100
				}
			}`,
			want: &api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "call1",
						ToolName:   "get_weather",
						Args:       json.RawMessage(`{"location":"New York"}`),
					},
				},
				Usage: api.Usage{
					InputTokens:  200,
					OutputTokens: 100,
					TotalTokens:  300,
				},
				FinishReason: api.FinishReasonToolCalls,
			},
		},
		{
			name: "response with provider metadata",
			input: `{
				"id": "resp_789",
				"output": [
					{
						"type": "message",
						"content": [
							{
								"type": "output_text",
								"text": "Test message"
							}
						]
					}
				],
				"usage": {
					"input_tokens": 100,
					"output_tokens": 50,
					"input_tokens_details": {
						"cached_tokens": 25
					},
					"output_tokens_details": {
						"reasoning_tokens": 30
					}
				}
			}`,
			want: &api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "Test message"},
				},
				Usage: api.Usage{
					InputTokens:       100,
					OutputTokens:      50,
					TotalTokens:       150,
					ReasoningTokens:   30,
					CachedInputTokens: 25,
				},
				FinishReason: api.FinishReasonStop,
				ProviderMetadata: api.NewProviderMetadata(map[string]any{
					"openai": &Metadata{
						ResponseID: "resp_789",
						Usage: Usage{
							InputTokens:           100,
							OutputTokens:          50,
							InputCachedTokens:     25,
							OutputReasoningTokens: 30,
						},
					},
				}),
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var msg responses.Response
			err := json.Unmarshal([]byte(testCase.input), &msg)
			require.NoError(t, err)

			got, err := DecodeResponse(&msg)
			require.NoError(t, err)

			// Ensure all slice fields are initialized as empty slices:
			if testCase.want.Content == nil {
				testCase.want.Content = []api.ContentBlock{}
			}
			if testCase.want.Warnings == nil {
				testCase.want.Warnings = []api.CallWarning{}
			}

			// Use ResponseContains to verify the response
			aitesting.ResponseContains(t, testCase.want, got)
		})
	}
}

func TestDecodeText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []api.ContentBlock
		wantErr string
	}{
		{
			name: "simple message",
			input: `{
				"output": [
					{
						"type": "message",
						"content": [
							{
								"type": "output_text",
								"text": "Hello world"
							}
						]
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.TextBlock{Text: "Hello world"},
			},
		},
		{
			name: "multiple text segments in same message",
			input: `{
				"output": [
					{
						"type": "message",
						"content": [
							{
								"type": "output_text",
								"text": "First"
							},
							{
								"type": "output_text",
								"text": "Second"
							}
						]
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.TextBlock{Text: "First"},
				&api.TextBlock{Text: "Second"},
			},
		},
		{
			name: "multiple messages with text",
			input: `{
				"output": [
					{
						"type": "message",
						"content": [
							{
								"type": "output_text",
								"text": "Message 1"
							}
						]
					},
					{
						"type": "message",
						"content": [
							{
								"type": "output_text",
								"text": "Message 2"
							}
						]
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.TextBlock{Text: "Message 1"},
				&api.TextBlock{Text: "Message 2"},
			},
		},
		{
			name: "mixed content types",
			input: `{
				"output": [
					{
						"type": "message",
						"content": [
							{
								"type": "output_text",
								"text": "Text before"
							},
							{
								"type": "other_type"
							},
							{
								"type": "output_text",
								"text": "Text after"
							}
						]
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.TextBlock{Text: "Text before"},
				&api.TextBlock{Text: "Text after"},
			},
		},
		{
			name: "empty response",
			input: `{
				"output": []
			}`,
			want: []api.ContentBlock{},
		},
		{
			name:  "nil response",
			input: "null",
			want:  []api.ContentBlock{},
		},
		{
			name: "non-message type",
			input: `{
				"output": [
					{
						"type": "not_a_message",
						"content": [
							{
								"type": "output_text",
								"text": "Hello world"
							}
						]
					}
				]
			}`,
			want:    []api.ContentBlock{},
			wantErr: "unknown output item type: not_a_message",
		},
		{
			name: "non-text content type",
			input: `{
				"output": [
					{
						"type": "message",
						"content": [
							{
								"type": "not_output_text",
								"text": "Hello world"
							}
						]
					}
				]
			}`,
			want: []api.ContentBlock{},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var msg *responses.Response
			if testCase.input != "null" {
				msg = &responses.Response{}
				err := json.Unmarshal([]byte(testCase.input), msg)
				require.NoError(t, err)
			}

			got, err := decodeContent(msg)
			if testCase.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.wantErr)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, testCase.want, got.Content)
		})
	}
}

func TestDecodeSources(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []api.ContentBlock
		wantErr string
	}{
		{
			name: "single url citation",
			input: `{
				"output": [
					{
						"type": "message",
						"content": [
							{
								"type": "output_text",
								"text": "Info from [example](https://example.com)",
								"annotations": [
									{
										"type": "url_citation",
										"text": "example",
										"url": "https://example.com",
										"title": "Example Site"
									}
								]
							}
						]
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.TextBlock{Text: "Info from [example](https://example.com)"},
				&api.SourceBlock{
					ID:    "source-0",
					URL:   "https://example.com",
					Title: "Example Site",
				},
			},
		},
		{
			name: "multiple url citations",
			input: `{
				"output": [
					{
						"type": "message",
						"content": [
							{
								"type": "output_text",
								"text": "Info from [example1](https://example1.com) and [example2](https://example2.com)",
								"annotations": [
									{
										"type": "url_citation",
										"text": "example1",
										"url": "https://example1.com",
										"title": "Example Site 1"
									},
									{
										"type": "url_citation",
										"text": "example2",
										"url": "https://example2.com",
										"title": "Example Site 2"
									}
								]
							}
						]
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.TextBlock{Text: "Info from [example1](https://example1.com) and [example2](https://example2.com)"},
				&api.SourceBlock{
					ID:    "source-0",
					URL:   "https://example1.com",
					Title: "Example Site 1",
				},
				&api.SourceBlock{
					ID:    "source-1",
					URL:   "https://example2.com",
					Title: "Example Site 2",
				},
			},
		},
		{
			name: "non-url annotation",
			input: `{
				"output": [
					{
						"type": "message",
						"content": [
							{
								"type": "output_text",
								"text": "Info from source",
								"annotations": [
									{
										"type": "not_url_citation",
										"text": "source"
									}
								]
							}
						]
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.TextBlock{Text: "Info from source"},
			},
		},
		{
			name: "non-message type",
			input: `{
				"output": [
					{
						"type": "not_a_message",
						"content": [
							{
								"type": "output_text",
								"text": "Info from [example](https://example.com)",
								"annotations": [
									{
										"type": "url_citation",
										"text": "example",
										"url": "https://example.com",
										"title": "Example Site"
									}
								]
							}
						]
					}
				]
			}`,
			want:    []api.ContentBlock{},
			wantErr: "unknown output item type: not_a_message",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var msg responses.Response
			err := json.Unmarshal([]byte(testCase.input), &msg)
			require.NoError(t, err)

			got, err := decodeContent(&msg)
			if testCase.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.wantErr)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, testCase.want, got.Content)
		})
	}
}

func TestDecodeToolCalls(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     []api.ContentBlock
		wantErr  string
	}{
		{
			name: "single function call",
			response: `{
				"output": [
					{
						"type": "function_call",
						"call_id": "call1",
						"name": "get_weather",
						"arguments": "{\"location\":\"New York\"}"
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.ToolCallBlock{
					ToolCallID: "call1",
					ToolName:   "get_weather",
					Args:       json.RawMessage(`{"location":"New York"}`),
				},
			},
		},
		{
			name: "multiple function calls",
			response: `{
				"output": [
					{
						"type": "function_call",
						"call_id": "call1",
						"name": "get_weather",
						"arguments": "{\"location\":\"New York\"}"
					},
					{
						"type": "function_call",
						"call_id": "call2",
						"name": "get_time",
						"arguments": "{\"timezone\":\"EST\"}"
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.ToolCallBlock{
					ToolCallID: "call1",
					ToolName:   "get_weather",
					Args:       json.RawMessage(`{"location":"New York"}`),
				},
				&api.ToolCallBlock{
					ToolCallID: "call2",
					ToolName:   "get_time",
					Args:       json.RawMessage(`{"timezone":"EST"}`),
				},
			},
		},
		{
			name: "file search call",
			response: `{
				"output": [
					{
						"type": "file_search_call",
						"id": "search1",
						"query": "find main.go",
						"include_pattern": "*.go",
						"exclude_pattern": "vendor/*"
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.ToolCallBlock{
					ToolCallID: "search1",
					ToolName:   "openai.file_search",
					Args:       json.RawMessage(`{"query":"find main.go","include_pattern":"*.go","exclude_pattern":"vendor/*","id":"search1","type":"file_search_call"}`),
				},
			},
		},
		{
			name: "web search call",
			response: `{
				"output": [
					{
						"type": "web_search_call",
						"id": "web1",
						"query": "golang best practices",
						"num_results": 5
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.ToolCallBlock{
					ToolCallID: "web1",
					ToolName:   "openai.web_search_preview",
					Args:       json.RawMessage(`{"query":"golang best practices","num_results":5,"id":"web1","type":"web_search_call"}`),
				},
			},
		},
		{
			name: "computer call",
			response: `{
				"output": [
					{
						"type": "computer_call",
						"id": "comp1",
						"command": "ls -la",
						"working_directory": "/tmp"
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.ToolCallBlock{
					ToolCallID: "comp1",
					ToolName:   "openai.computer_use_preview",
					Args:       json.RawMessage(`{"command":"ls -la","working_directory":"/tmp","id":"comp1","type":"computer_call"}`),
					ProviderMetadata: api.NewProviderMetadata(map[string]any{
						"openai": &Metadata{
							ComputerSafetyChecks: []ComputerSafetyCheck{},
						},
					}),
				},
			},
		},
		{
			name: "computer call with safety checks",
			response: `{
				"output": [
					{
						"type": "computer_call",
						"id": "comp2",
						"command": "rm -rf /",
						"working_directory": "/",
						"pending_safety_checks": [
							{
								"id": "check1",
								"code": "file_access",
								"message": "File access requested"
							},
							{
								"id": "check2",
								"code": "dangerous_command",
								"message": "Potentially dangerous command detected"
							}
						]
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.ToolCallBlock{
					ToolCallID: "comp2",
					ToolName:   "openai.computer_use_preview",
					Args:       json.RawMessage(`{"command":"rm -rf /","working_directory":"/","id":"comp2","type":"computer_call","pending_safety_checks":[{"id":"check1","code":"file_access","message":"File access requested"},{"id":"check2","code":"dangerous_command","message":"Potentially dangerous command detected"}]}`),
					ProviderMetadata: api.NewProviderMetadata(map[string]any{
						"openai": &Metadata{
							ComputerSafetyChecks: []ComputerSafetyCheck{
								{
									ID:      "check1",
									Code:    "file_access",
									Message: "File access requested",
								},
								{
									ID:      "check2",
									Code:    "dangerous_command",
									Message: "Potentially dangerous command detected",
								},
							},
						},
					}),
				},
			},
		},
		{
			name: "computer call with empty safety checks",
			response: `{
				"output": [
					{
						"type": "computer_call",
						"id": "comp3",
						"command": "echo hello",
						"working_directory": "/tmp",
						"pending_safety_checks": []
					}
				]
			}`,
			want: []api.ContentBlock{
				&api.ToolCallBlock{
					ToolCallID: "comp3",
					ToolName:   "openai.computer_use_preview",
					Args:       json.RawMessage(`{"command":"echo hello","working_directory":"/tmp","id":"comp3","type":"computer_call","pending_safety_checks":[]}`),
					ProviderMetadata: api.NewProviderMetadata(map[string]any{
						"openai": &Metadata{
							ComputerSafetyChecks: []ComputerSafetyCheck{},
						},
					}),
				},
			},
		},
		{
			name: "unknown tool call type",
			response: `{
				"output": [
					{
						"type": "unknown_call",
						"id": "unknown1"
					}
				]
			}`,
			wantErr: "unknown output item type: unknown_call",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var msg responses.Response
			err := json.Unmarshal([]byte(testCase.response), &msg)
			require.NoError(t, err)

			got, err := decodeContent(&msg)
			if testCase.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.wantErr)
				return
			}
			require.NoError(t, err)

			// Compare content blocks directly
			require.Equal(t, len(testCase.want), len(got.Content))
			for i, wantBlock := range testCase.want {
				if wantToolCall, ok := wantBlock.(*api.ToolCallBlock); ok {
					gotToolCall, ok := got.Content[i].(*api.ToolCallBlock)
					require.True(t, ok, "expected content block %d to be a ToolCallBlock", i)

					assert.Equal(t, wantToolCall.ToolCallID, gotToolCall.ToolCallID)
					assert.Equal(t, wantToolCall.ToolName, gotToolCall.ToolName)
					assert.JSONEq(t, string(wantToolCall.Args), string(gotToolCall.Args))
					assert.Equal(t, wantToolCall.ProviderMetadata, gotToolCall.ProviderMetadata)
				} else {
					assert.Equal(t, wantBlock, got.Content[i])
				}
			}
		})
	}
}

func TestDecodeUsage(t *testing.T) {
	tests := []struct {
		name     string
		input    responses.ResponseUsage
		expected api.Usage
	}{
		{
			name: "basic usage",
			input: responses.ResponseUsage{
				InputTokens:  100,
				OutputTokens: 50,
			},
			expected: api.Usage{
				InputTokens:  100,
				OutputTokens: 50,
				TotalTokens:  150,
			},
		},
		{
			name: "with total tokens",
			input: responses.ResponseUsage{
				InputTokens:  100,
				OutputTokens: 50,
				TotalTokens:  200,
			},
			expected: api.Usage{
				InputTokens:  100,
				OutputTokens: 50,
				TotalTokens:  200,
			},
		},
		{
			name: "total tokens differs from sum",
			input: responses.ResponseUsage{
				InputTokens:  100,
				OutputTokens: 50,
				TotalTokens:  175, // Different from sum (150)
			},
			expected: api.Usage{
				InputTokens:  100,
				OutputTokens: 50,
				TotalTokens:  175, // Should use provided total, not sum
			},
		},
		{
			name: "with token details",
			input: responses.ResponseUsage{
				InputTokens:  100,
				OutputTokens: 50,
				InputTokensDetails: responses.ResponseUsageInputTokensDetails{
					CachedTokens: 25,
				},
				OutputTokensDetails: responses.ResponseUsageOutputTokensDetails{
					ReasoningTokens: 30,
				},
			},
			expected: api.Usage{
				InputTokens:       100,
				OutputTokens:      50,
				TotalTokens:       150,
				ReasoningTokens:   30,
				CachedInputTokens: 25,
			},
		},
		{
			name: "zero values",
			input: responses.ResponseUsage{
				InputTokens:  0,
				OutputTokens: 0,
			},
			expected: api.Usage{
				InputTokens:  0,
				OutputTokens: 0,
				TotalTokens:  0,
			},
		},
		{
			name: "large values",
			input: responses.ResponseUsage{
				InputTokens:  1000000,
				OutputTokens: 500000,
			},
			expected: api.Usage{
				InputTokens:  1000000,
				OutputTokens: 500000,
				TotalTokens:  1500000,
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := decodeUsage(testCase.input)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestDecodeFinishReason(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     api.FinishReason
		wantErr  string
	}{
		{
			name: "max_output_tokens reason",
			response: `{
				"incomplete_details": {
					"reason": "max_output_tokens"
				}
			}`,
			want: api.FinishReasonLength,
		},
		{
			name: "content_filter reason",
			response: `{
				"incomplete_details": {
					"reason": "content_filter"
				}
			}`,
			want: api.FinishReasonContentFilter,
		},
		{
			name: "empty reason with tool calls",
			response: `{
				"incomplete_details": {
					"reason": ""
				},
				"output": [
					{
						"type": "function_call"
					}
				]
			}`,
			want:    api.FinishReasonToolCalls,
			wantErr: "function call missing name",
		},
		{
			name: "empty reason without tool calls",
			response: `{
				"incomplete_details": {
					"reason": ""
				},
				"output": [
					{
						"type": "message"
					}
				]
			}`,
			want: api.FinishReasonStop,
		},
		{
			name: "unknown reason with tool calls",
			response: `{
				"incomplete_details": {
					"reason": "unknown_reason"
				},
				"output": [
					{
						"type": "function_call"
					}
				]
			}`,
			want:    api.FinishReasonToolCalls,
			wantErr: "function call missing name",
		},
		{
			name: "unknown reason without tool calls",
			response: `{
				"incomplete_details": {
					"reason": "unknown_reason"
				},
				"output": [
					{
						"type": "message"
					}
				]
			}`,
			want: api.FinishReasonUnknown,
		},
		{
			name:     "nil response",
			response: `null`,
			want:     api.FinishReasonStop,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var msg *responses.Response
			if testCase.response != "null" {
				msg = &responses.Response{}
				err := json.Unmarshal([]byte(testCase.response), msg)
				require.NoError(t, err)
			}

			// Get hasToolCalls from decodeContent
			hasToolCalls := false
			if msg != nil {
				got, err := decodeContent(msg)
				if testCase.wantErr != "" {
					require.Error(t, err)
					assert.Contains(t, err.Error(), testCase.wantErr)
					return
				}
				require.NoError(t, err)
				hasToolCalls = got.HasTools
			}

			// Get the incomplete reason from the response
			var incompleteReason string
			if msg != nil && msg.IncompleteDetails.Reason != "" {
				incompleteReason = msg.IncompleteDetails.Reason
			}

			got := decodeFinishReason(incompleteReason, hasToolCalls)
			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestHasToolCalls(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     bool
		wantErr  string
	}{
		{
			name: "has tool calls",
			response: `{
				"output": [
					{
						"type": "function_call"
					}
				]
			}`,
			want:    true,
			wantErr: "function call missing name",
		},
		{
			name: "no tool calls",
			response: `{
				"output": [
					{
						"type": "message"
					}
				]
			}`,
			want: false,
		},
		{
			name: "mixed output",
			response: `{
				"output": [
					{
						"type": "message"
					},
					{
						"type": "function_call"
					}
				]
			}`,
			want:    true,
			wantErr: "function call missing name",
		},
		{
			name: "empty output",
			response: `{
				"output": []
			}`,
			want: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var msg responses.Response
			err := json.Unmarshal([]byte(testCase.response), &msg)
			require.NoError(t, err)

			got, err := decodeContent(&msg)
			if testCase.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, got.HasTools)
		})
	}
}

func TestDecodeProviderMetadata(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *Metadata
	}{
		{
			name: "complete metadata",
			input: `{
				"id": "resp_123",
				"usage": {
					"input_tokens": 100,
					"output_tokens": 50,
					"input_tokens_details": {
						"cached_tokens": 25
					},
					"output_tokens_details": {
						"reasoning_tokens": 30
					}
				}
			}`,
			want: &Metadata{
				ResponseID: "resp_123",
				Usage: Usage{
					InputTokens:           100,
					OutputTokens:          50,
					InputCachedTokens:     25,
					OutputReasoningTokens: 30,
				},
			},
		},
		{
			name: "minimal metadata",
			input: `{
				"id": "resp_456",
				"usage": {
					"input_tokens": 10,
					"output_tokens": 5
				}
			}`,
			want: &Metadata{
				ResponseID: "resp_456",
				Usage: Usage{
					InputTokens:  10,
					OutputTokens: 5,
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var msg responses.Response
			err := json.Unmarshal([]byte(testCase.input), &msg)
			require.NoError(t, err)

			got := decodeProviderMetadata(&msg)

			gotMetaAny, ok := got.Get("openai")
			require.True(t, ok, "got metadata should exist")
			gotMeta, ok := gotMetaAny.(*Metadata)
			require.True(t, ok, "got metadata should be of type *Metadata")

			assert.Equal(t, testCase.want, gotMeta, "metadata should match")
		})
	}
}

func TestDecodeReasoning(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *api.ReasoningBlock
		wantErr string
	}{
		{
			name: "valid reasoning with single summary",
			input: `{
				"type": "reasoning",
				"id": "reason_123",
				"summary": [
					{"text": "This is my reasoning", "type": "summary_text"}
				]
			}`,
			want: &api.ReasoningBlock{
				Text: "This is my reasoning",
			},
		},
		{
			name: "valid reasoning with multiple summaries",
			input: `{
				"type": "reasoning",
				"id": "reason_123",
				"summary": [
					{"text": "First point", "type": "summary_text"},
					{"text": "Second point", "type": "summary_text"}
				]
			}`,
			want: &api.ReasoningBlock{
				Text: "First point\nSecond point",
			},
		},
		{
			name: "empty summary array",
			input: `{
				"type": "reasoning",
				"id": "reason_123",
				"summary": []
			}`,
			wantErr: "reasoning item has no summary",
		},
		{
			name: "empty text in summary",
			input: `{
				"type": "reasoning",
				"id": "reason_123",
				"summary": [
					{"text": "", "type": "summary_text"}
				]
			}`,
			wantErr: "empty text in reasoning summary at index 0",
		},
		{
			name: "wrong type",
			input: `{
				"type": "not_reasoning"
			}`,
			wantErr: "unexpected item type for reasoning: not_reasoning",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var item responses.ResponseOutputItemUnion
			err := json.Unmarshal([]byte(testCase.input), &item)
			require.NoError(t, err)

			got, err := decodeReasoning(item)
			if testCase.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, got)
		})
	}
}
