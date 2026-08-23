package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"assembly/internal/config"
	"assembly/internal/service/consumer/order_consumer"
	"assembly/internal/service/producer/order_producer"
	"platform/pkg/kafka/consumer"
	"platform/pkg/kafka/producer"
	"platform/pkg/logger"
	middleware "platform/pkg/middleware/kafka"
)

type DI struct {
	Config        *config.Config
	OrderConsumer *order_consumer.OrderConsumer
	OrderProducer *order_producer.OrderProducer
}

func NewDI(cfg *config.Config) (*DI, error) {
	// logger.Init(level string, asJSON bool)
	logger.Init(cfg.Logger.Level(), cfg.Logger.AsJSON())
	log := logger.Logger()
	ctx := context.Background()

	di := &DI{Config: cfg}

	// 1. Kafka Producer
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	syncProducer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers(), saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create sync producer: %w", err)
	}

	kafkaProducer := producer.NewProducer(syncProducer, cfg.Producer.Topic(), log)
	di.OrderProducer = order_producer.NewOrderProducer(kafkaProducer, cfg.Producer.Topic())

	// 2. Kafka Consumer
	consumerGroup, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers(), cfg.Consumer.ConsumerGroup(), saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	kafkaConsumer := consumer.NewConsumer(
		consumerGroup,
		[]string{cfg.Consumer.Topic()},
		log,
		middleware.Logging(log),
	)

	di.OrderConsumer = order_consumer.NewOrderConsumer(kafkaConsumer, di.OrderProducer)

	log.Info(ctx, "AssemblyService DI initialized",
		zap.Strings("brokers", cfg.Kafka.Brokers()),
		zap.String("consumer_topic", cfg.Consumer.Topic()),
		zap.String("producer_topic", cfg.Producer.Topic()),
	)

	return di, nil
}