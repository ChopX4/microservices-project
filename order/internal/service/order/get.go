package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID string) (model.OrderByUUID, error) {
	order, err := s.orderRepository.Get(orderUUID)
	if err != nil {
		return model.OrderByUUID{}, err
	}

	return order, nil
}
