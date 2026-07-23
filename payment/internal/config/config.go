package config

import (
	"fmt"

	"payment/internal/config/env"
)

type Config struct {
	Logger  *env.LoggerConfig
	Payment *env.PaymentGRPCConfig
}

func Load() (*Config, error) {
	logger, err := env.LoadLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger: %w", err)
	}

	payment, err := env.LoadPaymentGRPCConfig()
	if err != nil {
		return nil, fmt.Errorf("load payment grpc: %w", err)
	}

	return &Config{
		Logger:  logger,
		Payment: payment,
	}, nil
}