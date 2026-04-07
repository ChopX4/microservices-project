package orderconsumer

import (
	"context"

	"go.uber.org/zap"

	decoder "github.com/ChopX4/raketka/assembly/internal/converter/kafka"
	src "github.com/ChopX4/raketka/assembly/internal/service"
	"github.com/ChopX4/raketka/platform/pkg/kafka"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

type service struct {
	orderConsumer kafka.Consumer
	orderDecoder  decoder.OrderPaidDecoder
	orderProducer src.ShipProducer
}

func NewOrderConsumer(orderConsumer kafka.Consumer, orderDecoder decoder.OrderPaidDecoder, orderProducer src.ShipProducer) *service {
	return &service{
		orderConsumer: orderConsumer,
		orderDecoder:  orderDecoder,
		orderProducer: orderProducer,
	}
}

func (s *service) RunOrderConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting OrderConsumer service")

	if err := s.orderConsumer.Consume(ctx, s.OrderHandler); err != nil {
		logger.Error(ctx, "failed to consume order paid topic", zap.Error(err))
		return err
	}

	return nil
}
