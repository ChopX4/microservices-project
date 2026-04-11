package orderproducer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/ChopX4/raketka/order/internal/converter"
	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/ChopX4/raketka/platform/pkg/kafka"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

type service struct {
	orderProducer kafka.Producer
}

func NewOrderProducer(orderProducer kafka.Producer) *service {
	return &service{
		orderProducer: orderProducer,
	}
}

func (s *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaid) error {
	msg := converter.OrderPaidToProto(event)

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderPaid", zap.Error(err))
		return err
	}

	if err := s.orderProducer.Send(ctx, []byte(msg.EventUuid), payload); err != nil {
		logger.Error(ctx, "failed to publish OrderPaid", zap.Error(err))
		return err
	}

	return nil
}
