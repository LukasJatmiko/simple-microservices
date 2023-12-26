package utils

import (
	"encoding/json"
)

type Map map[string]interface{}

func (m Map) Decode(output any) error {
	if raw, e := json.Marshal(m); e != nil {
		return e
	} else {
		return json.Unmarshal(raw, output)
	}
}
