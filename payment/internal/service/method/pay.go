package method

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/ChopX4/raketka/payment/internal/model"
)

func (s *service) Pay(ctx context.Context, req model.PayOrderRequest) (string, error) {
	transaction_uuid := uuid.New()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transaction_uuid.String())

	return transaction_uuid.String(), nil
}
