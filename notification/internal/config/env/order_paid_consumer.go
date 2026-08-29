package env

import (
	"os"
)

type OrderPaidConsumerConfig struct {
	topic         string
	consumerGroup string
}

func LoadOrderPaidConsumerConfig() (*OrderPaidConsumerConfig, error) {
	topic := os.Getenv("NOTIFICATION_ORDER_PAID_TOPIC_NAME")
	if topic == "" {
		topic = "order-paid"
	}

	consumerGroup := os.Getenv("NOTIFICATION_ORDER_PAID_CONSUMER_GROUP_ID")
	if consumerGroup == "" {
		consumerGroup = "notification-group-order-paid"
	}

	return &OrderPaidConsumerConfig{
		topic:         topic,
		consumerGroup: consumerGroup,
	}, nil
}

func (c *OrderPaidConsumerConfig) Topic() string         { return c.topic }
func (c *OrderPaidConsumerConfig) ConsumerGroup() string { return c.consumerGroup }