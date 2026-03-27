package order

import (
	"github.com/ChopX4/raketka/order/internal/clients/grpc"
	"github.com/ChopX4/raketka/order/internal/repository"
)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(orderRepository repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
