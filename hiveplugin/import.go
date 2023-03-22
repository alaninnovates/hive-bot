package hiveplugin

import (
	"context"
	"encoding/json"
	"github.com/qri-io/jsonschema"
)

var schemaData = []byte(`{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "patternProperties": {
    "^[a-zA-Z]+$": {
      "type": "object",
      "properties": {
        "amount": {
          "type": "integer"
        },
        "gifted": {
          "type": "boolean"
        }
      },
      "required": [
        "amount",
        "gifted"
      ],
      "additionalProperties": false
    }
  },
  "maxProperties": 50,
  "additionalProperties": false
}`)

type BeeData struct {
	Amount int  `json:"amount"`
	Gifted bool `json:"gifted"`
}

type ImportedHive map[string]BeeData

var cachedSchema *jsonschema.Schema

func GetHiveSchema() *jsonschema.Schema {
	if cachedSchema != nil {
		return cachedSchema
	}
	rs := &jsonschema.Schema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}
	cachedSchema = rs
	return rs
}

func ValidateHiveJson(jsonStr string) bool {
	schema := GetHiveSchema()
	errs, err := schema.ValidateBytes(context.Background(), []byte(jsonStr))
	if len(errs) > 0 || err != nil {
		return false
	}
	return true
}

func ParseHiveJson(jsonStr string) ImportedHive {
	var data ImportedHive
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		panic(err)
	}
	return data
}
