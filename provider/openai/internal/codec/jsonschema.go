package codec

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// encodeSchema converts a jsonschema.Schema to a map representation suitable for OpenAI API.
// It handles the conversion of JSON Schema 2020-12 format to the format expected by OpenAI.
// TODO: promote to a framework-level function
func encodeSchema(schema *jsonschema.Schema) (map[string]any, error) {
	if schema == nil {
		return nil, nil
	}

	// Marshal to JSON and unmarshal back to interface{} to convert the types
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal properties: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal properties: %w", err)
	}

	// Convert {"not": {}} patterns to false throughout the schema
	normalizeSchemaMap(result)

	return result, nil
}

// normalizeSchemaMap recursively converts {"not": {}} to false in a schema map.
// This is needed because OpenAI expects "additionalProperties": false
// but MCP jsonschema represents false as {"not": {}}
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
