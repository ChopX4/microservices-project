package orderproducer

import (
	"context"

	"github.com/ChopX4/raketka/assembly/internal/converter"
	"github.com/ChopX4/raketka/assembly/internal/model"
	"github.com/ChopX4/raketka/platform/pkg/kafka"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type service struct {
	shipProducer kafka.Producer
}

func NewShipProducer(producer kafka.Producer) *service {
	return &service{
		shipProducer: producer,
	}
}

func (s *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error {
	msg := converter.ShipAssembledToProto(event)

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal ShipAssembled", zap.Error(err))
		return err
	}

	if err := s.shipProducer.Send(ctx, []byte(msg.EventUuid), payload); err != nil {
		logger.Error(ctx, "failed to publish ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
