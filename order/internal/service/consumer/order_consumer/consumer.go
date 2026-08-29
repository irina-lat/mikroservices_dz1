package order_consumer

import (
	"context"

	"go.uber.org/zap"
	"order/internal/converter/kafka/decoder"
	"order/internal/model"
	"order/internal/service/order"
	"platform/pkg/kafka"
	"platform/pkg/kafka/consumer"
	"platform/pkg/logger"
)

type OrderAssembledConsumer struct {
	consumer     kafka.Consumer
	orderService order.Service
}

func NewOrderAssembledConsumer(consumer kafka.Consumer, orderService order.Service) *OrderAssembledConsumer {
	return &OrderAssembledConsumer{
		consumer:     consumer,
		orderService: orderService,
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

	// Обновляем статус заказа на ASSEMBLED
	if err := c.orderService.UpdateOrderStatus(ctx, event.OrderUUID, model.StatusAssembled); err != nil {
		log.Error(ctx, "Failed to update order status to ASSEMBLED", zap.Error(err))
		return err
	}

	log.Info(ctx, "Order status updated to ASSEMBLED",
		zap.String("order_uuid", event.OrderUUID),
	)

	return nil
}