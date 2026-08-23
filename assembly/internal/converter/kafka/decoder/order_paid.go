package decoder

import (
	"encoding/json"

	"assembly/internal/model"
)

// DecodeOrderPaid декодирует JSON в OrderPaidEvent
func DecodeOrderPaid(data []byte) (*model.OrderPaidEvent, error) {
	var event model.OrderPaidEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}