package v1

import (
	"github.com/ChopX4/raketka/order/internal/service"
)

type api struct {
	orderService service.OrderService
}

func NewApi(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
