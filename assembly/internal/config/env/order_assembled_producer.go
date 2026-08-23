package env

import (
	"os"
)

type OrderAssembledProducerConfig struct {
	topic string
}

func LoadOrderAssembledProducerConfig() (*OrderAssembledProducerConfig, error) {
	topic := os.Getenv("ASSEMBLY_ORDER_ASSEMBLED_TOPIC_NAME")
	if topic == "" {
		topic = "order.assembled"
	}

	return &OrderAssembledProducerConfig{
		topic: topic,
	}, nil
}

func (c *OrderAssembledProducerConfig) Topic() string { return c.topic }