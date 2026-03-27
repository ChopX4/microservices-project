package v1

import (
	"context"
	"net/http"

	order_v1 "github.com/ChopX4/raketka/shared/pkg/openapi/order/v1"
)

func (a *api) NewError(ctx context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}
