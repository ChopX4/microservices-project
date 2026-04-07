package order

import (
	"github.com/ChopX4/raketka/order/internal/clients/grpc"
	"github.com/ChopX4/raketka/order/internal/repository"
	src "github.com/ChopX4/raketka/order/internal/service"
)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
	orderProducer   src.OrderProducer
}

func NewService(orderRepository repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient, orderProducer src.OrderProducer) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderProducer:   orderProducer,
	}
}
