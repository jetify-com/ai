package codec

import (
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

func TestEncodeSchema(t *testing.T) {
	tests := []struct {
		name    string
		input   *jsonschema.Schema
		want    string // JSON string
		wantErr bool
	}{
		{
			name:  "nil schema",
			input: nil,
			want:  "null",
		},
		{
			name: "simple object schema",
			input: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"name": {
						Type:        "string",
						Description: "The name",
					},
					"age": {
						Type:        "number",
						Description: "The age",
					},
				},
				Required: []string{"name"},
			},
			want: `{
				"type": "object",
				"properties": {
					"name": {
						"type": "string",
						"description": "The name"
					},
					"age": {
						"type": "number",
						"description": "The age"
					}
				},
				"required": ["name"]
			}`,
		},
		{
			name: "schema with additionalProperties false",
			input: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"field": {
						Type: "string",
					},
				},
				AdditionalProperties: api.FalseSchema(),
			},
			want: `{
				"type": "object",
				"properties": {
					"field": {
						"type": "string"
					}
				},
				"additionalProperties": false
			}`,
		},
		{
			name: "nested schema with additionalProperties false",
			input: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"nested": {
						Type: "object",
						Properties: map[string]*jsonschema.Schema{
							"inner": {
								Type: "string",
							},
						},
						AdditionalProperties: api.FalseSchema(),
					},
				},
			},
			want: `{
				"type": "object",
				"properties": {
					"nested": {
						"type": "object",
						"properties": {
							"inner": {
								"type": "string"
							}
						},
						"additionalProperties": false
					}
				}
			}`,
		},
		{
			name: "schema with array of schemas",
			input: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"items": {
						Type: "array",
						Items: &jsonschema.Schema{
							Type: "object",
							Properties: map[string]*jsonschema.Schema{
								"id": {
									Type: "string",
								},
							},
							AdditionalProperties: api.FalseSchema(),
						},
					},
				},
			},
			want: `{
				"type": "object",
				"properties": {
					"items": {
						"type": "array",
						"items": {
							"type": "object",
							"properties": {
								"id": {
									"type": "string"
								}
							},
							"additionalProperties": false
						}
					}
				}
			}`,
		},
		{
			name: "schema with allOf containing additionalProperties",
			input: &jsonschema.Schema{
				AllOf: []*jsonschema.Schema{
					{
						Type:                 "object",
						AdditionalProperties: api.FalseSchema(),
					},
				},
			},
			want: `{
				"allOf": [{
					"type": "object",
					"additionalProperties": false
				}]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeSchema(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Convert result to JSON for comparison
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)

			assert.JSONEq(t, tt.want, string(gotJSON))
		})
	}
}

func TestNormalizeSchemaMap(t *testing.T) {
	tests := []struct {
		name  string
		input string // JSON string that will be unmarshaled
		want  string // Expected JSON output
	}{
		{
			name: "additionalProperties with not empty object",
			input: `{
				"additionalProperties": {
					"not": {}
				}
			}`,
			want: `{
				"additionalProperties": false
			}`,
		},
		{
			name: "additionalProperties with not non-empty object",
			input: `{
				"additionalProperties": {
					"not": {
						"type": "string"
					}
				}
			}`,
			want: `{
				"additionalProperties": {
					"not": {
						"type": "string"
					}
				}
			}`,
		},
		{
			name: "additionalProperties with non-not pattern",
			input: `{
				"additionalProperties": {
					"type": "string"
				}
			}`,
			want: `{
				"additionalProperties": {
					"type": "string"
				}
			}`,
		},
		{
			name: "nested object with additionalProperties",
			input: `{
				"properties": {
					"nested": {
						"type": "object",
						"additionalProperties": {
							"not": {}
						}
					}
				}
			}`,
			want: `{
				"properties": {
					"nested": {
						"type": "object",
						"additionalProperties": false
					}
				}
			}`,
		},
		{
			name: "array with schema items",
			input: `{
				"items": [
					{
						"type": "object",
						"additionalProperties": {
							"not": {}
						}
					},
					"string"
				]
			}`,
			want: `{
				"items": [
					{
						"type": "object",
						"additionalProperties": false
					},
					"string"
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input map[string]interface{}
			err := json.Unmarshal([]byte(tt.input), &input)
			require.NoError(t, err)

			normalizeSchemaMap(input)

			// Convert result to JSON for comparison
			gotJSON, err := json.Marshal(input)
			require.NoError(t, err)

			assert.JSONEq(t, tt.want, string(gotJSON))
		})
	}
}

func TestEncodeInputSchema(t *testing.T) {
	tests := []struct {
		name    string
		input   *jsonschema.Schema
		want    string // Expected JSON output
		wantErr bool
		errMsg  string
	}{
		{
			name: "simple object schema",
			input: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"name": {
						Type:        "string",
						Description: "The name",
					},
				},
				Required: []string{"name"},
			},
			want: `{
				"type": "object",
				"properties": {
					"name": {
						"type": "string",
						"description": "The name"
					}
				},
				"required": ["name"]
			}`,
		},
		{
			name: "object schema with empty type",
			input: &jsonschema.Schema{
				Properties: map[string]*jsonschema.Schema{
					"field": {
						Type: "string",
					},
				},
			},
			want: `{
				"type": "object",
				"properties": {
					"field": {
						"type": "string"
					}
				}
			}`,
		},
		{
			name: "object schema with multiple types including object",
			input: &jsonschema.Schema{
				Types: []string{"object", "null"},
				Properties: map[string]*jsonschema.Schema{
					"field": {
						Type: "string",
					},
				},
			},
			wantErr: true,
			errMsg:  "unsupported schema with multiple types: [object null], only single type 'object' is supported",
		},
		{
			name: "non-object schema type",
			input: &jsonschema.Schema{
				Type: "array",
			},
			wantErr: true,
			errMsg:  "unsupported schema type: array, only 'object' is supported",
		},
		{
			name: "multiple types without object",
			input: &jsonschema.Schema{
				Types: []string{"string", "number"},
			},
			wantErr: true,
			errMsg:  "unsupported schema with multiple types: [string number], only single type 'object' is supported",
		},
		{
			name: "schema with no extra fields",
			input: &jsonschema.Schema{
				Type: "object",
			},
			want: `{
				"type": "object"
			}`,
		},
		{
			name: "schema with additionalProperties false",
			input: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"field": {
						Type: "string",
					},
				},
				AdditionalProperties: api.FalseSchema(),
			},
			want: `{
				"type": "object",
				"properties": {
					"field": {
						"type": "string"
					}
				},
				"additionalProperties": false
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeInputSchema(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.EqualError(t, err, tt.errMsg)
				}
				return
			}
			require.NoError(t, err)

			// Marshal the result to JSON for comparison
			gotJSON, err := json.Marshal(got)
			require.NoError(t, err)

			assert.JSONEq(t, tt.want, string(gotJSON))
		})
	}
}

func TestEncodeInputSchema_Integration(t *testing.T) {
	// Test that encodeInputSchema properly integrates with encodeSchema
	schema := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"nested": {
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"field": {
						Type: "string",
					},
				},
				AdditionalProperties: api.FalseSchema(),
			},
		},
	}

	got, err := encodeInputSchema(schema)
	require.NoError(t, err)

	// Marshal the result and verify nested additionalProperties was normalized
	gotJSON, err := json.Marshal(got)
	require.NoError(t, err)

	expectedJSON := `{
		"type": "object",
		"properties": {
			"nested": {
				"type": "object",
				"properties": {
					"field": {
						"type": "string"
					}
				},
				"additionalProperties": false
			}
		}
	}`

	assert.JSONEq(t, expectedJSON, string(gotJSON))
}
