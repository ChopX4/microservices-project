package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/ChopX4/raketka/order/internal/repository/converter"
)

func (r *repository) Update(_ context.Context, order model.OrderByUUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	repoOrder := converter.OrderByUUIDToRepo(order)

	if _, ok := r.orders[repoOrder.OrderUUID.String()]; !ok {
		return model.ErrNotFound
	}

	r.orders[repoOrder.OrderUUID.String()] = repoOrder

	return nil
}
