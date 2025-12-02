package codec

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
)

// encodeSchema converts a jsonschema.Schema to a map representation suitable for OpenAI API.
// It handles the conversion of JSON Schema 2020-12 format to the format expected by OpenAI.
// TODO: promote to a framework-level function
func encodeSchema(schema *jsonschema.Schema) (map[string]any, error) {
	if schema == nil {
		return nil, nil
	}

	// Enforce OpenAI restrictions
	// https://platform.openai.com/docs/guides/structured-outputs#root-objects-must-not-be-anyof-and-must-be-an-object
	// NOTE: we could simply encode the input schema, pass it through to OpenAI and let it return an error, but there are
	// other encoding rules we want to enforce later, and limiting the scope here allows us to limit the scope later.
	if schema.Type != "object" {
		return nil, fmt.Errorf("schema root must be of type object, got: %s", schema.Type)
	}
	if schema.AnyOf != nil {
		return nil, fmt.Errorf("schema root cannot use AnyOf")
	}

	// Marshal to JSON and unmarshal back to interface{} to convert the types
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal properties: %w", err)
	}

	// Convert the schema "false" back to an empty object.
	if bytes.Equal(data, []byte("true")) {
		return make(map[string]any), nil
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal properties: %w\n\n%s", err, data)
	}

	// Ensure properties field is set, even if it's empty. It's unclear whether OpenAI requires
	// this to be set for nested schema objects too. For now we only set it at the top-level.
	if _, ok := result["properties"]; !ok {
		result["properties"] = map[string]any{}
	}

	// Convert {"not": {}} patterns to false throughout the schema
	normalizeSchemaMap(result)

	return result, nil
}

// normalizeSchemaMap recursively converts {"not": {}} to false in a schema map.
// This is needed because OpenAI expects "additionalProperties": false
// but MCP jsonschema represents false as {"not": {}}
func normalizeSchemaMap(schemaMap map[string]any) {
	for key, value := range schemaMap {
		switch v := value.(type) {
		case map[string]any:
			// Check if this is a {"not": {}} pattern
			if key == "additionalProperties" && len(v) == 1 {
				if not, hasNot := v["not"]; hasNot {
					if notMap, isMap := not.(map[string]any); isMap && len(notMap) == 0 {
						schemaMap[key] = false
						continue
					}
				}
			}
			// Recursively process nested objects
			normalizeSchemaMap(v)
		case []any:
			// Process arrays of schemas
			for _, item := range v {
				if itemMap, ok := item.(map[string]any); ok {
					normalizeSchemaMap(itemMap)
				}
			}
		}
	}
}
