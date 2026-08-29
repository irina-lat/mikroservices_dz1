package config

import (
	"fmt"

	"order/internal/config/env"
)

type Config struct {
	Logger              *env.LoggerConfig
	HTTP                *env.OrderHTTPConfig
	Postgres            *env.PostgresConfig
	Inventory           *env.InventoryGRPCConfig
	Payment             *env.PaymentGRPCConfig
	Kafka               *env.KafkaConfig
	OrderPaidProducer   *env.OrderPaidProducerConfig
	OrderAssembledConsumer *env.OrderAssembledConsumerConfig
}

func Load() (*Config, error) {
	logger, err := env.LoadLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger: %w", err)
	}

	http, err := env.LoadOrderHTTPConfig()
	if err != nil {
		return nil, fmt.Errorf("load http: %w", err)
	}

	pg, err := env.LoadPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("load postgres: %w", err)
	}

	inventory, err := env.LoadInventoryGRPCConfig()
	if err != nil {
		return nil, fmt.Errorf("load inventory grpc: %w", err)
	}

	payment, err := env.LoadPaymentGRPCConfig()
	if err != nil {
		return nil, fmt.Errorf("load payment grpc: %w", err)
	}

	kafka, err := env.LoadKafkaConfig()
	if err != nil {
		return nil, fmt.Errorf("load kafka: %w", err)
	}

	producer, err := env.LoadOrderPaidProducerConfig()
	if err != nil {
		return nil, fmt.Errorf("load order paid producer: %w", err)
	}

	consumer, err := env.LoadOrderAssembledConsumerConfig()
	if err != nil {
		return nil, fmt.Errorf("load order assembled consumer: %w", err)
	}

	return &Config{
		Logger:              logger,
		HTTP:                http,
		Postgres:            pg,
		Inventory:           inventory,
		Payment:             payment,
		Kafka:               kafka,
		OrderPaidProducer:   producer,
		OrderAssembledConsumer: consumer,
	}, nil
}