package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/assembly/internal/config"
	kafkaConverter "github.com/ChopX4/raketka/assembly/internal/converter/kafka"
	decoder "github.com/ChopX4/raketka/assembly/internal/converter/kafka/decoder"
	"github.com/ChopX4/raketka/assembly/internal/service"
	orderconsumer "github.com/ChopX4/raketka/assembly/internal/service/consumer/order_consumer"
	orderconsumermetrics "github.com/ChopX4/raketka/assembly/internal/service/consumer/order_consumer/metrics"
	orderproducer "github.com/ChopX4/raketka/assembly/internal/service/producer/order_producer"
	"github.com/ChopX4/raketka/platform/pkg/closer"
	kafkaConsumer "github.com/ChopX4/raketka/platform/pkg/kafka/consumer"
	kafkaProducer "github.com/ChopX4/raketka/platform/pkg/kafka/producer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	platformMetrics "github.com/ChopX4/raketka/platform/pkg/metrics"
)

type diContainer struct {
	orderConsumer   service.OrderConsumer
	orderProducer   service.ShipProducer
	orderDecoder    kafkaConverter.OrderPaidDecoder
	consumerMetrics *orderconsumermetrics.Metrics

	syncProducer  sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InitMetrics(ctx context.Context) error {
	if err := platformMetrics.InitProvider(
		ctx,
		config.AppConfig().Metrics.CollectorEndpoint(),
		config.AppConfig().Metrics.CollectorInterval(),
		"assembly",
	); err != nil {
		return err
	}

	closer.AddNamed("otel meter provider", platformMetrics.Shutdown)

	return nil
}

func (d *diContainer) OrderConsumer(ctx context.Context) service.OrderConsumer {
	if d.orderConsumer == nil {
		d.orderConsumer = orderconsumer.NewOrderConsumer(
			kafkaConsumer.NewConsumer(
				d.ConsumerGroup(ctx),
				[]string{config.AppConfig().OrderPaidConsumer.Topic()},
				logger.Logger(),
			),
			d.OrderDecoder(),
			d.OrderProducer(ctx),
			d.OrderConsumerMetrics(ctx),
		)
	}

	return d.orderConsumer
}

func (d *diContainer) OrderConsumerMetrics(ctx context.Context) *orderconsumermetrics.Metrics {
	if d.consumerMetrics == nil {
		m, err := orderconsumermetrics.New(ctx, otel.Meter("assembly.service"))
		if err != nil {
			logger.Error(ctx, "failed to create assembly metrics", zap.Error(err))
			panic(fmt.Sprintf("failed to create assembly metrics: %v", err))
		}

		d.consumerMetrics = m
	}

	return d.consumerMetrics
}

func (d *diContainer) OrderProducer(ctx context.Context) service.ShipProducer {
	if d.orderProducer == nil {
		d.orderProducer = orderproducer.NewShipProducer(
			kafkaProducer.NewProducer(
				d.SyncProducer(ctx),
				config.AppConfig().OrderAssembledProducer.Topic(),
				logger.Logger(),
			),
		)
	}

	return d.orderProducer
}

func (d *diContainer) OrderDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderDecoder == nil {
		d.orderDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderDecoder
}

func (d *diContainer) SyncProducer(ctx context.Context) sarama.SyncProducer {
	if d.syncProducer == nil {
		cfg := sarama.NewConfig()
		cfg.Version = sarama.V3_6_0_0
		cfg.Producer.Return.Successes = true

		producer, err := sarama.NewSyncProducer(config.AppConfig().Kafka.Brokers(), cfg)
		if err != nil {
			logger.Error(ctx, "failed to create kafka sync producer", zap.Error(err))
			panic(fmt.Sprintf("failed to create kafka sync producer: %v", err))
		}

		closer.AddNamed("kafka sync producer", func(context.Context) error {
			return producer.Close()
		})

		d.syncProducer = producer
	}

	return d.syncProducer
}

func (d *diContainer) ConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		cfg := sarama.NewConfig()
		cfg.Version = sarama.V3_6_0_0
		cfg.Consumer.Return.Errors = true

		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			cfg,
		)
		if err != nil {
			logger.Error(ctx, "failed to create kafka consumer group", zap.Error(err))
			panic(fmt.Sprintf("failed to create kafka consumer group: %v", err))
		}

		closer.AddNamed("kafka consumer group", func(context.Context) error {
			return group.Close()
		})

		d.consumerGroup = group
	}

	return d.consumerGroup
}
