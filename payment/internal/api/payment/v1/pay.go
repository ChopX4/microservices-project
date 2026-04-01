package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ChopX4/raketka/payment/internal/converter"
	"github.com/ChopX4/raketka/payment/internal/model"
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	transactionUuid, err := a.paymentService.Pay(ctx, converter.PayOrderRequestToModel(req))
	if err != nil {
		if errors.Is(err, model.ErrBadRequest) {
			return nil, status.Error(codes.InvalidArgument, "validation error")
		}

		return nil, err
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}
