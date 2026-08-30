package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type IAMGRPCConfig interface {
	Address() string
	Host() string
	Port() int
}

type PostgresConfig interface {
	DSN() string
	Host() string
	Port() int
	User() string
	Password() string
	Database() string
	SSLMode() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type SessionConfig interface {
	TTL() time.Duration
}