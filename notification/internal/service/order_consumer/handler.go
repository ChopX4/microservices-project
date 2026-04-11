package orderconsumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/ChopX4/raketka/platform/pkg/kafka/consumer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	if err := s.telegramService.SendOrderNotification(ctx, event); err != nil {
		logger.Error(ctx, "Failed to send Order notification", zap.Error(err))
		return err
	}

	return nil
}
