package orderconsumer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/assembly/internal/model"
	"github.com/ChopX4/raketka/platform/pkg/kafka/consumer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

const buildTimeSec = 10

var after = time.After

func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	startAt := time.Now()

	event, err := s.orderDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	if err := s.validateOrderPaidEvent(event); err != nil {
		logger.Error(ctx, "OrderPaid event validation failed", zap.Error(err))
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-after(time.Duration(buildTimeSec) * time.Second):
	}

	shipEvent := model.ShipAssembled{
		EventUuid:    uuid.NewString(),
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: int64(buildTimeSec),
	}

	if err := s.orderProducer.ProduceShipAssembled(ctx, shipEvent); err != nil {
		logger.Error(ctx, "Failed to produce ShipAssembled", zap.Error(err))
		return err
	}

	s.recordAssemblyDuration(ctx, time.Since(startAt).Seconds())

	return nil
}
