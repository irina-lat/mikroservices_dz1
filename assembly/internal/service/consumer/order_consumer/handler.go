package order_consumer

import (
	"context"

	"platform/pkg/kafka/consumer"
)

type MessageHandler interface {
	Handle(ctx context.Context, msg consumer.Message) error
}