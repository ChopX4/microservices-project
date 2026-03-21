package service

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/google/uuid"
)

type OrderService interface {
	Cancel(ctx context.Context, orderUUID string) error
	Create(ctx context.Context, order model.OrderByUUID) error
	Get(ctx context.Context, orderUUID string) (model.OrderByUUID, error)
	Pay(ctx context.Context, orderUUID, transactionUUID string, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}
