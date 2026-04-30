package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
)

func (s *service) Complete(ctx context.Context, orderUUID string) error {
	if err := s.validateOrderUUID(orderUUID); err != nil {
		return err
	}

	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return err
	}

	if err := s.validateOrderStatusForComplete(order.Status); err != nil {
		return err
	}

	order.Status = model.OrderStatusCompleted

	if err := s.orderRepository.Update(ctx, order); err != nil {
		return err
	}

	return nil
}
