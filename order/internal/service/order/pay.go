package order

import (
	"context"
	"fmt"

	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) Pay(ctx context.Context, orderUUID, transactionUUID string, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	order, err := s.orderRepository.Get(orderUUID)
	if err != nil {
		return uuid.Nil, err
	}

	UUIDtransactionUUID, err := uuid.Parse(transactionUUID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse payment transaction id: %w", err)
	}

	order.Status = model.OrderStatusPaid
	order.TransactionUUID = UUIDtransactionUUID
	order.PaymentMethod = paymentMethod

	if err := s.orderRepository.Update(ctx, order); err != nil {
		return uuid.Nil, err
	}

	return UUIDtransactionUUID, nil
}
