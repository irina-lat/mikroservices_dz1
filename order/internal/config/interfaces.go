package config

import "time"

type InventoryGRPCConfig interface {
	Address() string
	Host() string
	Port() int
}

type PaymentGRPCConfig interface {
	Address() string
	Host() string
	Port() int
}

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type OrderHTTPConfig interface {
	Address() string
	Host() string
	Port() int
	ReadTimeout() time.Duration
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
