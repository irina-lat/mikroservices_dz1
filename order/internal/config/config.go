package config

import (
	"fmt"

	"order/internal/config/env"
)

type Config struct {
	Logger    *env.LoggerConfig
	HTTP      *env.OrderHTTPConfig
	Postgres  *env.PostgresConfig
	Inventory *env.InventoryGRPCConfig
	Payment   *env.PaymentGRPCConfig
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

	return &Config{
		Logger:    logger,
		HTTP:      http,
		Postgres:  pg,
		Inventory: inventory,
		Payment:   payment,
	}, nil
}
