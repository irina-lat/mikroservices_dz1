package config

type PaymentGRPCConfig interface {
	Address() string
	Host() string
	Port() int
}

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}