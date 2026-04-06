package orderconsumer

import (
	"context"
	"time"

	"github.com/ChopX4/raketka/assembly/internal/model"
	"github.com/ChopX4/raketka/platform/pkg/kafka/consumer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	buildTime := 10
	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-time.After(time.Duration(buildTime) * time.Second):
	}

	shipEvent := model.ShipAssembled{
		EventUuid:    uuid.NewString(),
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: int64(buildTime),
	}

	if err := s.orderProducer.ProduceShipAssembled(ctx, shipEvent); err != nil {
		logger.Error(ctx, "Failed to produce ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
