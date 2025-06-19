package codec

import (
	"encoding/json"
	"testing"

	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

func TestEncodeTools(t *testing.T) {
	tests := []struct {
		name               string
		tools              []api.ToolDefinition
		toolChoice         *api.ToolChoice
		expectedTools      []string
		expectedToolChoice string
		expectedWarnings   []api.CallWarning
		expectedError      string
	}{
		{
			name:  "empty tools list",
			tools: []api.ToolDefinition{},
		},
		{
			name: "single function tool",
			tools: []api.ToolDefinition{
				api.FunctionTool{
					Name: "test_function",
					InputSchema: &jsonschema.Definition{
						Type: "object",
						Properties: map[string]jsonschema.Definition{
							"param1": {Type: "string"},
						},
					},
				},
			},
			expectedTools: []string{
				`{
					"type": "function",
					"name": "test_function",
					"parameters": {
						"type": "object",
						"properties": {
							"param1": {
								"type": "string"
							}
						}
					},
					"strict": true
				}`,
			},
		},
		{
			name: "multiple function tools",
			tools: []api.ToolDefinition{
				api.FunctionTool{
					Name: "function1",
					InputSchema: &jsonschema.Definition{
						Type: "object",
						Properties: map[string]jsonschema.Definition{
							"param1": {Type: "string"},
						},
					},
				},
				api.FunctionTool{
					Name: "function2",
					InputSchema: &jsonschema.Definition{
						Type: "object",
						Properties: map[string]jsonschema.Definition{
							"param2": {Type: "number"},
						},
						Required: []string{"param2"},
					},
				},
			},
			expectedTools: []string{
				`{
					"type": "function",
					"name": "function1",
					"parameters": {
						"type": "object",
						"properties": {
							"param1": {
								"type": "string"
							}
						}
					},
					"strict": true
				}`,
				`{
					"type": "function",
					"name": "function2",
					"parameters": {
						"type": "object",
						"properties": {
							"param2": {
								"type": "number"
							}
						},
						"required": ["param2"]
					},
					"strict": true
				}`,
			},
		},
		{
			name: "file search tool (constructor)",
			tools: []api.ToolDefinition{
				FileSearchTool(WithVectorStoreIDs("store1", "store2")),
			},
			expectedTools: []string{
				`{
					"type": "file_search",
					"vector_store_ids": ["store1", "store2"]
				}`,
			},
		},
		{
			name: "file search tool (map args)",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "openai.file_search",
					Name: "file_search",
					Args: map[string]any{
						"vector_store_ids": []string{"store1", "store2"},
					},
				},
			},
			expectedTools: []string{
				`{
					"type": "file_search",
					"vector_store_ids": ["store1", "store2"]
				}`,
			},
		},
		{
			name: "mix of function and provider-defined tools",
			tools: []api.ToolDefinition{
				api.FunctionTool{
					Name: "mixed_function",
					InputSchema: &jsonschema.Definition{
						Type: "object",
						Properties: map[string]jsonschema.Definition{
							"param1": {Type: "string"},
						},
					},
				},
				FileSearchTool(WithVectorStoreIDs("store3")),
			},
			expectedTools: []string{
				`{
					"type": "function",
					"name": "mixed_function",
					"parameters": {
						"type": "object",
						"properties": {
							"param1": {
								"type": "string"
							}
						}
					},
					"strict": true
				}`,
				`{
					"type": "file_search",
					"vector_store_ids": ["store3"]
				}`,
			},
		},
		{
			name: "unsupported tool type with warning",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "unsupported_tool",
					Name: "unsupported",
					Args: &mockUnsupportedTool{id: "unsupported_tool"},
				},
			},
			expectedWarnings: []api.CallWarning{
				{
					Type: "unsupported-tool",
					Tool: api.ProviderDefinedTool{
						ID:   "unsupported_tool",
						Name: "unsupported",
						Args: &mockUnsupportedTool{id: "unsupported_tool"},
					},
				},
			},
		},
		{
			name: "tool choice auto",
			tools: []api.ToolDefinition{
				api.FunctionTool{
					Name: "function_with_choice",
					InputSchema: &jsonschema.Definition{
						Type: "object",
						Properties: map[string]jsonschema.Definition{
							"param1": {Type: "string"},
						},
					},
				},
			},
			toolChoice: &api.ToolChoice{
				Type: "auto",
			},
			expectedTools: []string{
				`{
					"type": "function",
					"name": "function_with_choice",
					"parameters": {
						"type": "object",
						"properties": {
							"param1": {
								"type": "string"
							}
						}
					},
					"strict": true
				}`,
			},
			expectedToolChoice: `"auto"`,
		},
		{
			name: "tool choice none",
			tools: []api.ToolDefinition{
				api.FunctionTool{
					Name: "function_with_choice",
					InputSchema: &jsonschema.Definition{
						Type: "object",
						Properties: map[string]jsonschema.Definition{
							"param1": {Type: "string"},
						},
					},
				},
			},
			toolChoice: &api.ToolChoice{
				Type: "none",
			},
			expectedTools: []string{
				`{
					"type": "function",
					"name": "function_with_choice",
					"parameters": {
						"type": "object",
						"properties": {
							"param1": {
								"type": "string"
							}
						}
					},
					"strict": true
				}`,
			},
			expectedToolChoice: `"none"`,
		},
		{
			name: "tool choice specific function",
			tools: []api.ToolDefinition{
				api.FunctionTool{
					Name: "function1",
					InputSchema: &jsonschema.Definition{
						Type: "object",
						Properties: map[string]jsonschema.Definition{
							"param1": {Type: "string"},
						},
					},
				},
				api.FunctionTool{
					Name: "function2",
					InputSchema: &jsonschema.Definition{
						Type: "object",
						Properties: map[string]jsonschema.Definition{
							"param2": {Type: "number"},
						},
					},
				},
			},
			toolChoice: &api.ToolChoice{
				Type:     "tool",
				ToolName: "function2",
			},
			expectedTools: []string{
				`{
					"type": "function",
					"name": "function1",
					"parameters": {
						"type": "object",
						"properties": {
							"param1": {
								"type": "string"
							}
						}
					},
					"strict": true
				}`,
				`{
					"type": "function",
					"name": "function2",
					"parameters": {
						"type": "object",
						"properties": {
							"param2": {
								"type": "number"
							}
						}
					},
					"strict": true
				}`,
			},
			expectedToolChoice: `{"type":"function","name":"function2"}`,
		},
		{
			name: "tool choice provider-defined tool (constructor)",
			tools: []api.ToolDefinition{
				FileSearchTool(WithVectorStoreIDs("store1")),
			},
			toolChoice: &api.ToolChoice{
				Type:     "tool",
				ToolName: "file_search",
			},
			expectedTools: []string{
				`{
					"type": "file_search",
					"vector_store_ids": ["store1"]
				}`,
			},
			expectedToolChoice: `{"type":"file_search"}`,
		},
		{
			name: "tool choice provider-defined tool (map args)",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "openai.file_search",
					Name: "file_search",
					Args: map[string]any{
						"vector_store_ids": []string{"store1"},
					},
				},
			},
			toolChoice: &api.ToolChoice{
				Type:     "tool",
				ToolName: "file_search",
			},
			expectedTools: []string{
				`{
					"type": "file_search",
					"vector_store_ids": ["store1"]
				}`,
			},
			expectedToolChoice: `{"type":"file_search"}`,
		},
		{
			name: "web search tool with minimal settings (constructor)",
			tools: []api.ToolDefinition{
				WebSearchTool(),
			},
			expectedTools: []string{
				`{
					"type": "web_search_preview"
				}`,
			},
		},
		{
			name: "web search tool with minimal settings (map args)",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "openai.web_search_preview",
					Name: "web_search_preview",
					Args: map[string]any{},
				},
			},
			expectedTools: []string{
				`{
					"type": "web_search_preview"
				}`,
			},
		},
		{
			name: "web search tool with search context size (constructor)",
			tools: []api.ToolDefinition{
				WebSearchTool(WithSearchContextSize("large")),
			},
			expectedTools: []string{
				`{
					"type": "web_search_preview",
					"search_context_size": "large"
				}`,
			},
		},
		{
			name: "web search tool with search context size (map args)",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "openai.web_search_preview",
					Name: "web_search_preview",
					Args: map[string]any{
						"search_context_size": "large",
					},
				},
			},
			expectedTools: []string{
				`{
					"type": "web_search_preview",
					"search_context_size": "large"
				}`,
			},
		},
		{
			name: "web search tool with user location (constructor)",
			tools: []api.ToolDefinition{
				WebSearchTool(WithUserLocation(&WebSearchUserLocation{
					City:     "San Francisco",
					Country:  "US",
					Region:   "CA",
					Timezone: "America/Los_Angeles",
				})),
			},
			expectedTools: []string{
				`{
					"type": "web_search_preview",
					"user_location": {
						"city": "San Francisco",
						"country": "US",
						"region": "CA",
						"timezone": "America/Los_Angeles",
						"type": "approximate"
					}
				}`,
			},
		},
		{
			name: "web search tool with user location (map args)",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "openai.web_search_preview",
					Name: "web_search_preview",
					Args: map[string]any{
						"user_location": map[string]any{
							"city":     "San Francisco",
							"country":  "US",
							"region":   "CA",
							"timezone": "America/Los_Angeles",
						},
					},
				},
			},
			expectedTools: []string{
				`{
					"type": "web_search_preview",
					"user_location": {
						"city": "San Francisco",
						"country": "US",
						"region": "CA",
						"timezone": "America/Los_Angeles",
						"type": "approximate"
					}
				}`,
			},
		},
		{
			name: "computer use tool (constructor)",
			tools: []api.ToolDefinition{
				ComputerUseTool(1366, 768, "windows"),
			},
			expectedTools: []string{
				`{
					"type": "computer_use_preview",
					"display_height": 768,
					"display_width": 1366,
					"environment": "windows"
				}`,
			},
		},
		{
			name: "computer use tool (map args)",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "openai.computer_use_preview",
					Name: "computer_use_preview",
					Args: map[string]any{
						"display_height": 768,
						"display_width":  1366,
						"environment":    "windows",
					},
				},
			},
			expectedTools: []string{
				`{
					"type": "computer_use_preview",
					"display_height": 768,
					"display_width": 1366,
					"environment": "windows"
				}`,
			},
		},
		{
			name: "computer use tool missing required display width",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "openai.computer_use_preview",
					Name: "computer_use_preview",
					Args: map[string]any{
						"display_height": 768,
						"environment":    "windows",
					},
				},
			},
			expectedError: "displayWidth is required and must be positive",
		},
		{
			name: "computer use tool missing required display height",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "openai.computer_use_preview",
					Name: "computer_use_preview",
					Args: map[string]any{
						"display_width": 1366,
						"environment":   "windows",
					},
				},
			},
			expectedError: "displayHeight is required and must be positive",
		},
		{
			name: "web search tool with partial user location (constructor)",
			tools: []api.ToolDefinition{
				WebSearchTool(
					WithSearchContextSize("medium"),
					WithUserLocation(&WebSearchUserLocation{
						City:    "London",
						Country: "UK",
					}),
				),
			},
			expectedTools: []string{
				`{
					"type": "web_search_preview",
					"search_context_size": "medium",
					"user_location": {
						"city": "London",
						"country": "UK",
						"type": "approximate"
					}
				}`,
			},
		},
		{
			name: "web search tool with partial user location (map args)",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "openai.web_search_preview",
					Name: "web_search_preview",
					Args: map[string]any{
						"search_context_size": "medium",
						"user_location": map[string]any{
							"city":    "London",
							"country": "UK",
						},
					},
				},
			},
			expectedTools: []string{
				`{
					"type": "web_search_preview",
					"search_context_size": "medium",
					"user_location": {
						"city": "London",
						"country": "UK",
						"type": "approximate"
					}
				}`,
			},
		},
		{
			name: "computer use tool with mac environment (constructor)",
			tools: []api.ToolDefinition{
				ComputerUseTool(1200, 800, "mac"),
			},
			expectedTools: []string{
				`{
					"type": "computer_use_preview",
					"display_height": 800,
					"display_width": 1200,
					"environment": "mac"
				}`,
			},
		},
		{
			name: "computer use tool with mac environment (map args)",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "openai.computer_use_preview",
					Name: "computer_use_preview",
					Args: map[string]any{
						"display_height": 800,
						"display_width":  1200,
						"environment":    "mac",
					},
				},
			},
			expectedTools: []string{
				`{
					"type": "computer_use_preview",
					"display_height": 800,
					"display_width": 1200,
					"environment": "mac"
				}`,
			},
		},
		{
			name: "computer use tool with browser environment (constructor)",
			tools: []api.ToolDefinition{
				ComputerUseTool(1920, 1080, "browser"),
			},
			expectedTools: []string{
				`{
					"type": "computer_use_preview",
					"display_height": 1080,
					"display_width": 1920,
					"environment": "browser"
				}`,
			},
		},
		{
			name: "computer use tool with browser environment (map args)",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "openai.computer_use_preview",
					Name: "computer_use_preview",
					Args: map[string]any{
						"display_height": 1080,
						"display_width":  1920,
						"environment":    "browser",
					},
				},
			},
			expectedTools: []string{
				`{
					"type": "computer_use_preview",
					"display_height": 1080,
					"display_width": 1920,
					"environment": "browser"
				}`,
			},
		},
		{
			name: "computer use tool with invalid environment",
			tools: []api.ToolDefinition{
				api.ProviderDefinedTool{
					ID:   "openai.computer_use_preview",
					Name: "computer_use_preview",
					Args: map[string]any{
						"display_height": 768,
						"display_width":  1366,
						"environment":    "invalid_env",
					},
				},
			},
			expectedError: "environment must be one of: mac, windows, ubuntu, browser",
		},
		{
			name: "invalid tool choice type",
			tools: []api.ToolDefinition{
				api.FunctionTool{
					Name: "function1",
					InputSchema: &jsonschema.Definition{
						Type: "object",
						Properties: map[string]jsonschema.Definition{
							"param1": {Type: "string"},
						},
					},
				},
			},
			toolChoice: &api.ToolChoice{
				Type: "invalid",
			},
			expectedError: "unsupported tool choice type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := EncodeTools(tc.tools, tc.toolChoice, api.CallOptions{})

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				return
			}

			require.NoError(t, err)

			// Check if the warnings match the expected warnings if specified
			if len(tc.expectedWarnings) > 0 {
				assert.ElementsMatch(t, tc.expectedWarnings, result.Warnings, "Warning mismatch")
			}

			assert.Equal(t, len(tc.expectedTools), len(result.Tools), "Tools count mismatch")
			for i, expectedTool := range tc.expectedTools {
				toolBytes, err := json.Marshal(result.Tools[i])
				require.NoError(t, err)
				assert.JSONEq(t, expectedTool, string(toolBytes))
			}

			// Check tool choice if expected
			if tc.expectedToolChoice != "" {
				toolChoiceBytes, err := json.Marshal(result.ToolChoice)
				require.NoError(t, err)
				assert.JSONEq(t, tc.expectedToolChoice, string(toolChoiceBytes))
			}
		})
	}
}

func TestEncodeToolChoice(t *testing.T) {
	tests := []struct {
		name          string
		input         *api.ToolChoice
		expected      string
		expectedError string
	}{
		{
			name:     "nil tool choice",
			input:    nil,
			expected: `null`,
		},
		{
			name: "auto tool choice",
			input: &api.ToolChoice{
				Type: "auto",
			},
			expected: `"auto"`,
		},
		{
			name: "none tool choice",
			input: &api.ToolChoice{
				Type: "none",
			},
			expected: `"none"`,
		},
		{
			name: "function tool choice",
			input: &api.ToolChoice{
				Type:     "tool",
				ToolName: "my_function",
			},
			expected: `{
				"type": "function",
				"name": "my_function"
			}`,
		},
		{
			name: "provider-defined tool choice",
			input: &api.ToolChoice{
				Type:     "tool",
				ToolName: "web_search_preview",
			},
			expected: `{
				"type": "web_search_preview"
			}`,
		},
		{
			name: "invalid tool choice type",
			input: &api.ToolChoice{
				Type: "invalid",
			},
			expectedError: "unsupported tool choice type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := encodeToolChoice(tc.input)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				return
			}

			require.NoError(t, err)

			resultJSON, err := json.Marshal(result)
			require.NoError(t, err)

			assert.JSONEq(t, tc.expected, string(resultJSON))
		})
	}
}

func TestJsonSchemaAsMap(t *testing.T) {
	tests := []struct {
		name          string
		input         *jsonschema.Definition
		expected      string
		expectedError string
	}{
		{
			name:     "nil schema",
			input:    nil,
			expected: `null`,
		},
		{
			name: "simple schema",
			input: &jsonschema.Definition{
				Type: "object",
				Properties: map[string]jsonschema.Definition{
					"name": {Type: "string"},
					"age":  {Type: "number"},
				},
				Required: []string{"name"},
			},
			expected: `{
				"type": "object",
				"properties": {
					"name": {"type": "string"},
					"age": {"type": "number"}
				},
				"required": ["name"]
			}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := jsonSchemaAsMap(tc.input)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				return
			}

			require.NoError(t, err)

			if result == nil {
				assert.Equal(t, "null", tc.expected)
			} else {
				resultJSON, err := json.Marshal(result)
				require.NoError(t, err)

				assert.JSONEq(t, tc.expected, string(resultJSON))
			}
		})
	}
}

func TestEncodeProviderDefinedTool(t *testing.T) {
	tests := []struct {
		name             string
		input            api.ProviderDefinedTool
		expected         string
		expectedWarnings []api.CallWarning
		expectedError    string
	}{
		{
			name:  "file search tool (constructor)",
			input: FileSearchTool(WithVectorStoreIDs("store1", "store2")),
			expected: `{
				"type": "file_search",
				"vector_store_ids": ["store1", "store2"]
			}`,
			expectedWarnings: nil,
		},
		{
			name: "file search tool (map args)",
			input: api.ProviderDefinedTool{
				ID:   "openai.file_search",
				Name: "file_search",
				Args: map[string]any{
					"vector_store_ids": []string{"store1", "store2"},
				},
			},
			expected: `{
				"type": "file_search",
				"vector_store_ids": ["store1", "store2"]
			}`,
			expectedWarnings: nil,
		},
		{
			name: "unsupported provider tool",
			input: api.ProviderDefinedTool{
				ID:   "unsupported_tool",
				Name: "unsupported",
				Args: &mockUnsupportedTool{id: "unsupported_tool"},
			},
			expected: `null`,
			expectedWarnings: []api.CallWarning{
				{
					Type: "unsupported-tool",
					Tool: api.ProviderDefinedTool{
						ID:   "unsupported_tool",
						Name: "unsupported",
						Args: &mockUnsupportedTool{id: "unsupported_tool"},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, warnings, err := encodeProviderDefinedTool(tc.input)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				return
			}

			require.NoError(t, err)

			assert.ElementsMatch(t, tc.expectedWarnings, warnings, "Warning mismatch")

			if tc.expected == `null` {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				resultJSON, err := json.Marshal(result)
				require.NoError(t, err)

				assert.JSONEq(t, tc.expected, string(resultJSON))
			}
		})
	}
}

// Mock unsupported tool for testing
type mockUnsupportedTool struct {
	id string
}
