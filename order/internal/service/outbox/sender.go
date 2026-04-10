package outbox

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/order/internal/repository"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

const (
	defaultBatchSize    = 100
	defaultPollInterval = time.Second
)

type sender struct {
	outboxRepository repository.OutboxRepository
	producer         sarama.SyncProducer
}

func NewSender(outboxRepository repository.OutboxRepository, producer sarama.SyncProducer) *sender {
	return &sender{
		outboxRepository: outboxRepository,
		producer:         producer,
	}
}

func (s *sender) Run(ctx context.Context) error {
	ticker := time.NewTicker(defaultPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.flush(ctx); err != nil {
				logger.Error(ctx, "failed to flush outbox", zap.Error(err))
			}
		}
	}
}

func (s *sender) flush(ctx context.Context) error {
	events, err := s.outboxRepository.ListPending(ctx, defaultBatchSize)
	if err != nil {
		return err
	}

	for _, event := range events {
		msg := &sarama.ProducerMessage{
			Topic: event.Topic,
			Key:   sarama.ByteEncoder([]byte(event.Key)),
			Value: sarama.ByteEncoder(event.Payload),
		}

		if _, _, err := s.producer.SendMessage(msg); err != nil {
			logger.Error(ctx, "failed to publish outbox event", zap.String("event_uuid", event.EventUUID), zap.Error(err))
			continue
		}

		if err := s.outboxRepository.MarkPublished(ctx, event.EventUUID); err != nil {
			logger.Error(ctx, "failed to mark outbox event published", zap.String("event_uuid", event.EventUUID), zap.Error(err))
		}
	}

	return nil
}
