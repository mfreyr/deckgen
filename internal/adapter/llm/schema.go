package llm

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
)

// GenerateSchema creates a JSON schema definition from a Go struct.
func GenerateSchema[T any]() map[string]interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	schemaJson, err := schema.MarshalJSON()
	if err != nil {
		panic(err)
	}
	var schemaObj map[string]interface{}
	err = json.Unmarshal(schemaJson, &schemaObj)
	if err != nil {
		panic(err)
	}
	return schemaObj
}
