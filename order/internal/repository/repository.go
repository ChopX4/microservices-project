package repository

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.OrderByUUID) error
	Get(ctx context.Context, orderUUID string) (model.OrderByUUID, error)
	Update(ctx context.Context, order model.OrderByUUID) error
}
