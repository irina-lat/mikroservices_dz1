package config

import (
	"fmt"

	"assembly/internal/config/env"
)

type Config struct {
	Logger      *env.LoggerConfig
	Kafka       *env.KafkaConfig
	Consumer    *env.OrderPaidConsumerConfig
	Producer    *env.OrderAssembledProducerConfig
}

func Load() (*Config, error) {
	logger, err := env.LoadLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger: %w", err)
	}

	kafka, err := env.LoadKafkaConfig()
	if err != nil {
		return nil, fmt.Errorf("load kafka: %w", err)
	}

	consumer, err := env.LoadOrderPaidConsumerConfig()
	if err != nil {
		return nil, fmt.Errorf("load consumer: %w", err)
	}

	producer, err := env.LoadOrderAssembledProducerConfig()
	if err != nil {
		return nil, fmt.Errorf("load producer: %w", err)
	}

	return &Config{
		Logger:   logger,
		Kafka:    kafka,
		Consumer: consumer,
		Producer: producer,
	}, nil
}