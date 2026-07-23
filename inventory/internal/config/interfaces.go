package config

// InventoryGRPCConfig интерфейс для gRPC конфигурации
type InventoryGRPCConfig interface {
	Address() string
	Host() string
	Port() int
}

// LoggerConfig интерфейс для конфигурации логгера
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// MongoConfig интерфейс для конфигурации MongoDB
type MongoConfig interface {
	URI() string
	Database() string
	Host() string
	Port() int
}