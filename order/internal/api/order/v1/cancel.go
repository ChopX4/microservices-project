package v1

import (
	"context"
	"errors"

	"github.com/ChopX4/raketka/order/internal/model"
	order_v1 "github.com/ChopX4/raketka/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	if err := a.orderService.Cancel(ctx, params.OrderUUID.String()); err != nil {
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

	return &order_v1.NoContentError{
		Code:    204,
		Message: "Order successfully canceled",
	}, nil
}
