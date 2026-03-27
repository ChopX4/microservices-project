package repository

import (
	"context"

	"github.com/ChopX4/raketka/inventory/internal/model"
)

type InventoryRepository interface {
	Get(ctx context.Context, uuid string) (model.Part, error)
	List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
