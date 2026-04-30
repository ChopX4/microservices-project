package assembledconsumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/ChopX4/raketka/platform/pkg/kafka/consumer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

func (s *service) AssembledHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.assembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled", zap.Error(err))
		return err
	}

	if err := s.validateShipAssembledEvent(event); err != nil {
		logger.Error(ctx, "ShipAssembled event validation failed", zap.Error(err))
		return err
	}

	if err := s.telegramService.SendShipNotification(ctx, event); err != nil {
		logger.Error(ctx, "Failed to send Ship notification", zap.Error(err))
		return err
	}

	return nil
}
