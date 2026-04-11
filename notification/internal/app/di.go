package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"
	"go.uber.org/zap"

	notificationClient "github.com/ChopX4/raketka/notification/internal/client"
	telegramclient "github.com/ChopX4/raketka/notification/internal/client/telegram"
	"github.com/ChopX4/raketka/notification/internal/config"
	kafkaConverter "github.com/ChopX4/raketka/notification/internal/converter/kafka"
	decoder "github.com/ChopX4/raketka/notification/internal/converter/kafka/decoder"
	"github.com/ChopX4/raketka/notification/internal/service"
	assembledconsumer "github.com/ChopX4/raketka/notification/internal/service/assembled_consumer"
	orderconsumer "github.com/ChopX4/raketka/notification/internal/service/order_consumer"
	telegramservice "github.com/ChopX4/raketka/notification/internal/service/telegram"
	"github.com/ChopX4/raketka/platform/pkg/closer"
	kafkaConsumer "github.com/ChopX4/raketka/platform/pkg/kafka/consumer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

type diContainer struct {
	telegramService   service.TelegramService
	orderConsumer     service.OrderConsumer
	assembledConsumer service.AssembledConsumer

	telegramClient notificationClient.TelegramClient
	telegramBot    *bot.Bot

	orderDecoder     kafkaConverter.OrderDecoder
	assembledDecoder kafkaConverter.AssembledDecoder

	orderConsumerGroup     sarama.ConsumerGroup
	assembledConsumerGroup sarama.ConsumerGroup
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramservice.NewTelegramService(d.TelegramClient(ctx))
	}

	return d.telegramService
}

func (d *diContainer) OrderConsumer(ctx context.Context) service.OrderConsumer {
	if d.orderConsumer == nil {
		d.orderConsumer = orderconsumer.NewOrderConsumer(
			kafkaConsumer.NewConsumer(
				d.OrderConsumerGroup(ctx),
				[]string{config.AppConfig().OrderConsumer.Topic()},
				logger.Logger(),
			),
			d.OrderDecoder(),
			d.TelegramService(ctx),
		)
	}

	return d.orderConsumer
}

func (d *diContainer) AssembledConsumer(ctx context.Context) service.AssembledConsumer {
	if d.assembledConsumer == nil {
		d.assembledConsumer = assembledconsumer.NewAssembledConsumer(
			kafkaConsumer.NewConsumer(
				d.AssembledConsumerGroup(ctx),
				[]string{config.AppConfig().AssembledConsumer.Topic()},
				logger.Logger(),
			),
			d.AssembledDecoder(),
			d.TelegramService(ctx),
		)
	}

	return d.assembledConsumer
}

func (d *diContainer) TelegramClient(ctx context.Context) notificationClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramclient.NewTelegramClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		telegramBot, err := bot.New(
			config.AppConfig().TelegramConfig.Token(),
			bot.WithCheckInitTimeout(10*time.Second),
			bot.WithHTTPClient(
				10*time.Second,
				&http.Client{Timeout: 10 * time.Second},
			),
		)
		if err != nil {
			logger.Error(ctx, "failed to create telegram bot", zap.Error(err))
			panic(fmt.Sprintf("failed to create telegram bot: %v", err))
		}

		d.telegramBot = telegramBot
	}

	return d.telegramBot
}

func (d *diContainer) OrderDecoder() kafkaConverter.OrderDecoder {
	if d.orderDecoder == nil {
		d.orderDecoder = decoder.NewOrderDecoder()
	}

	return d.orderDecoder
}

func (d *diContainer) AssembledDecoder() kafkaConverter.AssembledDecoder {
	if d.assembledDecoder == nil {
		d.assembledDecoder = decoder.NewAssembledDecoder()
	}

	return d.assembledDecoder
}

func (d *diContainer) OrderConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if d.orderConsumerGroup == nil {
		cfg := sarama.NewConfig()
		cfg.Version = sarama.V3_6_0_0
		cfg.Consumer.Return.Errors = true

		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderConsumer.GroupID(),
			cfg,
		)
		if err != nil {
			logger.Error(ctx, "failed to create order consumer group", zap.Error(err))
			panic(fmt.Sprintf("failed to create order consumer group: %v", err))
		}

		closer.AddNamed("order consumer group", func(context.Context) error {
			return group.Close()
		})

		d.orderConsumerGroup = group
	}

	return d.orderConsumerGroup
}

func (d *diContainer) AssembledConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if d.assembledConsumerGroup == nil {
		cfg := sarama.NewConfig()
		cfg.Version = sarama.V3_6_0_0
		cfg.Consumer.Return.Errors = true

		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().AssembledConsumer.GroupID(),
			cfg,
		)
		if err != nil {
			logger.Error(ctx, "failed to create assembled consumer group", zap.Error(err))
			panic(fmt.Sprintf("failed to create assembled consumer group: %v", err))
		}

		closer.AddNamed("assembled consumer group", func(context.Context) error {
			return group.Close()
		})

		d.assembledConsumerGroup = group
	}

	return d.assembledConsumerGroup
}
