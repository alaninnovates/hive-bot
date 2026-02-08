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
    "type": {
       "type": "string"
    },
    "basic|bomber|brave|bumble|cool|hasty|looker|rad|rascal|stubborn|bubble|bucko|commander|demo|exhausted|fire|frosty|honey|rage|riley|shocked|baby|carpenter|demon|diamond|lion|music|ninja|shy|buoyant|fuzzy|precise|spicy|tadpole|vector|bear|cobalt|crimson|digital|festive|gummy|photon|puppy|tabby|vicious|windy": {
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
  "maxProperties": 51,
  "additionalProperties": false
}`)

type BeeData struct {
	Amount int  `json:"amount"`
	Gifted bool `json:"gifted"`
}

type ImportedHive map[string]any

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
