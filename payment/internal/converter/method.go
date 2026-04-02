package converter

import (
	"github.com/ChopX4/raketka/payment/internal/model"
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
)

func PayOrderRequestToModel(payment *payment_v1.PayOrderRequest) model.PayOrderRequest {
	return model.PayOrderRequest{
		OrderUuid:     payment.GetOrderUuid(),
		UserUuid:      payment.GetUserUuid(),
		PaymentMethod: PaymentMethodToModel(payment.PaymentMethod),
	}
}

func PaymentMethodToModel(method payment_v1.PaymentMethod) model.PaymentMethod {
	paymentMethod := model.PaymentMethod(method)
	if !paymentMethod.IsValid() {
		return model.PaymentMethodUnknown
	}

	return paymentMethod
}
