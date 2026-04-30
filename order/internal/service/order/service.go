package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/clients/grpc"
	"github.com/ChopX4/raketka/order/internal/repository"
	orderMetrics "github.com/ChopX4/raketka/order/internal/service/order/metrics"
	"github.com/ChopX4/raketka/platform/pkg/pgxtx"
)

type service struct {
	orderRepository  repository.OrderRepository
	outboxRepository repository.OutboxRepository
	txManager        pgxtx.TxManager
	inventoryClient  grpc.InventoryClient
	paymentClient    grpc.PaymentClient
	orderPaidTopic   string
	metrics          *orderMetrics.Metrics
}

func NewService(orderRepository repository.OrderRepository, outboxRepository repository.OutboxRepository, txManager pgxtx.TxManager, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient, orderPaidTopic string, m *orderMetrics.Metrics) *service {
	return &service{
		orderRepository:  orderRepository,
		outboxRepository: outboxRepository,
		txManager:        txManager,
		inventoryClient:  inventoryClient,
		paymentClient:    paymentClient,
		orderPaidTopic:   orderPaidTopic,
		metrics:          m,
	}
}

func (s *service) addOrdersTotal(ctx context.Context, value int64) {
	if s.metrics == nil {
		return
	}
	s.metrics.OrdersTotal.Add(ctx, value)
}

func (s *service) addOrdersRevenue(ctx context.Context, value float64) {
	if s.metrics == nil {
		return
	}
	s.metrics.OrdersRevenue.Add(ctx, value)
}
