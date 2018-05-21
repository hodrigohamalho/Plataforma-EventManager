package infra

import (
	"encoding/json"
)

//Clone make a deep clone from objects
func Clone(from interface{}, to interface{}) error {
	if data, err := json.Marshal(from); err != nil {
		return err
	} else {
		return json.Unmarshal(data, to)
	}
}
