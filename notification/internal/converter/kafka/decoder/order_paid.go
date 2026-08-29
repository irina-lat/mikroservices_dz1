package decoder

import (
	"encoding/json"

	"notification/internal/model"
)

func DecodeOrderPaid(data []byte) (*model.OrderPaidEvent, error) {
	var event model.OrderPaidEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}