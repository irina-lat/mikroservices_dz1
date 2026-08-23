package order_consumer

import (
	"context"
	"math/rand"
	"time"

	"go.uber.org/zap"
	"assembly/internal/converter/kafka/decoder"
	"assembly/internal/service/producer/order_producer"
	"platform/pkg/kafka"
	"platform/pkg/kafka/consumer"
	"platform/pkg/logger"
)

type OrderConsumer struct {
	consumer kafka.Consumer
	producer *order_producer.OrderProducer
}

func NewOrderConsumer(cons kafka.Consumer, producer *order_producer.OrderProducer) *OrderConsumer {
	return &OrderConsumer{
		consumer: cons,
		producer: producer,
	}
}

func (c *OrderConsumer) Start(ctx context.Context) error {
	return c.consumer.Consume(ctx, c.Handle)
}

func (c *OrderConsumer) Handle(ctx context.Context, msg consumer.Message) error {
	log := logger.Logger()

	event, err := decoder.DecodeOrderPaid(msg.Value)
	if err != nil {
		log.Error(ctx, "Failed to decode OrderPaid event", zap.Error(err))
		return nil
	}

	log.Info(ctx, "Received OrderPaid event",
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
	)

	buildTime := rand.Intn(10) + 1
	log.Info(ctx, "Starting assembly...", zap.Int64("build_time_sec", int64(buildTime)))
	time.Sleep(time.Duration(buildTime) * time.Second)

	if err := c.producer.SendShipAssembled(ctx, event.OrderUUID, event.UserUUID, int64(buildTime)); err != nil {
		log.Error(ctx, "Failed to send ShipAssembled", zap.Error(err))
		return err
	}

	log.Info(ctx, "Assembly completed and ShipAssembled sent",
		zap.String("order_uuid", event.OrderUUID),
		zap.Int64("build_time_sec", int64(buildTime)),
	)

	return nil
}