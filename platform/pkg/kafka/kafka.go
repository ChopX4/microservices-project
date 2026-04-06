package kafka

import (
	"context"

	"github.com/ChopX4/raketka/platform/pkg/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, handler consumer.MessageHandler) error
}

type Producer interface {
	Send(ctx context.Context, key, value []byte) error
}
