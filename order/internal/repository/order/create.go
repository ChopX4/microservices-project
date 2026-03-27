package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/ChopX4/raketka/order/internal/repository/converter"
)

func (r *repository) Create(_ context.Context, order model.OrderByUUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.OrderUUID.String()] = converter.OrderByUUIDToRepo(order)

	return nil
}
