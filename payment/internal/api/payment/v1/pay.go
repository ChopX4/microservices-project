package v1

import (
	"context"

	"github.com/ChopX4/raketka/payment/internal/converter"
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	transactionUuid, err := a.paymentService.Pay(ctx, converter.PayOrderRequestToModel(req))
	if err != nil {
		return nil, err
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}
