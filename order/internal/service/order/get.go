package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID string) (model.OrderByUUID, error) {
	if err := s.validateOrderUUID(orderUUID); err != nil {
		return model.OrderByUUID{}, err
	}

	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return model.OrderByUUID{}, err
	}

	return order, nil
}
