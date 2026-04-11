package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
)

func (s *service) Complete(ctx context.Context, orderUUID string) error {
	if !model.IsValidUUID(orderUUID) {
		return model.ErrBadRequest
	}

	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return err
	}

	if order.Status == model.OrderStatusCanceled || order.Status == model.OrderStatusCompleted || order.Status == model.OrderStatusPendingPayment {
		return model.ErrConflict
	}

	order.Status = model.OrderStatusCompleted

	if err := s.orderRepository.Update(ctx, order); err != nil {
		return err
	}

	return nil
}
