package service

import (
	"assembly/internal/service/consumer/order_consumer"
	"assembly/internal/service/producer/order_producer"
)

type OrderConsumer = order_consumer.OrderConsumer
type OrderProducer = order_producer.OrderProducer