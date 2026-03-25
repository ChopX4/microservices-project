package payment

import (
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
)

type paymentClient struct {
	generatedClient payment_v1.PaymentServiceClient
}

func NewPaymentClient(generatedClient payment_v1.PaymentServiceClient) *paymentClient {
	return &paymentClient{
		generatedClient: generatedClient,
	}
}
