package app

import (
	"context"

	paymentAPI "github.com/ChopX4/raketka/payment/internal/api/payment/v1"
	"github.com/ChopX4/raketka/payment/internal/service"
	paymentService "github.com/ChopX4/raketka/payment/internal/service/method"
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API payment_v1.PaymentServiceServer

	paymentService service.PaymentService
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) payment_v1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentAPI.NewApi(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentService(context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}

	return d.paymentService
}
