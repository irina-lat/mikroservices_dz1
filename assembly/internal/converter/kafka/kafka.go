package kafka

import (
	"encoding/json"
)

// EncodeEvent кодирует событие в JSON
func EncodeEvent(event interface{}) ([]byte, error) {
	return json.Marshal(event)
}

// DecodeEvent декодирует JSON в событие
func DecodeEvent(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}