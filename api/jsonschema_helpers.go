package api

import (
	"github.com/google/jsonschema-go/jsonschema"
)

// TODO: remove ... it really should be part of jsonschema.

// FalseSchema returns a new Schema that fails to validate any value.
// This represents "additionalProperties": false in JSON Schema.
// In JSON Schema 2020-12, false can also be represented as {"not": {}}
// and that's what the `jsonschema` package does.
func FalseSchema() *jsonschema.Schema {
	return &jsonschema.Schema{Not: &jsonschema.Schema{}}
}
