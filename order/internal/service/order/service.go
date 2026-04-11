package order

import (
	"github.com/ChopX4/raketka/order/internal/clients/grpc"
	"github.com/ChopX4/raketka/order/internal/repository"
)

type service struct {
	orderRepository  repository.OrderRepository
	outboxRepository repository.OutboxRepository
	txManager        repository.TxManager
	inventoryClient  grpc.InventoryClient
	paymentClient    grpc.PaymentClient
	orderPaidTopic   string
}

func NewService(orderRepository repository.OrderRepository, outboxRepository repository.OutboxRepository, txManager repository.TxManager, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient, orderPaidTopic string) *service {
	return &service{
		orderRepository:  orderRepository,
		outboxRepository: outboxRepository,
		txManager:        txManager,
		inventoryClient:  inventoryClient,
		paymentClient:    paymentClient,
		orderPaidTopic:   orderPaidTopic,
	}
}
