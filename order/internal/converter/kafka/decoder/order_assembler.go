package decoder

import (
	"encoding/json"

	"order/internal/model"
)

// DecodeOrderAssembled декодирует JSON в OrderAssembledEvent
func DecodeOrderAssembled(data []byte) (*model.OrderAssembledEvent, error) {
	var event model.OrderAssembledEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}