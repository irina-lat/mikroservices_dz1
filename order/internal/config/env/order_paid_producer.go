package env

import (
	"os"
)

type OrderPaidProducerConfig struct {
	topic string
}

func LoadOrderPaidProducerConfig() (*OrderPaidProducerConfig, error) {
	topic := os.Getenv("ORDER_ORDER_PAID_TOPIC_NAME")
	if topic == "" {
		topic = "order.paid"
	}

	return &OrderPaidProducerConfig{
		topic: topic,
	}, nil
}

func (c *OrderPaidProducerConfig) Topic() string { return c.topic }