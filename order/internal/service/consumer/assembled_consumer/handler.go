package assembledconsumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/ChopX4/raketka/platform/pkg/kafka/consumer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

func (s *service) AssembledHandler(ctx context.Context, msg consumer.Message) error {
	message, err := s.assembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled", zap.Error(err))
		return err
	}

	if err := s.orderService.Complete(ctx, message.OrderUuid); err != nil {
		logger.Error(ctx, "Failed to complete order", zap.Error(err))
		return err
	}

	return nil
}
