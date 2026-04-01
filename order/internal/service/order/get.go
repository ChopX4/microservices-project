package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID string) (model.OrderByUUID, error) {
	if !model.IsValidUUID(orderUUID) {
		return model.OrderByUUID{}, model.ErrBadRequest
	}

	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return model.OrderByUUID{}, err
	}

	return order, nil
}
