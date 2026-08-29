package decoder

import (
	"encoding/json"

	"notification/internal/model"
)

func DecodeOrderAssembled(data []byte) (*model.OrderAssembledEvent, error) {
	var event model.OrderAssembledEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}