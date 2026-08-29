package config

import (
	"fmt"

	"notification/internal/config/env"
)

type Config struct {
	Logger          *env.LoggerConfig
	Kafka           *env.KafkaConfig
	OrderPaid       *env.OrderPaidConsumerConfig
	OrderAssembled  *env.OrderAssembledConsumerConfig
	Telegram        *env.TelegramBotConfig
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

	orderPaid, err := env.LoadOrderPaidConsumerConfig()
	if err != nil {
		return nil, fmt.Errorf("load order paid consumer: %w", err)
	}

	orderAssembled, err := env.LoadOrderAssembledConsumerConfig()
	if err != nil {
		return nil, fmt.Errorf("load order assembled consumer: %w", err)
	}

	telegram, err := env.LoadTelegramBotConfig()
	if err != nil {
		return nil, fmt.Errorf("load telegram bot: %w", err)
	}

	return &Config{
		Logger:         logger,
		Kafka:          kafka,
		OrderPaid:      orderPaid,
		OrderAssembled: orderAssembled,
		Telegram:       telegram,
	}, nil
}