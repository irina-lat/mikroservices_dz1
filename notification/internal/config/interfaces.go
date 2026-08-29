package config

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidConsumerConfig interface {
	Topic() string
	ConsumerGroup() string
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	ConsumerGroup() string
}

type TelegramBotConfig interface {
	Token() string
	ChatID() string
}