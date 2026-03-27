package v1

import (
	"context"
	"errors"

	"github.com/ChopX4/raketka/order/internal/converter"
	"github.com/ChopX4/raketka/order/internal/model"
	order_v1 "github.com/ChopX4/raketka/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderById(ctx context.Context, params order_v1.GetOrderByIdParams) (order_v1.GetOrderByIdRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Order not found",
			}, nil
		}

		return nil, err
	}

	return converter.OrderToHttp(order), nil
}
