package v1

import (
	"context"
	"errors"

	"github.com/ChopX4/raketka/order/internal/converter"
	"github.com/ChopX4/raketka/order/internal/model"
	order_v1 "github.com/ChopX4/raketka/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	if req == nil {
		return &order_v1.BadRequestError{
			Code:    400,
			Message: "Bad request - request is required",
		}, nil
	}

	modelReq, err := a.orderService.Create(ctx, converter.OrderRequestToModel(req))
	if err != nil {
		if errors.Is(err, model.ErrBadRequest) {
			return &order_v1.BadRequestError{
				Code:    400,
				Message: "Bad request - validation error",
			}, nil
		}

		return nil, err
	}

	return &order_v1.CreateOrderResponse{
		OrderUUID:  modelReq.OrderUUID,
		TotalPrice: modelReq.TotalPrice,
	}, nil
}
