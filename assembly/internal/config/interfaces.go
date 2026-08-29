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

type OrderAssembledProducerConfig interface {
	Topic() string
}