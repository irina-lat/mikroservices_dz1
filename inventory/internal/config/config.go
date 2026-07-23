package config

import (
	"fmt"

	"inventory/internal/config/env"
)

// Config объединяет все конфигурации сервиса
type Config struct {
	GRPC    *env.InventoryGRPCConfig
	Logger  *env.LoggerConfig
	Mongo   *env.MongoConfig
}

// Load загружает все конфигурации из переменных окружения
func Load() (*Config, error) {
	grpcConfig, err := env.LoadInventoryGRPCConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load gRPC config: %w", err)
	}

	loggerConfig, err := env.LoadLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load logger config: %w", err)
	}

	mongoConfig, err := env.LoadMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load mongo config: %w", err)
	}

	return &Config{
		GRPC:   grpcConfig,
		Logger: loggerConfig,
		Mongo:  mongoConfig,
	}, nil
}