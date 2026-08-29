package order_producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"platform/pkg/kafka"
	"platform/pkg/logger"
)

type OrderProducer struct {
	producer kafka.Producer
	topic    string
}

func NewOrderProducer(producer kafka.Producer, topic string) *OrderProducer {
	return &OrderProducer{
		producer: producer,
		topic:    topic,
	}
}

func (p *OrderProducer) SendShipAssembled(ctx context.Context, orderUUID, userUUID string, buildTimeSec int64) error {
	log := logger.Logger()

	event := map[string]interface{}{
		"event_uuid":     uuid.New().String(),
		"order_uuid":     orderUUID,
		"user_uuid":      userUUID,
		"build_time_sec": buildTimeSec,
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Error(ctx, "Failed to marshal ShipAssembled event", zap.Error(err))
		return err
	}

	log.Info(ctx, "Sending ShipAssembled event",
		zap.String("topic", p.topic),
		zap.String("order_uuid", orderUUID),
	)

	return p.producer.Send(ctx, []byte(orderUUID), data)
}