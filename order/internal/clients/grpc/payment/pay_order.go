package payment

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/clients/converter"
	"github.com/ChopX4/raketka/order/internal/model"
)

func (c *paymentClient) Pay(ctx context.Context, req model.PayOrderRequest) (string, error) {
	resp, err := c.generatedClient.PayOrder(ctx, converter.PayOrderRequestToProto(req))
	if err != nil {
		return "", err
	}

	return resp.GetTransactionUuid(), nil
}
