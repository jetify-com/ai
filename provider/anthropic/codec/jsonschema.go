package codec

import (
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/google/jsonschema-go/jsonschema"
)

// encodeSchema converts a jsonschema.Schema to map[string]any for Anthropic API.
func encodeSchema(schema *jsonschema.Schema) (map[string]any, error) {
	if schema == nil {
		return nil, nil
	}

	// Marshal to JSON and unmarshal back to interface{} to convert the types
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	// Convert {"not": {}} patterns to false for additionalProperties
	// This is needed because Anthropic expects "additionalProperties": false
	// but MCP jsonschema represents false as {"not": {}}
	normalizeSchemaMap(result)

	return result, nil
}

// normalizeSchemaMap recursively converts {"not": {}} to false in a schema map
func normalizeSchemaMap(schemaMap map[string]interface{}) {
	for key, value := range schemaMap {
		switch v := value.(type) {
		case map[string]interface{}:
			// Check if this is a {"not": {}} pattern
			if key == "additionalProperties" && len(v) == 1 {
				if not, hasNot := v["not"]; hasNot {
					if notMap, isMap := not.(map[string]interface{}); isMap && len(notMap) == 0 {
						schemaMap[key] = false
						continue
					}
				}
			}
			// Recursively process nested objects
			normalizeSchemaMap(v)
		case []interface{}:
			// Process arrays of schemas
			for _, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					normalizeSchemaMap(itemMap)
				}
			}
		}
	}
}

// encodeInputSchema converts the JSON schema definition to Anthropic's schema format
func encodeInputSchema(schema *jsonschema.Schema) (anthropic.BetaToolInputSchemaParam, error) {
	// Verify the schema type is "object"
	// TODO: When the Anthropic SDK supports other type values or union types,
	// we should update this to pass through the original type information.
	if schema.Type != "" && schema.Type != "object" {
		return anthropic.BetaToolInputSchemaParam{}, fmt.Errorf("unsupported schema type: %s, only 'object' is supported", schema.Type)
	}
	// Reject schemas with multiple types (union types)
	if len(schema.Types) > 0 {
		return anthropic.BetaToolInputSchemaParam{}, fmt.Errorf("unsupported schema with multiple types: %v, only single type 'object' is supported", schema.Types)
	}

	// Convert schema to map
	schemaMap, err := encodeSchema(schema)
	if err != nil {
		return anthropic.BetaToolInputSchemaParam{}, err
	}

	// Create the input schema with the type field
	inputSchema := anthropic.BetaToolInputSchemaParam{
		Type: "object",
	}

	// Add properties and required fields from the schema map
	if properties, ok := schemaMap["properties"]; ok {
		inputSchema.Properties = properties
	}
	if required, ok := schemaMap["required"]; ok {
		if reqArray, ok := required.([]interface{}); ok {
			strRequired := make([]string, len(reqArray))
			for i, r := range reqArray {
				if str, ok := r.(string); ok {
					strRequired[i] = str
				}
			}
			inputSchema.Required = strRequired
		}
	}

	// Add any other fields to ExtraFields
	for k, v := range schemaMap {
		if k != "type" && k != "properties" && k != "required" {
			if inputSchema.ExtraFields == nil {
				inputSchema.ExtraFields = make(map[string]interface{})
			}
			inputSchema.ExtraFields[k] = v
		}
	}

	return inputSchema, nil
}
