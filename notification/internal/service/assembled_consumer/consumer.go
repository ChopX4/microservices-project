package assembledconsumer

import (
	"context"

	"go.uber.org/zap"

	decoder "github.com/ChopX4/raketka/notification/internal/converter/kafka"
	src "github.com/ChopX4/raketka/notification/internal/service"
	"github.com/ChopX4/raketka/platform/pkg/kafka"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

type service struct {
	assembledConsumer kafka.Consumer
	assembledDecoder  decoder.AssembledDecoder
	telegramService   src.TelegramService
}

func NewAssembledConsumer(assembledConsumer kafka.Consumer, assembledDecoder decoder.AssembledDecoder, telegramService src.TelegramService) *service {
	return &service{
		assembledConsumer: assembledConsumer,
		assembledDecoder:  assembledDecoder,
		telegramService:   telegramService,
	}
}

func (s *service) RunAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting AssembledConsumer service")

	if err := s.assembledConsumer.Consume(ctx, s.AssembledHandler); err != nil {
		logger.Error(ctx, "failed to consume order assembled topic", zap.Error(err))
		return err
	}

	return nil
}
