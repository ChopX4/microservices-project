package metrics

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/platform/pkg/logger"
)

type Metrics struct {
	AssemblyDuration metric.Float64Histogram
}

func New(ctx context.Context, m metric.Meter) (*Metrics, error) {
	assemblyDuration, err := m.Float64Histogram("assembly_duration_seconds")
	if err != nil {
		logger.Error(ctx, "failed to create metric assembly_duration_seconds", zap.Error(err))
		return nil, err
	}

	return &Metrics{
		AssemblyDuration: assemblyDuration,
	}, nil
}
