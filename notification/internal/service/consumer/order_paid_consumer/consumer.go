package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"
	"notification/internal/converter/kafka/decoder"
	"notification/internal/service"
	"platform/pkg/kafka"
	"platform/pkg/kafka/consumer"
	"platform/pkg/logger"
)

type OrderPaidConsumer struct {
	consumer kafka.Consumer
	service  *service.NotificationService
}

func NewOrderPaidConsumer(consumer kafka.Consumer, service *service.NotificationService) *OrderPaidConsumer {
	return &OrderPaidConsumer{
		consumer: consumer,
		service:  service,
	}
}

func (c *OrderPaidConsumer) Start(ctx context.Context) error {
	return c.consumer.Consume(ctx, c.Handle)
}

func (c *OrderPaidConsumer) Handle(ctx context.Context, msg consumer.Message) error {
	log := logger.Logger()

	event, err := decoder.DecodeOrderPaid(msg.Value)
	if err != nil {
		log.Error(ctx, "Failed to decode OrderPaid event", zap.Error(err))
		return nil
	}

	log.Info(ctx, "Received OrderPaid event",
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
	)

	// Отправляем уведомление в Telegram
	if err := c.service.SendOrderPaidNotification(ctx, event); err != nil {
		log.Error(ctx, "Failed to send notification", zap.Error(err))
		return err
	}

	return nil
}