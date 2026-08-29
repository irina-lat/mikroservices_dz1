package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"notification/internal/client/http/telegram"
	"notification/internal/config"
	"notification/internal/service"
	"notification/internal/service/consumer/order_assembled_consumer"
	"notification/internal/service/consumer/order_paid_consumer"
	"platform/pkg/kafka/consumer"
	"platform/pkg/logger"
	middleware "platform/pkg/middleware/kafka"
)

type DI struct {
	Config                 *config.Config
	NotificationService    *service.NotificationService
	OrderPaidConsumer      *order_paid_consumer.OrderPaidConsumer
	OrderAssembledConsumer *order_assembled_consumer.OrderAssembledConsumer
}

func NewDI(cfg *config.Config) (*DI, error) {
	logger.Init(cfg.Logger.Level(), cfg.Logger.AsJSON())
	log := logger.Logger()
	ctx := context.Background()

	di := &DI{Config: cfg}

	// 1. Telegram клиент
	telegramClient := telegram.NewClient(cfg.Telegram.Token(), cfg.Telegram.ChatID())
	di.NotificationService = service.NewNotificationService(telegramClient)

	// 2. Kafka Consumer для OrderPaid
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroupPaid, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers(), cfg.OrderPaid.ConsumerGroup(), saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group for order paid: %w", err)
	}

	kafkaConsumerPaid := consumer.NewConsumer(
		consumerGroupPaid,
		[]string{cfg.OrderPaid.Topic()},
		log,
		middleware.Logging(log),
	)
	di.OrderPaidConsumer = order_paid_consumer.NewOrderPaidConsumer(kafkaConsumerPaid, di.NotificationService)

	// 3. Kafka Consumer для OrderAssembled
	consumerGroupAssembled, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers(), cfg.OrderAssembled.ConsumerGroup(), saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group for order assembled: %w", err)
	}

	kafkaConsumerAssembled := consumer.NewConsumer(
		consumerGroupAssembled,
		[]string{cfg.OrderAssembled.Topic()},
		log,
		middleware.Logging(log),
	)
	di.OrderAssembledConsumer = order_assembled_consumer.NewOrderAssembledConsumer(kafkaConsumerAssembled, di.NotificationService)

	log.Info(ctx, "NotificationService DI initialized",
		zap.Strings("kafka_brokers", cfg.Kafka.Brokers()),
		zap.String("order_paid_topic", cfg.OrderPaid.Topic()),
		zap.String("order_assembled_topic", cfg.OrderAssembled.Topic()),
	)

	return di, nil
}