package codec

import (
	"encoding/json"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
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
		{
			name: "complex schema with multiple features",
			input: &jsonschema.Schema{
				Type:        "object",
				Title:       "Complex Schema",
				Description: "A complex schema with various features",
				Properties: map[string]*jsonschema.Schema{
					"id": {
						Type:    "string",
						Pattern: "^[a-zA-Z0-9]+$",
					},
					"metadata": {
						Type: "object",
						Properties: map[string]*jsonschema.Schema{
							"tags": {
								Type: "array",
								Items: &jsonschema.Schema{
									Type: "string",
								},
							},
						},
						AdditionalProperties: api.FalseSchema(),
					},
				},
				Required: []string{"id"},
			},
			want: `{
				"type": "object",
				"title": "Complex Schema",
				"description": "A complex schema with various features",
				"properties": {
					"id": {
						"type": "string",
						"pattern": "^[a-zA-Z0-9]+$"
					},
					"metadata": {
						"type": "object",
						"properties": {
							"tags": {
								"type": "array",
								"items": {
									"type": "string"
								}
							}
						},
						"additionalProperties": false
					}
				},
				"required": ["id"]
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
			name: "additionalProperties with multiple keys",
			input: `{
				"additionalProperties": {
					"not": {},
					"type": "string"
				}
			}`,
			want: `{
				"additionalProperties": {
					"not": {},
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
					"string",
					42,
					null
				]
			}`,
			want: `{
				"items": [
					{
						"type": "object",
						"additionalProperties": false
					},
					"string",
					42,
					null
				]
			}`,
		},
		{
			name: "deeply nested schemas",
			input: `{
				"properties": {
					"level1": {
						"properties": {
							"level2": {
								"properties": {
									"level3": {
										"type": "object",
										"additionalProperties": {
											"not": {}
										}
									}
								},
								"additionalProperties": {
									"not": {}
								}
							}
						}
					}
				}
			}`,
			want: `{
				"properties": {
					"level1": {
						"properties": {
							"level2": {
								"properties": {
									"level3": {
										"type": "object",
										"additionalProperties": false
									}
								},
								"additionalProperties": false
							}
						}
					}
				}
			}`,
		},
		{
			name: "allOf with nested additionalProperties",
			input: `{
				"allOf": [
					{
						"type": "object",
						"additionalProperties": {
							"not": {}
						}
					},
					{
						"properties": {
							"field": {
								"type": "object",
								"additionalProperties": {
									"not": {}
								}
							}
						}
					}
				]
			}`,
			want: `{
				"allOf": [
					{
						"type": "object",
						"additionalProperties": false
					},
					{
						"properties": {
							"field": {
								"type": "object",
								"additionalProperties": false
							}
						}
					}
				]
			}`,
		},
		{
			name: "not pattern in non-additionalProperties context",
			input: `{
				"someOtherKey": {
					"not": {}
				}
			}`,
			want: `{
				"someOtherKey": {
					"not": {}
				}
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

func TestEncodeSchema_EdgeCases(t *testing.T) {
	t.Run("schema with only additionalProperties", func(t *testing.T) {
		schema := &jsonschema.Schema{
			AdditionalProperties: api.FalseSchema(),
		}

		got, err := encodeSchema(schema)
		require.NoError(t, err)

		gotJSON, err := json.Marshal(got)
		require.NoError(t, err)

		expectedJSON := `{
			"additionalProperties": false
		}`

		assert.JSONEq(t, expectedJSON, string(gotJSON))
	})

	t.Run("empty schema", func(t *testing.T) {
		schema := &jsonschema.Schema{}

		got, err := encodeSchema(schema)
		require.NoError(t, err)

		gotJSON, err := json.Marshal(got)
		require.NoError(t, err)

		assert.JSONEq(t, "{}", string(gotJSON))
	})
}
