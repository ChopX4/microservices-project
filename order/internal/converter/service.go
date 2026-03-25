package converter

import (
	"github.com/ChopX4/raketka/order/internal/model"
	order_v1 "github.com/ChopX4/raketka/shared/pkg/openapi/order/v1"
)

func OrderRequestToModel(req *order_v1.CreateOrderRequest) model.OrderRequest {
	return model.OrderRequest{
		UserUUID:  req.GetUserUUID(),
		PartUUIDs: req.GetPartUuids(),
	}
}

func PaymentMethodToModel(method order_v1.PaymentMethod) model.PaymentMethod {
	return model.PaymentMethod(method)
}

func PaymentMethodToHttp(method model.PaymentMethod) order_v1.PaymentMethod {
	return order_v1.PaymentMethod(method)
}

func StatusToHttp(status model.OrderStatus) order_v1.OrderStatus {
	return order_v1.OrderStatus(status)
}

func OrderToHttp(order model.OrderByUUID) *order_v1.OrderByUUID {
	return &order_v1.OrderByUUID{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   PaymentMethodToHttp(order.PaymentMethod),
		Status:          StatusToHttp(order.Status),
	}
}
