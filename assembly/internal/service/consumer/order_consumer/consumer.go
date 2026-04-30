package orderconsumer

import (
	"context"

	"go.uber.org/zap"

	decoder "github.com/ChopX4/raketka/assembly/internal/converter/kafka"
	src "github.com/ChopX4/raketka/assembly/internal/service"
	orderMetrics "github.com/ChopX4/raketka/assembly/internal/service/consumer/order_consumer/metrics"
	"github.com/ChopX4/raketka/platform/pkg/kafka"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

type service struct {
	orderConsumer kafka.Consumer
	orderDecoder  decoder.OrderPaidDecoder
	orderProducer src.ShipProducer
	metrics       *orderMetrics.Metrics
}

func NewOrderConsumer(orderConsumer kafka.Consumer, orderDecoder decoder.OrderPaidDecoder, orderProducer src.ShipProducer, metrics *orderMetrics.Metrics) *service {
	return &service{
		orderConsumer: orderConsumer,
		orderDecoder:  orderDecoder,
		orderProducer: orderProducer,
		metrics:       metrics,
	}
}

func (s *service) recordAssemblyDuration(ctx context.Context, seconds float64) {
	if s.metrics == nil {
		return
	}

	s.metrics.AssemblyDuration.Record(ctx, seconds)
}

func (s *service) RunOrderConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting OrderConsumer service")

	if err := s.orderConsumer.Consume(ctx, s.OrderHandler); err != nil {
		logger.Error(ctx, "failed to consume order paid topic", zap.Error(err))
		return err
	}

	return nil
}
