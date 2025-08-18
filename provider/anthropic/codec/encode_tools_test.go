package codec

import (
	"encoding/json"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

func TestEncodeFunctionTool(t *testing.T) {
	tests := []struct {
		name    string
		input   api.FunctionTool
		wantErr bool
		want    string // Expected JSON output
	}{
		{
			name: "simple function tool",
			input: api.FunctionTool{
				Name:        "test_function",
				Description: "A test function",
				InputSchema: &jsonschema.Schema{
					Type: "object",
					Properties: map[string]*jsonschema.Schema{
						"param1": {
							Type:        "string",
							Description: "First parameter",
						},
					},
					Required: []string{"param1"},
				},
			},
			want: `{
				"type": "tool",
				"name": "test_function",
				"description": "A test function",
				"input_schema": {
					"type": "object",
					"properties": {
						"param1": {
							"type": "string",
							"description": "First parameter"
						}
					},
					"required": ["param1"]
				}
			}`,
		},
		{
			name: "function tool with additionalProperties false",
			input: api.FunctionTool{
				Name:        "test_function",
				Description: "A test function",
				InputSchema: &jsonschema.Schema{
					Type: "object",
					Properties: map[string]*jsonschema.Schema{
						"param1": {Type: "string"},
					},
					AdditionalProperties: api.FalseSchema(),
				},
			},
			want: `{
				"type": "tool",
				"name": "test_function",
				"description": "A test function",
				"input_schema": {
					"type": "object",
					"properties": {
						"param1": {
							"type": "string"
						}
					},
					"additionalProperties": false
				}
			}`,
		},
		{
			name: "function tool with minimal fields",
			input: api.FunctionTool{
				Name: "minimal_function",
				InputSchema: &jsonschema.Schema{
					Type: "object",
				},
			},
			want: `{
				"type": "tool",
				"name": "minimal_function",
				"description": "",
				"input_schema": {
					"type": "object"
				}
			}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := EncodeFunctionTool(&tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			resultJSON, err := json.Marshal(result)
			require.NoError(t, err)
			assert.JSONEq(t, tc.want, string(resultJSON))
		})
	}
}

func TestEncodeToolChoice(t *testing.T) {
	tests := []struct {
		name      string
		input     *api.ToolChoice
		expectNil bool
		wantErr   bool
		wantJSON  string // Expected JSON output, only used when not nil
	}{
		{
			name:      "nil input",
			input:     nil,
			expectNil: true,
		},
		{
			name: "auto choice",
			input: &api.ToolChoice{
				Type: "auto",
			},
			wantJSON: `{"type": "auto"}`,
		},
		{
			name: "required choice",
			input: &api.ToolChoice{
				Type: "required",
			},
			wantJSON: `{"type": "any"}`,
		},
		{
			name: "none choice",
			input: &api.ToolChoice{
				Type: "none",
			},
			expectNil: true,
		},
		{
			name: "specific tool choice",
			input: &api.ToolChoice{
				Type:     "tool",
				ToolName: "test_tool",
			},
			wantJSON: `{"type": "tool", "name": "test_tool"}`,
		},
		{
			name: "unsupported choice type",
			input: &api.ToolChoice{
				Type: "invalid_type",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := EncodeToolChoice(tc.input)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tc.expectNil {
				assert.Zero(t, result)
				return
			}

			assert.NotZero(t, result)
			resultJSON, err := json.Marshal(result)
			require.NoError(t, err)
			assert.JSONEq(t, tc.wantJSON, string(resultJSON))
		})
	}
}

func TestEncodeProviderDefinedTool(t *testing.T) {
	tests := []struct {
		name         string
		input        *api.ProviderDefinedTool
		expectNil    bool
		expectBetas  []string
		wantWarnings []api.CallWarning // Expected warnings (empty means no warnings)
		wantErrMsg   string            // Empty means no error, non-empty means expect error containing this string
		want         anthropic.BetaToolUnionParam
	}{
		// Computer tool tests using constructor functions
		{
			name:         "computer tool with version 20250124 (constructor)",
			input:        ComputerTool(800, 600, WithComputerVersion("20250124"), WithDisplayNumber(1)),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfComputerUseTool20250124: &anthropic.BetaToolComputerUse20250124Param{
					DisplayWidthPx:  800,
					DisplayHeightPx: 600,
					DisplayNumber:   anthropic.Int(1),
				},
			},
		},
		{
			name:         "computer tool with version 20241022 (constructor)",
			input:        ComputerTool(800, 600, WithComputerVersion("20241022"), WithDisplayNumber(1)),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfComputerUseTool20241022: &anthropic.BetaToolComputerUse20241022Param{
					DisplayWidthPx:  800,
					DisplayHeightPx: 600,
					DisplayNumber:   anthropic.Int(1),
				},
			},
		},
		{
			name:         "computer tool with default version (constructor)",
			input:        ComputerTool(800, 600, WithDisplayNumber(1)),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfComputerUseTool20250124: &anthropic.BetaToolComputerUse20250124Param{
					DisplayWidthPx:  800,
					DisplayHeightPx: 600,
					DisplayNumber:   anthropic.Int(1),
				},
			},
		},

		// Computer tool tests using map[string]any args
		{
			name: "computer tool with version 20250124 (map args)",
			input: &api.ProviderDefinedTool{
				ID:   "anthropic.computer",
				Name: "computer",
				Args: map[string]any{
					"version":           "20250124",
					"display_width_px":  800,
					"display_height_px": 600,
					"display_number":    1,
				},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfComputerUseTool20250124: &anthropic.BetaToolComputerUse20250124Param{
					DisplayWidthPx:  800,
					DisplayHeightPx: 600,
					DisplayNumber:   anthropic.Int(1),
				},
			},
		},
		{
			name: "computer tool with version 20241022 (map args)",
			input: &api.ProviderDefinedTool{
				ID:   "anthropic.computer",
				Name: "computer",
				Args: map[string]any{
					"version":           "20241022",
					"display_width_px":  800,
					"display_height_px": 600,
					"display_number":    1,
				},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfComputerUseTool20241022: &anthropic.BetaToolComputerUse20241022Param{
					DisplayWidthPx:  800,
					DisplayHeightPx: 600,
					DisplayNumber:   anthropic.Int(1),
				},
			},
		},
		{
			name: "computer tool with default version (map args)",
			input: &api.ProviderDefinedTool{
				ID:   "anthropic.computer",
				Name: "computer",
				Args: map[string]any{
					"display_width_px":  800,
					"display_height_px": 600,
					"display_number":    1,
				},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfComputerUseTool20250124: &anthropic.BetaToolComputerUse20250124Param{
					DisplayWidthPx:  800,
					DisplayHeightPx: 600,
					DisplayNumber:   anthropic.Int(1),
				},
			},
		},

		// Text editor tool tests using constructor functions
		{
			name:         "text editor tool with version 20250124 (constructor)",
			input:        TextEditorTool(WithTextEditorVersion("20250124")),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfTextEditor20250124: &anthropic.BetaToolTextEditor20250124Param{},
			},
		},
		{
			name:         "text editor tool with version 20241022 (constructor)",
			input:        TextEditorTool(WithTextEditorVersion("20241022")),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfTextEditor20241022: &anthropic.BetaToolTextEditor20241022Param{},
			},
		},
		{
			name:         "text editor tool with default version (constructor)",
			input:        TextEditorTool(),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfTextEditor20250124: &anthropic.BetaToolTextEditor20250124Param{},
			},
		},

		// Text editor tool tests using map[string]any args
		{
			name: "text editor tool with version 20250124 (map args)",
			input: &api.ProviderDefinedTool{
				ID:   "anthropic.text_editor",
				Name: "str_replace_editor",
				Args: map[string]any{
					"version": "20250124",
				},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfTextEditor20250124: &anthropic.BetaToolTextEditor20250124Param{},
			},
		},
		{
			name: "text editor tool with version 20241022 (map args)",
			input: &api.ProviderDefinedTool{
				ID:   "anthropic.text_editor",
				Name: "str_replace_editor",
				Args: map[string]any{
					"version": "20241022",
				},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfTextEditor20241022: &anthropic.BetaToolTextEditor20241022Param{},
			},
		},
		{
			name: "text editor tool with default version (map args)",
			input: &api.ProviderDefinedTool{
				ID:   "anthropic.text_editor",
				Name: "str_replace_editor",
				Args: map[string]any{},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfTextEditor20250124: &anthropic.BetaToolTextEditor20250124Param{},
			},
		},

		// Bash tool tests using constructor functions
		{
			name:         "bash tool with version 20250124 (constructor)",
			input:        BashTool(WithBashVersion("20250124")),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfBashTool20250124: &anthropic.BetaToolBash20250124Param{},
			},
		},
		{
			name:         "bash tool with version 20241022 (constructor)",
			input:        BashTool(WithBashVersion("20241022")),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfBashTool20241022: &anthropic.BetaToolBash20241022Param{},
			},
		},
		{
			name:         "bash tool with default version (constructor)",
			input:        BashTool(),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfBashTool20250124: &anthropic.BetaToolBash20250124Param{},
			},
		},

		// Bash tool tests using map[string]any args
		{
			name: "bash tool with version 20250124 (map args)",
			input: &api.ProviderDefinedTool{
				ID:   "anthropic.bash",
				Name: "bash",
				Args: map[string]any{
					"version": "20250124",
				},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfBashTool20250124: &anthropic.BetaToolBash20250124Param{},
			},
		},
		{
			name: "bash tool with version 20241022 (map args)",
			input: &api.ProviderDefinedTool{
				ID:   "anthropic.bash",
				Name: "bash",
				Args: map[string]any{
					"version": "20241022",
				},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfBashTool20241022: &anthropic.BetaToolBash20241022Param{},
			},
		},
		{
			name: "bash tool with default version (map args)",
			input: &api.ProviderDefinedTool{
				ID:   "anthropic.bash",
				Name: "bash",
				Args: map[string]any{},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolUnionParam{
				OfBashTool20250124: &anthropic.BetaToolBash20250124Param{},
			},
		},

		// Error cases using map[string]any args
		{
			name: "computer tool with invalid version (map args)",
			input: &api.ProviderDefinedTool{
				ID:   "anthropic.computer",
				Name: "computer",
				Args: map[string]any{
					"version":           "invalid",
					"display_width_px":  800,
					"display_height_px": 600,
				},
			},
			wantErrMsg: "unsupported computer tool version",
		},
		{
			name: "text editor tool with invalid version (map args)",
			input: &api.ProviderDefinedTool{
				ID:   "anthropic.text_editor",
				Name: "str_replace_editor",
				Args: map[string]any{
					"version": "invalid",
				},
			},
			wantErrMsg: "unsupported text editor tool version",
		},
		{
			name: "bash tool with invalid version (map args)",
			input: &api.ProviderDefinedTool{
				ID:   "anthropic.bash",
				Name: "bash",
				Args: map[string]any{
					"version": "invalid",
				},
			},
			wantErrMsg: "unsupported bash tool version",
		},

		// Unsupported tool type
		{
			name: "unsupported tool type",
			input: &api.ProviderDefinedTool{
				ID:   "mock.unsupported",
				Name: "unsupported",
				Args: &unsupportedArgs{},
			},
			expectNil:   true,
			expectBetas: []string{},
			wantWarnings: []api.CallWarning{
				{
					Type: "unsupported-tool",
					Tool: &api.ProviderDefinedTool{
						ID:   "mock.unsupported",
						Name: "unsupported",
						Args: &unsupportedArgs{},
					},
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			tool, betas, warnings, err := EncodeProviderDefinedTool(testCase.input)

			if testCase.wantErrMsg != "" {
				assert.Error(t, err, "Expected an error")
				assert.Contains(t, err.Error(), testCase.wantErrMsg, "Error message should contain expected substring")
				return
			}

			require.NoError(t, err)

			// Check warnings
			if len(testCase.wantWarnings) == 0 {
				assert.Empty(t, warnings, "No warnings should be returned")
			} else {
				assert.ElementsMatch(t, testCase.wantWarnings, warnings, "Warnings mismatch")
			}

			// Check betas
			assert.ElementsMatch(t, testCase.expectBetas, betas, "Betas mismatch")

			// Check if tool should be empty (empty union)
			if testCase.expectNil {
				// Check that the tool is empty using GetType()
				assert.Nil(t, tool.GetType(), "Tool should be empty")
				return
			}

			require.NotNil(t, tool, "Tool should not be nil")

			// Validate that GetType() returns a non-empty value (implementation sets it correctly)
			require.NotNil(t, tool.GetType(), "Tool type should not be nil")
			require.NotEmpty(t, *tool.GetType(), "Tool type should not be empty")

			// Validate the JSON representation matches
			// This will verify our type matches the SDK's default when marshaled
			expectedJSON, err := json.Marshal(testCase.want)
			require.NoError(t, err, "Failed to marshal expected tool to JSON")

			actualJSON, err := json.Marshal(tool)
			require.NoError(t, err, "Failed to marshal actual tool to JSON")

			assert.JSONEq(t, string(expectedJSON), string(actualJSON), "Tool JSON content mismatch for %s", testCase.name)
		})
	}
}

func TestEncodeTools(t *testing.T) {
	functionTool := &api.FunctionTool{
		Name:        "test_function",
		Description: "A test function",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"param1": {
					Type:        "string",
					Description: "First parameter",
				},
			},
			Required: []string{"param1"},
		},
	}

	computerTool := &api.ProviderDefinedTool{
		ID:   "anthropic.computer",
		Name: "computer",
		Args: &ComputerToolArgs{
			DisplayWidthPx:  800,
			DisplayHeightPx: 600,
			DisplayNumber:   1,
		},
	}

	// Use a concrete tool type that we know won't be handled correctly
	unsupportedTool := &api.ProviderDefinedTool{
		ID:   "mock.unsupported",
		Name: "unsupported",
		Args: &unsupportedArgs{},
	}

	// Helper to create expected tool choice
	autoChoice := anthropic.BetaToolChoiceUnionParam{
		OfAuto: &anthropic.BetaToolChoiceAutoParam{
			Type: "auto",
		},
	}
	anyChoice := anthropic.BetaToolChoiceUnionParam{
		OfAny: &anthropic.BetaToolChoiceAnyParam{
			Type: "any",
		},
	}
	specificChoice := func(toolName string) anthropic.BetaToolChoiceUnionParam {
		return anthropic.BetaToolChoiceParamOfTool(toolName)
	}

	tests := []struct {
		name    string
		tools   []api.ToolDefinition
		choice  *api.ToolChoice
		want    AnthropicTools
		wantErr bool
	}{
		{
			name:   "no tools",
			tools:  nil,
			choice: nil,
			want:   AnthropicTools{},
		},
		{
			name:  "function tool",
			tools: []api.ToolDefinition{functionTool},
			want: AnthropicTools{
				Tools:      []anthropic.BetaToolUnionParam{mustEncodeFunctionTool(functionTool)},
				ToolChoice: anthropic.BetaToolChoiceUnionParam{},
				Betas:      []anthropic.AnthropicBeta{},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name:   "provider tool with auto choice",
			tools:  []api.ToolDefinition{computerTool},
			choice: &api.ToolChoice{Type: "auto"},
			want: AnthropicTools{
				Tools:      []anthropic.BetaToolUnionParam{mustEncodeComputerTool()},
				ToolChoice: autoChoice,
				Betas:      []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2025_01_24},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name:   "mixed tools with required choice",
			tools:  []api.ToolDefinition{functionTool, computerTool},
			choice: &api.ToolChoice{Type: "required"},
			want: AnthropicTools{
				Tools:      []anthropic.BetaToolUnionParam{mustEncodeFunctionTool(functionTool), mustEncodeComputerTool()},
				ToolChoice: anyChoice,
				Betas:      []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2025_01_24},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name:  "unsupported tool",
			tools: []api.ToolDefinition{unsupportedTool},
			want: AnthropicTools{
				Tools:      []anthropic.BetaToolUnionParam{},
				ToolChoice: anthropic.BetaToolChoiceUnionParam{},
				Betas:      []anthropic.AnthropicBeta{},
				Warnings: []api.CallWarning{
					{Type: "unsupported-tool", Tool: unsupportedTool},
				},
			},
		},
		{
			name:    "unsupported tool choice",
			tools:   []api.ToolDefinition{functionTool},
			choice:  &api.ToolChoice{Type: "invalid_type"},
			wantErr: true,
		},
		{
			name:   "none tool choice",
			tools:  []api.ToolDefinition{functionTool, computerTool},
			choice: &api.ToolChoice{Type: "none"},
			want: AnthropicTools{
				Tools:      nil,
				ToolChoice: anthropic.BetaToolChoiceUnionParam{},
				Betas:      []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2025_01_24},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name:   "specific tool choice",
			tools:  []api.ToolDefinition{functionTool},
			choice: &api.ToolChoice{Type: "tool", ToolName: "test_function"},
			want: AnthropicTools{
				Tools:      []anthropic.BetaToolUnionParam{mustEncodeFunctionTool(functionTool)},
				ToolChoice: specificChoice("test_function"),
				Betas:      []anthropic.AnthropicBeta{},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name: "object tool mode equivalent",
			tools: []api.ToolDefinition{
				&api.FunctionTool{
					Name:        "test_tool",
					Description: "A test tool",
					InputSchema: &jsonschema.Schema{Type: "object"},
				},
			},
			choice: &api.ToolChoice{Type: "tool", ToolName: "test_tool"},
			want: AnthropicTools{
				Tools: []anthropic.BetaToolUnionParam{
					mustEncodeFunctionTool(&api.FunctionTool{
						Name:        "test_tool",
						Description: "A test tool",
						InputSchema: &jsonschema.Schema{Type: "object"},
					}),
				},
				ToolChoice: specificChoice("test_tool"),
				Betas:      []anthropic.AnthropicBeta{},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name: "multiple provider tools with same beta",
			tools: []api.ToolDefinition{
				&api.ProviderDefinedTool{
					ID:   "anthropic.computer",
					Name: "computer",
					Args: &ComputerToolArgs{DisplayWidthPx: 800, DisplayHeightPx: 600, DisplayNumber: 1},
				},
				&api.ProviderDefinedTool{
					ID:   "anthropic.text_editor",
					Name: "str_replace_editor",
					Args: &TextEditorToolArgs{},
				},
				&api.ProviderDefinedTool{
					ID:   "anthropic.bash",
					Name: "bash",
					Args: &BashToolArgs{},
				},
			},
			want: AnthropicTools{
				Tools: []anthropic.BetaToolUnionParam{
					mustEncodeComputerTool(),
					mustEncodeTextEditorTool(),
					mustEncodeBashTool(),
				},
				ToolChoice: anthropic.BetaToolChoiceUnionParam{},
				Betas:      []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2025_01_24},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name:  "mixed supported and unsupported tools",
			tools: []api.ToolDefinition{functionTool, unsupportedTool, computerTool},
			want: AnthropicTools{
				Tools:      []anthropic.BetaToolUnionParam{mustEncodeFunctionTool(functionTool), mustEncodeComputerTool()},
				ToolChoice: anthropic.BetaToolChoiceUnionParam{},
				Betas:      []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2025_01_24},
				Warnings: []api.CallWarning{
					{Type: "unsupported-tool", Tool: unsupportedTool},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := EncodeTools(tc.tools, tc.choice)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Compare each field separately for better error messages
			assert.ElementsMatch(t, tc.want.Tools, result.Tools, "Tools mismatch")
			assert.Equal(t, tc.want.ToolChoice, result.ToolChoice, "ToolChoice mismatch")
			assert.ElementsMatch(t, tc.want.Betas, result.Betas, "Betas mismatch")
			assert.ElementsMatch(t, tc.want.Warnings, result.Warnings, "Warnings mismatch")
		})
	}
}

// Helper functions to create expected tool encodings
func mustEncodeFunctionTool(tool *api.FunctionTool) anthropic.BetaToolUnionParam {
	result, err := EncodeFunctionTool(tool)
	if err != nil {
		panic(err)
	}
	return result
}

func mustEncodeComputerTool() anthropic.BetaToolUnionParam {
	return anthropic.BetaToolUnionParam{
		OfComputerUseTool20250124: &anthropic.BetaToolComputerUse20250124Param{
			Type:            "computer_20250124",
			DisplayWidthPx:  800,
			DisplayHeightPx: 600,
			DisplayNumber:   anthropic.Int(1),
		},
	}
}

func mustEncodeTextEditorTool() anthropic.BetaToolUnionParam {
	return anthropic.BetaToolUnionParam{
		OfTextEditor20250124: &anthropic.BetaToolTextEditor20250124Param{
			Type: "text_editor_20250124",
		},
	}
}

func mustEncodeBashTool() anthropic.BetaToolUnionParam {
	return anthropic.BetaToolUnionParam{
		OfBashTool20250124: &anthropic.BetaToolBash20250124Param{
			Type: "bash_20250124",
		},
	}
}

// unsupportedArgs is used as Args for testing unsupported tools
type unsupportedArgs struct{}
