package metrics

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/platform/pkg/logger"
)

type Metrics struct {
	OrdersTotal   metric.Int64Counter
	OrdersRevenue metric.Float64Counter
}

func New(ctx context.Context, m metric.Meter) (*Metrics, error) {
	ordersTotal, err := m.Int64Counter("orders_total")
	if err != nil {
		logger.Error(ctx, "failed to create metric orders_total", zap.Error(err))
		return nil, err
	}

	ordersRevenue, err := m.Float64Counter("orders_revenue_total")
	if err != nil {
		logger.Error(ctx, "failed to create metric orders_revenue_total", zap.Error(err))
		return nil, err
	}

	return &Metrics{
		OrdersTotal:   ordersTotal,
		OrdersRevenue: ordersRevenue,
	}, nil
}
