package service

import (
	"context"

	"github.com/ChopX4/raketka/payment/internal/model"
)

type PaymentService interface {
	Pay(ctx context.Context, req model.PayOrderRequest) (string, error)
}
