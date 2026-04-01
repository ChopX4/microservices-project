package payment

import (
	"context"

	"go.uber.org/zap"

	"github.com/ChopX4/raketka/order/internal/clients/converter"
	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

func (c *paymentClient) Pay(ctx context.Context, req model.PayOrderRequest) (string, error) {
	resp, err := c.generatedClient.PayOrder(ctx, converter.PayOrderRequestToProto(req))
	if err != nil {
		logger.Error(ctx, "failed to pay order via payment grpc", zap.String("order_uuid", req.OrderUuid), zap.Error(err))
		return "", err
	}

	return resp.GetTransactionUuid(), nil
}
