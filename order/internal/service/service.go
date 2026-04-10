package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ChopX4/raketka/order/internal/model"
)

type OrderService interface {
	Cancel(ctx context.Context, orderUUID string) error
	Create(ctx context.Context, order model.OrderRequest) (model.OrderResponse, error)
	Get(ctx context.Context, orderUUID string) (model.OrderByUUID, error)
	Pay(ctx context.Context, req model.PayOrderRequest) (uuid.UUID, error)
	Complete(ctx context.Context, orderUUID string) error
}

type OrderProducer interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaid) error
}

type AssembledConsumer interface {
	RunAssembledConsumer(ctx context.Context) error
}

type OutboxSender interface {
	Run(ctx context.Context) error
}
