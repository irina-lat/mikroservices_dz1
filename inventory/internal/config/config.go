package config

import (
	"fmt"

	"inventory/internal/config/env"
)

type Config struct {
	Logger      *env.LoggerConfig
	GRPC        *env.InventoryGRPCConfig
	IAM         *env.IAMGRPCConfig
	Mongo       *env.MongoConfig
}

func Load() (*Config, error) {
	logger, err := env.LoadLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger: %w", err)
	}

	grpc, err := env.LoadInventoryGRPCConfig()
	if err != nil {
		return nil, fmt.Errorf("load grpc: %w", err)
	}

	iam, err := env.LoadIAMGRPCConfig()
	if err != nil {
		return nil, fmt.Errorf("load iam grpc: %w", err)
	}

	mongo, err := env.LoadMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("load mongo: %w", err)
	}

	return &Config{
		Logger: logger,
		GRPC:   grpc,
		IAM:    iam,
		Mongo:  mongo,
	}, nil
}