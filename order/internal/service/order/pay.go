package order

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/ChopX4/raketka/order/internal/converter"
	"github.com/ChopX4/raketka/order/internal/model"
	repoModel "github.com/ChopX4/raketka/order/internal/repository/model"
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

	req.UserUuid = order.UserUUID.String()

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

	message := model.OrderPaid{
		EventUuid:       uuid.NewString(),
		OrderUuid:       order.OrderUUID.String(),
		UserUuid:        order.UserUUID.String(),
		TransactionUuid: transactionUUID.String(),
	}

	payload, err := proto.Marshal(converter.OrderPaidToProto(message))
	if err != nil {
		return uuid.Nil, err
	}

	outboxMsg := repoModel.OutboxMessage{
		EventUUID: message.EventUuid,
		Topic:     s.orderPaidTopic,
		Key:       message.OrderUuid,
		Payload:   payload,
		Status:    repoModel.OutboxStatusPending,
	}

	if err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
		if err := s.orderRepository.Update(ctx, order); err != nil {
			return err
		}

		return s.outboxRepository.Create(ctx, outboxMsg)
	}); err != nil {
		return uuid.Nil, err
	}
	s.addOrdersRevenue(ctx, float64(order.TotalPrice))

	return transactionUUID, nil
}
