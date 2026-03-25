package grpc

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
