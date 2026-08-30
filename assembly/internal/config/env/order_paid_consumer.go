package env

import (
	"os"
)

type OrderPaidConsumerConfig struct {
	topic         string
	consumerGroup string
}

func LoadOrderPaidConsumerConfig() (*OrderPaidConsumerConfig, error) {
	topic := os.Getenv("ASSEMBLY_KAFKA_ORDER_PAID_TOPIC")
	if topic == "" {
		topic = "order.paid"
	}

	consumerGroup := os.Getenv("ASSEMBLY_KAFKA_CONSUMER_GROUP")
	if consumerGroup == "" {
		consumerGroup = "assembly-group"
	}

	return &OrderPaidConsumerConfig{
		topic:         topic,
		consumerGroup: consumerGroup,
	}, nil
}

func (c *OrderPaidConsumerConfig) Topic() string         { return c.topic }
func (c *OrderPaidConsumerConfig) ConsumerGroup() string { return c.consumerGroup }