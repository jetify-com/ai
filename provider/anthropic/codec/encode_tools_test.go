package codec

import (
	"encoding/json"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/sashabaranov/go-openai/jsonschema"
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
				InputSchema: &jsonschema.Definition{
					Type: "object",
					Properties: map[string]jsonschema.Definition{
						"param1": {
							Type:        "string",
							Description: "First parameter",
						},
					},
					Required: []string{"param1"},
				},
			},
			want: `{
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
			name: "function tool with minimal fields",
			input: api.FunctionTool{
				Name: "minimal_function",
				InputSchema: &jsonschema.Definition{
					Type: "object",
				},
			},
			want: `{
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
			result, err := EncodeFunctionTool(tc.input)
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
			wantJSON: `[{"type": "auto"}]`,
		},
		{
			name: "required choice",
			input: &api.ToolChoice{
				Type: "required",
			},
			wantJSON: `[{"type": "any"}]`,
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
			wantJSON: `[{"type": "tool", "name": "test_tool"}]`,
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
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			resultJSON, err := json.Marshal(result)
			require.NoError(t, err)
			assert.JSONEq(t, tc.wantJSON, string(resultJSON))
		})
	}
}

func TestEncodeProviderDefinedTool(t *testing.T) {
	tests := []struct {
		name         string
		input        api.ProviderDefinedTool
		expectNil    bool
		expectBetas  []string
		wantWarnings []api.CallWarning // Expected warnings (empty means no warnings)
		wantErrMsg   string            // Empty means no error, non-empty means expect error containing this string
		want         anthropic.BetaToolUnionUnionParam
	}{
		// Computer tool tests using constructor functions
		{
			name:         "computer tool with version 20250124 (constructor)",
			input:        ComputerTool(800, 600, WithComputerVersion("20250124"), WithDisplayNumber(1)),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolComputerUse20250124Param{
				Name:            anthropic.F(anthropic.BetaToolComputerUse20250124Name("computer")),
				Type:            anthropic.F(anthropic.BetaToolComputerUse20250124TypeComputer20250124),
				DisplayWidthPx:  anthropic.Int(800),
				DisplayHeightPx: anthropic.Int(600),
				DisplayNumber:   anthropic.Int(1),
			},
		},
		{
			name:         "computer tool with version 20241022 (constructor)",
			input:        ComputerTool(800, 600, WithComputerVersion("20241022"), WithDisplayNumber(1)),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolComputerUse20241022Param{
				Name:            anthropic.F(anthropic.BetaToolComputerUse20241022Name("computer")),
				Type:            anthropic.F(anthropic.BetaToolComputerUse20241022TypeComputer20241022),
				DisplayWidthPx:  anthropic.Int(800),
				DisplayHeightPx: anthropic.Int(600),
				DisplayNumber:   anthropic.Int(1),
			},
		},
		{
			name:         "computer tool with default version (constructor)",
			input:        ComputerTool(800, 600, WithDisplayNumber(1)),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolComputerUse20250124Param{
				Name:            anthropic.F(anthropic.BetaToolComputerUse20250124Name("computer")),
				Type:            anthropic.F(anthropic.BetaToolComputerUse20250124TypeComputer20250124),
				DisplayWidthPx:  anthropic.Int(800),
				DisplayHeightPx: anthropic.Int(600),
				DisplayNumber:   anthropic.Int(1),
			},
		},

		// Computer tool tests using map[string]any args
		{
			name: "computer tool with version 20250124 (map args)",
			input: api.ProviderDefinedTool{
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
			want: anthropic.BetaToolComputerUse20250124Param{
				Name:            anthropic.F(anthropic.BetaToolComputerUse20250124Name("computer")),
				Type:            anthropic.F(anthropic.BetaToolComputerUse20250124TypeComputer20250124),
				DisplayWidthPx:  anthropic.Int(800),
				DisplayHeightPx: anthropic.Int(600),
				DisplayNumber:   anthropic.Int(1),
			},
		},
		{
			name: "computer tool with version 20241022 (map args)",
			input: api.ProviderDefinedTool{
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
			want: anthropic.BetaToolComputerUse20241022Param{
				Name:            anthropic.F(anthropic.BetaToolComputerUse20241022Name("computer")),
				Type:            anthropic.F(anthropic.BetaToolComputerUse20241022TypeComputer20241022),
				DisplayWidthPx:  anthropic.Int(800),
				DisplayHeightPx: anthropic.Int(600),
				DisplayNumber:   anthropic.Int(1),
			},
		},
		{
			name: "computer tool with default version (map args)",
			input: api.ProviderDefinedTool{
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
			want: anthropic.BetaToolComputerUse20250124Param{
				Name:            anthropic.F(anthropic.BetaToolComputerUse20250124Name("computer")),
				Type:            anthropic.F(anthropic.BetaToolComputerUse20250124TypeComputer20250124),
				DisplayWidthPx:  anthropic.Int(800),
				DisplayHeightPx: anthropic.Int(600),
				DisplayNumber:   anthropic.Int(1),
			},
		},

		// Text editor tool tests using constructor functions
		{
			name:         "text editor tool with version 20250124 (constructor)",
			input:        TextEditorTool(WithTextEditorVersion("20250124")),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolTextEditor20250124Param{
				Name: anthropic.F(anthropic.BetaToolTextEditor20250124Name("str_replace_editor")),
				Type: anthropic.F(anthropic.BetaToolTextEditor20250124TypeTextEditor20250124),
			},
		},
		{
			name:         "text editor tool with version 20241022 (constructor)",
			input:        TextEditorTool(WithTextEditorVersion("20241022")),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolTextEditor20241022Param{
				Name: anthropic.F(anthropic.BetaToolTextEditor20241022Name("str_replace_editor")),
				Type: anthropic.F(anthropic.BetaToolTextEditor20241022TypeTextEditor20241022),
			},
		},
		{
			name:         "text editor tool with default version (constructor)",
			input:        TextEditorTool(),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolTextEditor20250124Param{
				Name: anthropic.F(anthropic.BetaToolTextEditor20250124Name("str_replace_editor")),
				Type: anthropic.F(anthropic.BetaToolTextEditor20250124TypeTextEditor20250124),
			},
		},

		// Text editor tool tests using map[string]any args
		{
			name: "text editor tool with version 20250124 (map args)",
			input: api.ProviderDefinedTool{
				ID:   "anthropic.text_editor",
				Name: "str_replace_editor",
				Args: map[string]any{
					"version": "20250124",
				},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolTextEditor20250124Param{
				Name: anthropic.F(anthropic.BetaToolTextEditor20250124Name("str_replace_editor")),
				Type: anthropic.F(anthropic.BetaToolTextEditor20250124TypeTextEditor20250124),
			},
		},
		{
			name: "text editor tool with version 20241022 (map args)",
			input: api.ProviderDefinedTool{
				ID:   "anthropic.text_editor",
				Name: "str_replace_editor",
				Args: map[string]any{
					"version": "20241022",
				},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolTextEditor20241022Param{
				Name: anthropic.F(anthropic.BetaToolTextEditor20241022Name("str_replace_editor")),
				Type: anthropic.F(anthropic.BetaToolTextEditor20241022TypeTextEditor20241022),
			},
		},
		{
			name: "text editor tool with default version (map args)",
			input: api.ProviderDefinedTool{
				ID:   "anthropic.text_editor",
				Name: "str_replace_editor",
				Args: map[string]any{},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolTextEditor20250124Param{
				Name: anthropic.F(anthropic.BetaToolTextEditor20250124Name("str_replace_editor")),
				Type: anthropic.F(anthropic.BetaToolTextEditor20250124TypeTextEditor20250124),
			},
		},

		// Bash tool tests using constructor functions
		{
			name:         "bash tool with version 20250124 (constructor)",
			input:        BashTool(WithBashVersion("20250124")),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolBash20250124Param{
				Name: anthropic.F(anthropic.BetaToolBash20250124Name("bash")),
				Type: anthropic.F(anthropic.BetaToolBash20250124TypeBash20250124),
			},
		},
		{
			name:         "bash tool with version 20241022 (constructor)",
			input:        BashTool(WithBashVersion("20241022")),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolBash20241022Param{
				Name: anthropic.F(anthropic.BetaToolBash20241022Name("bash")),
				Type: anthropic.F(anthropic.BetaToolBash20241022TypeBash20241022),
			},
		},
		{
			name:         "bash tool with default version (constructor)",
			input:        BashTool(),
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolBash20250124Param{
				Name: anthropic.F(anthropic.BetaToolBash20250124Name("bash")),
				Type: anthropic.F(anthropic.BetaToolBash20250124TypeBash20250124),
			},
		},

		// Bash tool tests using map[string]any args
		{
			name: "bash tool with version 20250124 (map args)",
			input: api.ProviderDefinedTool{
				ID:   "anthropic.bash",
				Name: "bash",
				Args: map[string]any{
					"version": "20250124",
				},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolBash20250124Param{
				Name: anthropic.F(anthropic.BetaToolBash20250124Name("bash")),
				Type: anthropic.F(anthropic.BetaToolBash20250124TypeBash20250124),
			},
		},
		{
			name: "bash tool with version 20241022 (map args)",
			input: api.ProviderDefinedTool{
				ID:   "anthropic.bash",
				Name: "bash",
				Args: map[string]any{
					"version": "20241022",
				},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolBash20241022Param{
				Name: anthropic.F(anthropic.BetaToolBash20241022Name("bash")),
				Type: anthropic.F(anthropic.BetaToolBash20241022TypeBash20241022),
			},
		},
		{
			name: "bash tool with default version (map args)",
			input: api.ProviderDefinedTool{
				ID:   "anthropic.bash",
				Name: "bash",
				Args: map[string]any{},
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolBash20250124Param{
				Name: anthropic.F(anthropic.BetaToolBash20250124Name("bash")),
				Type: anthropic.F(anthropic.BetaToolBash20250124TypeBash20250124),
			},
		},

		// Error cases using map[string]any args
		{
			name: "computer tool with invalid version (map args)",
			input: api.ProviderDefinedTool{
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
			input: api.ProviderDefinedTool{
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
			input: api.ProviderDefinedTool{
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
			input: api.ProviderDefinedTool{
				ID:   "mock.unsupported",
				Name: "unsupported",
				Args: &unsupportedArgs{},
			},
			expectNil:   true,
			expectBetas: []string{},
			wantWarnings: []api.CallWarning{
				{
					Type: "unsupported-tool",
					Tool: api.ProviderDefinedTool{
						ID:   "mock.unsupported",
						Name: "unsupported",
						Args: &unsupportedArgs{},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tool, betas, warnings, err := EncodeProviderDefinedTool(tc.input)

			if tc.wantErrMsg != "" {
				assert.Error(t, err, "Expected an error")
				assert.Contains(t, err.Error(), tc.wantErrMsg, "Error message should contain expected substring")
				return
			}

			require.NoError(t, err)

			// Check warnings
			if len(tc.wantWarnings) == 0 {
				assert.Empty(t, warnings, "No warnings should be returned")
			} else {
				assert.ElementsMatch(t, tc.wantWarnings, warnings, "Warnings mismatch")
			}

			// Check betas
			assert.ElementsMatch(t, tc.expectBetas, betas, "Betas mismatch")

			// Check if tool should be nil
			if tc.expectNil {
				assert.Nil(t, tool, "Tool should be nil")
				return
			}

			require.NotNil(t, tool, "Tool should not be nil")

			// Validate the returned type matches the expected type
			assert.IsType(t, tc.want, tool, "Tool type mismatch")

			// Validate the JSON representation matches
			expectedJSON, err := json.Marshal(tc.want)
			require.NoError(t, err, "Failed to marshal expected tool to JSON")

			actualJSON, err := json.Marshal(tool)
			require.NoError(t, err, "Failed to marshal actual tool to JSON")

			assert.JSONEq(t, string(expectedJSON), string(actualJSON), "Tool JSON content mismatch for %s", tc.name)
		})
	}
}

func TestEncodeTools(t *testing.T) {
	functionTool := api.FunctionTool{
		Name:        "test_function",
		Description: "A test function",
		InputSchema: &jsonschema.Definition{
			Type: "object",
			Properties: map[string]jsonschema.Definition{
				"param1": {
					Type:        "string",
					Description: "First parameter",
				},
			},
			Required: []string{"param1"},
		},
	}

	computerTool := api.ProviderDefinedTool{
		ID:   "anthropic.computer",
		Name: "computer",
		Args: &ComputerToolArgs{
			DisplayWidthPx:  800,
			DisplayHeightPx: 600,
			DisplayNumber:   1,
		},
	}

	// Use a concrete tool type that we know won't be handled correctly
	unsupportedTool := api.ProviderDefinedTool{
		ID:   "mock.unsupported",
		Name: "unsupported",
		Args: &unsupportedArgs{},
	}

	// Helper to create expected tool choice
	autoChoice := []anthropic.BetaToolChoiceUnionParam{
		anthropic.BetaToolChoiceAutoParam{
			Type: anthropic.F(anthropic.BetaToolChoiceAutoTypeAuto),
		},
	}
	anyChoice := []anthropic.BetaToolChoiceUnionParam{
		anthropic.BetaToolChoiceAnyParam{
			Type: anthropic.F(anthropic.BetaToolChoiceAnyTypeAny),
		},
	}
	specificChoice := func(toolName string) []anthropic.BetaToolChoiceUnionParam {
		return []anthropic.BetaToolChoiceUnionParam{
			anthropic.BetaToolChoiceToolParam{
				Type: anthropic.F(anthropic.BetaToolChoiceToolTypeTool),
				Name: anthropic.String(toolName),
			},
		}
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
				Tools:      []anthropic.BetaToolUnionUnionParam{mustEncodeFunctionTool(functionTool)},
				ToolChoice: nil,
				Betas:      []anthropic.AnthropicBeta{},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name:   "provider tool with auto choice",
			tools:  []api.ToolDefinition{computerTool},
			choice: &api.ToolChoice{Type: "auto"},
			want: AnthropicTools{
				Tools:      []anthropic.BetaToolUnionUnionParam{mustEncodeComputerTool()},
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
				Tools:      []anthropic.BetaToolUnionUnionParam{mustEncodeFunctionTool(functionTool), mustEncodeComputerTool()},
				ToolChoice: anyChoice,
				Betas:      []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2025_01_24},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name:  "unsupported tool",
			tools: []api.ToolDefinition{unsupportedTool},
			want: AnthropicTools{
				Tools:      []anthropic.BetaToolUnionUnionParam{},
				ToolChoice: nil,
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
				ToolChoice: nil,
				Betas:      []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2025_01_24},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name:   "specific tool choice",
			tools:  []api.ToolDefinition{functionTool},
			choice: &api.ToolChoice{Type: "tool", ToolName: "test_function"},
			want: AnthropicTools{
				Tools:      []anthropic.BetaToolUnionUnionParam{mustEncodeFunctionTool(functionTool)},
				ToolChoice: specificChoice("test_function"),
				Betas:      []anthropic.AnthropicBeta{},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name: "object tool mode equivalent",
			tools: []api.ToolDefinition{
				api.FunctionTool{
					Name:        "test_tool",
					Description: "A test tool",
					InputSchema: &jsonschema.Definition{Type: "object"},
				},
			},
			choice: &api.ToolChoice{Type: "tool", ToolName: "test_tool"},
			want: AnthropicTools{
				Tools: []anthropic.BetaToolUnionUnionParam{
					mustEncodeFunctionTool(api.FunctionTool{
						Name:        "test_tool",
						Description: "A test tool",
						InputSchema: &jsonschema.Definition{Type: "object"},
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
				api.ProviderDefinedTool{
					ID:   "anthropic.computer",
					Name: "computer",
					Args: &ComputerToolArgs{DisplayWidthPx: 800, DisplayHeightPx: 600, DisplayNumber: 1},
				},
				api.ProviderDefinedTool{
					ID:   "anthropic.text_editor",
					Name: "str_replace_editor",
					Args: &TextEditorToolArgs{},
				},
				api.ProviderDefinedTool{
					ID:   "anthropic.bash",
					Name: "bash",
					Args: &BashToolArgs{},
				},
			},
			want: AnthropicTools{
				Tools: []anthropic.BetaToolUnionUnionParam{
					mustEncodeComputerTool(),
					mustEncodeTextEditorTool(),
					mustEncodeBashTool(),
				},
				ToolChoice: nil,
				Betas:      []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2025_01_24},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name:  "mixed supported and unsupported tools",
			tools: []api.ToolDefinition{functionTool, unsupportedTool, computerTool},
			want: AnthropicTools{
				Tools:      []anthropic.BetaToolUnionUnionParam{mustEncodeFunctionTool(functionTool), mustEncodeComputerTool()},
				ToolChoice: nil,
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
func mustEncodeFunctionTool(tool api.FunctionTool) anthropic.BetaToolUnionUnionParam {
	result, err := EncodeFunctionTool(tool)
	if err != nil {
		panic(err)
	}
	return result
}

func mustEncodeComputerTool() anthropic.BetaToolUnionUnionParam {
	return anthropic.BetaToolComputerUse20250124Param{
		Name:            anthropic.F(anthropic.BetaToolComputerUse20250124Name("computer")),
		Type:            anthropic.F(anthropic.BetaToolComputerUse20250124TypeComputer20250124),
		DisplayWidthPx:  anthropic.Int(800),
		DisplayHeightPx: anthropic.Int(600),
		DisplayNumber:   anthropic.Int(1),
	}
}

func mustEncodeTextEditorTool() anthropic.BetaToolUnionUnionParam {
	return anthropic.BetaToolTextEditor20250124Param{
		Name: anthropic.F(anthropic.BetaToolTextEditor20250124Name("str_replace_editor")),
		Type: anthropic.F(anthropic.BetaToolTextEditor20250124TypeTextEditor20250124),
	}
}

func mustEncodeBashTool() anthropic.BetaToolUnionUnionParam {
	return anthropic.BetaToolBash20250124Param{
		Name: anthropic.F(anthropic.BetaToolBash20250124Name("bash")),
		Type: anthropic.F(anthropic.BetaToolBash20250124TypeBash20250124),
	}
}

// unsupportedArgs is used as Args for testing unsupported tools
type unsupportedArgs struct{}
