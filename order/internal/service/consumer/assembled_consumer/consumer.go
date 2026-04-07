package assembledconsumer

import (
	"context"

	"go.uber.org/zap"

	decoder "github.com/ChopX4/raketka/order/internal/converter/kafka"
	src "github.com/ChopX4/raketka/order/internal/service"
	"github.com/ChopX4/raketka/platform/pkg/kafka"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

type service struct {
	assembledConsumer kafka.Consumer
	assembledDecoder  decoder.ShipAssembledDecoder
	orderService      src.OrderService
}

func NewAssembledConsumer(assembledConsumer kafka.Consumer, assebmledDecoder decoder.ShipAssembledDecoder, orderService src.OrderService) *service {
	return &service{
		assembledConsumer: assembledConsumer,
		assembledDecoder:  assebmledDecoder,
		orderService:      orderService,
	}
}

func (s *service) RunAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting AssembledConsumer service")

	if err := s.assembledConsumer.Consume(ctx, s.AssembledHandler); err != nil {
		logger.Error(ctx, "failed to consume ship assembled topic", zap.Error(err))
		return err
	}

	return nil
}
