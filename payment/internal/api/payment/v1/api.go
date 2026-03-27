package v1

import (
	"github.com/ChopX4/raketka/payment/internal/service"
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
)

type api struct {
	payment_v1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewApi(service service.PaymentService) *api {
	return &api{
		paymentService: service,
	}
}
