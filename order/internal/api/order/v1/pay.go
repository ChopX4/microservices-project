package v1

import (
	"context"
	"errors"

	"github.com/ChopX4/raketka/order/internal/converter"
	"github.com/ChopX4/raketka/order/internal/model"
	order_v1 "github.com/ChopX4/raketka/shared/pkg/openapi/order/v1"
)

func (a *api) OrderPay(ctx context.Context, req *order_v1.OrderPayReq, params order_v1.OrderPayParams) (order_v1.OrderPayRes, error) {
	if req == nil {
		return &order_v1.BadRequestError{
			Code:    400,
			Message: "Bad request - request is required",
		}, nil
	}

	modelReq := model.PayOrderRequest{
		OrderUuid:     params.OrderUUID.String(),
		PaymentMethod: converter.PaymentMethodToModel(req.PaymentMethod.Value),
	}

	transactionUUID, err := a.orderService.Pay(ctx, modelReq)
	if err != nil {
		if errors.Is(err, model.ErrBadRequest) {
			return &order_v1.BadRequestError{
				Code:    400,
				Message: "Bad request - validation error",
			}, nil
		}
		if errors.Is(err, model.ErrNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Order not found",
			}, nil
		}
		if errors.Is(err, model.ErrConflict) {
			return &order_v1.ConflictError{
				Code:    409,
				Message: "Conflict",
			}, nil
		}
		return nil, err
	}

	return &order_v1.TransactionUUID{
		TransactionUUID: transactionUUID,
	}, nil
}
