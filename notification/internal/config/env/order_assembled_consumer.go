package env

import (
	"os"
)

type OrderAssembledConsumerConfig struct {
	topic         string
	consumerGroup string
}

func LoadOrderAssembledConsumerConfig() (*OrderAssembledConsumerConfig, error) {
	topic := os.Getenv("NOTIFICATION_ORDER_ASSEMBLED_TOPIC_NAME")
	if topic == "" {
		topic = "order.assembled"
	}

	consumerGroup := os.Getenv("NOTIFICATION_ORDER_ASSEMBLED_CONSUMER_GROUP_ID")
	if consumerGroup == "" {
		consumerGroup = "notification-group-order-assembled"
	}

	return &OrderAssembledConsumerConfig{
		topic:         topic,
		consumerGroup: consumerGroup,
	}, nil
}

func (c *OrderAssembledConsumerConfig) Topic() string         { return c.topic }
func (c *OrderAssembledConsumerConfig) ConsumerGroup() string { return c.consumerGroup }