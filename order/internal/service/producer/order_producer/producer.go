package order_producer

import (
	"context"
	"encoding/json"
	"time"

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

func (p *OrderProducer) SendOrderPaid(ctx context.Context, orderUUID, userUUID, paymentMethod, transactionUUID string) error {
	log := logger.Logger()

	event := map[string]interface{}{
		"event_uuid":       uuid.New().String(),
		"order_uuid":       orderUUID,
		"user_uuid":        userUUID,
		"payment_method":   paymentMethod,
		"transaction_uuid": transactionUUID,
		"paid_at":          time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Error(ctx, "Failed to marshal OrderPaid event", zap.Error(err))
		return err
	}

	log.Info(ctx, "Sending OrderPaid event",
		zap.String("topic", p.topic),
		zap.String("order_uuid", orderUUID),
	)

	return p.producer.Send(ctx, []byte(orderUUID), data)
}