package env

import (
	"os"
	"strings"
)

type KafkaConfig struct {
	brokers []string
}

func LoadKafkaConfig() (*KafkaConfig, error) {
	brokersStr := os.Getenv("ORDER_KAFKA_BROKERS")
	if brokersStr == "" {
		brokersStr = "localhost:9092"
	}

	brokers := strings.Split(brokersStr, ",")
	for i := range brokers {
		brokers[i] = strings.TrimSpace(brokers[i])
	}

	return &KafkaConfig{
		brokers: brokers,
	}, nil
}

func (c *KafkaConfig) Brokers() []string { return c.brokers }