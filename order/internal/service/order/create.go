package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
)

func (s *service) Create(ctx context.Context, order model.OrderByUUID) error {
	order.Status = model.OrderStatusPendingPayment

	if err := s.orderRepository.Create(ctx, order); err != nil {
		return err
	}

	return nil
}
