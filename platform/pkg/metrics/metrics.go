package metrics

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	defaultTimeout = 5 * time.Second
)

var (
	exporter      *otlpmetricgrpc.Exporter
	meterProvider *metric.MeterProvider
)

// InitProvider инициализирует глобальный провайдер метрик OpenTelemetry
func InitProvider(ctx context.Context, collectorEndpoint string, collectorInterval time.Duration, serviceName string) error {
	var err error

	if serviceName == "" {
		serviceName = "unknown_service"
	}

	// Создаем экспортер для отправки метрик в OTLP коллектор
	exporter, err = otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(collectorEndpoint),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithTimeout(defaultTimeout),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create metrics exporter")
	}

	// Создаем провайдер метрик
	meterProvider = metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(
				exporter,
				metric.WithInterval(collectorInterval),
			),
		),
		metric.WithResource(resource.NewWithAttributes(
			"",
			attribute.String("service.name", serviceName),
		)),
	)

	// Устанавливаем глобальный провайдер метрик
	otel.SetMeterProvider(meterProvider)

	return nil
}

// GetMeterProvider возвращает текущий провайдер метрик
func GetMeterProvider() *metric.MeterProvider {
	return meterProvider
}

// Shutdown закрывает провайдер метрик и экспортер
func Shutdown(ctx context.Context) error {
	if meterProvider == nil && exporter == nil {
		return nil
	}

	if meterProvider != nil {
		if err := meterProvider.Shutdown(ctx); err != nil {
			return errors.Wrap(err, "failed to shutdown meter provider")
		}
	}

	// meterProvider.Shutdown already shuts down its readers/exporters.
	// Calling exporter.Shutdown again can produce false-positive "already shutdown" errors.
	meterProvider = nil
	exporter = nil

	return nil
}
