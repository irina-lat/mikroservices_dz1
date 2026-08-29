package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"
	"notification/internal/converter/kafka/decoder"
	"notification/internal/service"
	"platform/pkg/kafka"
	"platform/pkg/kafka/consumer"
	"platform/pkg/logger"
)

type OrderAssembledConsumer struct {
	consumer kafka.Consumer
	service  *service.NotificationService
}

func NewOrderAssembledConsumer(consumer kafka.Consumer, service *service.NotificationService) *OrderAssembledConsumer {
	return &OrderAssembledConsumer{
		consumer: consumer,
		service:  service,
	}
}

func (c *OrderAssembledConsumer) Start(ctx context.Context) error {
	return c.consumer.Consume(ctx, c.Handle)
}

func (c *OrderAssembledConsumer) Handle(ctx context.Context, msg consumer.Message) error {
	log := logger.Logger()

	event, err := decoder.DecodeOrderAssembled(msg.Value)
	if err != nil {
		log.Error(ctx, "Failed to decode OrderAssembled event", zap.Error(err))
		return nil
	}

	log.Info(ctx, "Received OrderAssembled event",
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
	)

	// Отправляем уведомление в Telegram
	if err := c.service.SendOrderAssembledNotification(ctx, event); err != nil {
		log.Error(ctx, "Failed to send notification", zap.Error(err))
		return err
	}

	return nil
}