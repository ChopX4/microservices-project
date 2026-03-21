package order

import (
	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/ChopX4/raketka/order/internal/repository/converter"
)

func (r *repository) Get(orderUUID string) (model.OrderByUUID, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return model.OrderByUUID{}, model.ErrNotFound
	}

	return converter.OrderByUUIDToModel(order), nil
}
