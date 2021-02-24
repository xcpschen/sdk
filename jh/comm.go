package jh

import (
	"encoding/json"
)

// InArr check elment is in arrry b
func InArr(a string, b []string) bool {
	for _, j := range b {
		if a == j {
			return true
		}
	}
	return false
}

// JsonDecode
func JsonDecode(b []byte, obj interface{}) error {
	return json.Unmarshal(b, obj)
}

func JsonEncode(obj interface{}) string {
	b, _ := json.Marshal(obj)
	return string(b)
}
