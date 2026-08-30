package config

import (
	"fmt"

	"iam/internal/config/env"
)

type Config struct {
	Logger   *env.LoggerConfig
	GRPC     *env.IAMGRPCConfig
	Postgres *env.PostgresConfig
	Redis    *env.RedisConfig
	Session  *env.SessionConfig
}

func Load() (*Config, error) {
	logger, err := env.LoadLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger: %w", err)
	}

	grpc, err := env.LoadIAMGRPCConfig()
	if err != nil {
		return nil, fmt.Errorf("load grpc: %w", err)
	}

	postgres, err := env.LoadPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("load postgres: %w", err)
	}

	redis, err := env.LoadRedisConfig()
	if err != nil {
		return nil, fmt.Errorf("load redis: %w", err)
	}

	session, err := env.LoadSessionConfig()
	if err != nil {
		return nil, fmt.Errorf("load session: %w", err)
	}

	return &Config{
		Logger:   logger,
		GRPC:     grpc,
		Postgres: postgres,
		Redis:    redis,
		Session:  session,
	}, nil
}