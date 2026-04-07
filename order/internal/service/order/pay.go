package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/ChopX4/raketka/order/internal/model"
)

func (s *service) Pay(ctx context.Context, req model.PayOrderRequest) (uuid.UUID, error) {
	if !model.IsValidUUID(req.OrderUuid) {
		return uuid.Nil, model.ErrBadRequest
	}

	if !req.PaymentMethod.IsValid() {
		return uuid.Nil, model.ErrBadRequest
	}

	order, err := s.orderRepository.Get(ctx, req.OrderUuid)
	if err != nil {
		return uuid.Nil, err
	}

	if order.Status == model.OrderStatusCanceled ||
		order.Status == model.OrderStatusPaid ||
		order.Status == model.OrderStatusCompleted {
		return uuid.Nil, model.ErrConflict
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

	message := model.OrderPaid{
		EventUuid:       uuid.NewString(),
		OrderUuid:       order.OrderUUID.String(),
		UserUuid:        order.UserUUID.String(),
		TransactionUuid: transactionUUID.String(),
	}

	if err := s.orderProducer.ProduceOrderPaid(ctx, message); err != nil {
		return uuid.Nil, err
	}

	return transactionUUID, nil
}
