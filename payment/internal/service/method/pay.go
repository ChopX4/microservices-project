package method

import (
	"context"

	"github.com/google/uuid"

	"github.com/ChopX4/raketka/payment/internal/model"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

func (s *service) Pay(ctx context.Context, req model.PayOrderRequest) (string, error) {
	if err := s.validatePayRequest(req); err != nil {
		return "", err
	}

	transactionUUID := uuid.New()

	logger.Info(ctx, "payment completed successfully")

	return transactionUUID.String(), nil
}
