package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderUUID string) error {
	if err := s.validateOrderUUID(orderUUID); err != nil {
		return err
	}

	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return err
	}

	if err := s.validateOrderStatusForCancel(order.Status); err != nil {
		return err
	}

	order.Status = model.OrderStatusCanceled

	if err := s.orderRepository.Update(ctx, order); err != nil {
		return err
	}

	return nil
}
