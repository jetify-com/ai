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
		{
			name: "computer tool with version 20250124",
			input: &ComputerUseTool{
				Version:         "20250124",
				DisplayWidthPx:  800,
				DisplayHeightPx: 600,
				DisplayNumber:   1,
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
			name: "computer tool with version 20241022",
			input: &ComputerUseTool{
				Version:         "20241022",
				DisplayWidthPx:  800,
				DisplayHeightPx: 600,
				DisplayNumber:   1,
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
			name: "computer tool with default version",
			input: &ComputerUseTool{
				DisplayWidthPx:  800,
				DisplayHeightPx: 600,
				DisplayNumber:   1,
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
			name: "text editor tool with version 20250124",
			input: &TextEditorTool{
				Version: "20250124",
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolTextEditor20250124Param{
				Name: anthropic.F(anthropic.BetaToolTextEditor20250124Name("str_replace_editor")),
				Type: anthropic.F(anthropic.BetaToolTextEditor20250124TypeTextEditor20250124),
			},
		},
		{
			name: "text editor tool with version 20241022",
			input: &TextEditorTool{
				Version: "20241022",
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolTextEditor20241022Param{
				Name: anthropic.F(anthropic.BetaToolTextEditor20241022Name("str_replace_editor")),
				Type: anthropic.F(anthropic.BetaToolTextEditor20241022TypeTextEditor20241022),
			},
		},
		{
			name:         "text editor tool with default version",
			input:        &TextEditorTool{},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolTextEditor20250124Param{
				Name: anthropic.F(anthropic.BetaToolTextEditor20250124Name("str_replace_editor")),
				Type: anthropic.F(anthropic.BetaToolTextEditor20250124TypeTextEditor20250124),
			},
		},
		{
			name: "bash tool with version 20250124",
			input: &BashTool{
				Version: "20250124",
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolBash20250124Param{
				Name: anthropic.F(anthropic.BetaToolBash20250124Name("bash")),
				Type: anthropic.F(anthropic.BetaToolBash20250124TypeBash20250124),
			},
		},
		{
			name: "bash tool with version 20241022",
			input: &BashTool{
				Version: "20241022",
			},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2024_10_22},
			wantWarnings: nil,
			want: anthropic.BetaToolBash20241022Param{
				Name: anthropic.F(anthropic.BetaToolBash20241022Name("bash")),
				Type: anthropic.F(anthropic.BetaToolBash20241022TypeBash20241022),
			},
		},
		{
			name:         "bash tool with default version",
			input:        &BashTool{},
			expectBetas:  []string{anthropic.AnthropicBetaComputerUse2025_01_24},
			wantWarnings: nil,
			want: anthropic.BetaToolBash20250124Param{
				Name: anthropic.F(anthropic.BetaToolBash20250124Name("bash")),
				Type: anthropic.F(anthropic.BetaToolBash20250124TypeBash20250124),
			},
		},
		{
			name: "computer tool with invalid version",
			input: &ComputerUseTool{
				Version:         "invalid",
				DisplayWidthPx:  800,
				DisplayHeightPx: 600,
			},
			wantErrMsg: "unsupported computer tool version",
		},
		{
			name: "text editor tool with invalid version",
			input: &TextEditorTool{
				Version: "invalid",
			},
			wantErrMsg: "unsupported text editor tool version",
		},
		{
			name: "bash tool with invalid version",
			input: &BashTool{
				Version: "invalid",
			},
			wantErrMsg: "unsupported bash tool version",
		},
		{
			name:        "unsupported tool type",
			input:       &mockUnsupportedTool{},
			expectNil:   true,
			expectBetas: []string{},
			wantWarnings: []api.CallWarning{
				{
					Type: "unsupported-tool",
					Tool: &mockUnsupportedTool{},
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

func TestEncodeRegularMode(t *testing.T) {
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

	computerTool := &ComputerUseTool{
		DisplayWidthPx:  800,
		DisplayHeightPx: 600,
		DisplayNumber:   1,
	}

	// Use a concrete tool type that we know won't be handled correctly
	unsupportedTool := &mockUnsupportedTool{}

	tests := []struct {
		name             string
		input            api.RegularMode
		expectTools      int
		expectBetas      int
		expectWarnings   int
		expectNilTools   bool
		expectNilChoice  bool
		expectedBetas    []anthropic.AnthropicBeta
		expectedWarnings []api.CallWarning
		wantErr          bool
	}{
		{
			name: "no tools",
			input: api.RegularMode{
				Tools:      nil,
				ToolChoice: nil,
			},
			expectNilTools:   true,
			expectNilChoice:  true,
			expectedBetas:    []anthropic.AnthropicBeta{},
			expectedWarnings: []api.CallWarning{},
		},
		{
			name: "function tool",
			input: api.RegularMode{
				Tools: []api.ToolDefinition{
					functionTool,
				},
				ToolChoice: nil,
			},
			expectTools:      1,
			expectBetas:      0,
			expectWarnings:   0,
			expectNilChoice:  true,
			expectedBetas:    []anthropic.AnthropicBeta{},
			expectedWarnings: []api.CallWarning{},
		},
		{
			name: "provider tool",
			input: api.RegularMode{
				Tools: []api.ToolDefinition{
					computerTool,
				},
				ToolChoice: &api.ToolChoice{
					Type: "auto",
				},
			},
			expectTools:      1,
			expectBetas:      1,
			expectWarnings:   0,
			expectedBetas:    []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2025_01_24},
			expectedWarnings: []api.CallWarning{},
		},
		{
			name: "mixed tools",
			input: api.RegularMode{
				Tools: []api.ToolDefinition{
					functionTool,
					computerTool,
				},
				ToolChoice: &api.ToolChoice{
					Type: "required",
				},
			},
			expectTools:      2,
			expectBetas:      1,
			expectWarnings:   0,
			expectedBetas:    []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2025_01_24},
			expectedWarnings: []api.CallWarning{},
		},
		{
			name: "unsupported tool",
			input: api.RegularMode{
				Tools: []api.ToolDefinition{
					unsupportedTool,
				},
				ToolChoice: nil,
			},
			expectTools:     0,
			expectBetas:     0,
			expectWarnings:  1,
			expectNilChoice: true,
			expectedBetas:   []anthropic.AnthropicBeta{},
			expectedWarnings: []api.CallWarning{
				{
					Type: "unsupported-tool",
					Tool: unsupportedTool,
				},
			},
		},
		{
			name: "unsupported tool choice",
			input: api.RegularMode{
				Tools: []api.ToolDefinition{
					functionTool,
				},
				ToolChoice: &api.ToolChoice{
					Type: "invalid_type",
				},
			},
			wantErr: true,
		},
		{
			name: "none tool choice",
			input: api.RegularMode{
				Tools: []api.ToolDefinition{
					functionTool,
					computerTool,
				},
				ToolChoice: &api.ToolChoice{
					Type: "none",
				},
			},
			expectTools:      0, // For "none" tool choice, tools should be nil
			expectBetas:      1, // But we still expect betas
			expectWarnings:   0,
			expectNilTools:   true,
			expectNilChoice:  true,
			expectedBetas:    []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2025_01_24},
			expectedWarnings: []api.CallWarning{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := encodeRegularMode(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tc.expectNilTools && tc.expectNilChoice && tc.expectBetas == 0 && tc.expectWarnings == 0 {
				assert.Equal(t, AnthropicTools{}, result, "Result should be an empty struct when no tools, no choice, no betas, and no warnings")
				return
			}

			require.NotNil(t, result)

			if tc.expectNilTools {
				assert.Nil(t, result.Tools, "Tools should be nil")
			} else {
				assert.Len(t, result.Tools, tc.expectTools, "Unexpected number of tools")
			}

			if tc.expectNilChoice {
				assert.Nil(t, result.ToolChoice, "ToolChoice should be nil")
			} else {
				assert.NotNil(t, result.ToolChoice, "ToolChoice should not be nil")
			}

			assert.Len(t, result.Betas, tc.expectBetas, "Unexpected number of betas")
			assert.Len(t, result.Warnings, tc.expectWarnings, "Unexpected number of warnings")

			// Check specific beta values
			assert.ElementsMatch(t, tc.expectedBetas, result.Betas, "Beta values mismatch")

			// Check specific warning content
			assert.ElementsMatch(t, tc.expectedWarnings, result.Warnings, "Warnings mismatch")
		})
	}
}

func TestEncodeToolMode(t *testing.T) {
	tests := []struct {
		name    string
		input   api.ModeConfig
		want    AnthropicTools
		wantErr bool
	}{
		{
			name:  "regular mode without tools",
			input: api.RegularMode{},
			want:  AnthropicTools{},
		},
		{
			name: "regular mode with function tool",
			input: api.RegularMode{
				Tools: []api.ToolDefinition{
					api.FunctionTool{
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
				},
			},
			want: AnthropicTools{
				Tools: []anthropic.BetaToolUnionUnionParam{
					anthropic.BetaToolParam{
						Name:        anthropic.String("test_function"),
						Description: anthropic.String("A test function"),
						InputSchema: anthropic.F(anthropic.BetaToolInputSchemaParam{
							Type: anthropic.F(anthropic.BetaToolInputSchemaTypeObject),
							Properties: anthropic.F[interface{}](map[string]interface{}{
								"param1": map[string]interface{}{
									"type":        "string",
									"description": "First parameter",
								},
							}),
							ExtraFields: map[string]interface{}{
								"required": []string{"param1"},
							},
						}),
					},
				},
				Betas:    []anthropic.AnthropicBeta{},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "object json mode",
			input:   api.ObjectJSONMode{},
			wantErr: true,
		},
		{
			name: "object tool mode",
			input: api.ObjectToolMode{
				Tool: api.FunctionTool{
					Name:        "test_tool",
					Description: "A test tool",
					InputSchema: &jsonschema.Definition{
						Type: "object",
					},
				},
			},
			want: AnthropicTools{
				Tools: []anthropic.BetaToolUnionUnionParam{
					anthropic.BetaToolParam{
						Name:        anthropic.String("test_tool"),
						Description: anthropic.String("A test tool"),
						InputSchema: anthropic.F(anthropic.BetaToolInputSchemaParam{
							Type: anthropic.F(anthropic.BetaToolInputSchemaTypeObject),
						}),
					},
				},
				ToolChoice: []anthropic.BetaToolChoiceUnionParam{
					anthropic.BetaToolChoiceToolParam{
						Type: anthropic.F(anthropic.BetaToolChoiceToolTypeTool),
						Name: anthropic.String("test_tool"),
					},
				},
				Betas:    []anthropic.AnthropicBeta{},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name: "regular mode with tool choice none",
			input: api.RegularMode{
				Tools: []api.ToolDefinition{
					api.FunctionTool{
						Name:        "test_function",
						Description: "A test function",
						InputSchema: &jsonschema.Definition{
							Type: "object",
						},
					},
				},
				ToolChoice: &api.ToolChoice{
					Type: "none",
				},
			},
			want: AnthropicTools{
				Tools:      nil,
				ToolChoice: nil,
				Betas:      []anthropic.AnthropicBeta{},
				Warnings:   []api.CallWarning{},
			},
		},
		{
			name: "regular mode with provider defined tool",
			input: api.RegularMode{
				Tools: []api.ToolDefinition{
					&ComputerUseTool{
						DisplayWidthPx:  800,
						DisplayHeightPx: 600,
						DisplayNumber:   1,
					},
				},
			},
			want: AnthropicTools{
				Tools: []anthropic.BetaToolUnionUnionParam{
					anthropic.BetaToolComputerUse20250124Param{
						Name:            anthropic.F(anthropic.BetaToolComputerUse20250124Name("computer")),
						Type:            anthropic.F(anthropic.BetaToolComputerUse20250124TypeComputer20250124),
						DisplayWidthPx:  anthropic.Int(800),
						DisplayHeightPx: anthropic.Int(600),
						DisplayNumber:   anthropic.Int(1),
					},
				},
				Betas: []anthropic.AnthropicBeta{
					anthropic.AnthropicBetaComputerUse2025_01_24,
				},
				Warnings: []api.CallWarning{},
			},
		},
		{
			name:    "unsupported mode type",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := EncodeToolMode(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, result)
		})
	}
}

// mockUnsupportedTool implements the ProviderDefinedTool interface for testing unsupported tools
type mockUnsupportedTool struct{}

func (t *mockUnsupportedTool) ToolType() string { return "provider-defined" }
func (t *mockUnsupportedTool) ID() string       { return "mock.unsupported" }
func (t *mockUnsupportedTool) Name() string     { return "unsupported" }
