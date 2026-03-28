package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) Pay(ctx context.Context, req model.PayOrderRequest) (uuid.UUID, error) {
	order, err := s.orderRepository.Get(ctx, req.OrderUuid)
	if err != nil {
		return uuid.Nil, err
	}

	stringTransactionUUID, err := s.paymentClient.Pay(ctx, req)
	if err != nil {
		return uuid.Nil, err
	}

	transactionUUID, err := uuid.Parse(stringTransactionUUID)
	if err != nil {
		return uuid.Nil, err
	}

	order.Status = model.OrderStatusPaid
	order.TransactionUUID = transactionUUID
	order.PaymentMethod = req.PaymentMethod

	if err := s.orderRepository.Update(ctx, order); err != nil {
		return uuid.Nil, err
	}

	return transactionUUID, nil
}
