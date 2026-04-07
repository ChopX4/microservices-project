package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"go.uber.org/zap"

	"github.com/ChopX4/raketka/notification/internal/client"
	"github.com/ChopX4/raketka/notification/internal/model"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

const chatId = 943519931

//go:embed templates/*.tmpl
var templatesFS embed.FS

type shipAssembledTemplate struct {
	EventUuid    string
	OrderUuid    string
	UserUuid     string
	BuildTimeSec int64
}

type orderPaidTemplate struct {
	EventUuid       string
	OrderUuid       string
	UserUuid        string
	TransactionUuid string
}

var (
	shipTemplate  = template.Must(template.ParseFS(templatesFS, "templates/ship_assembled.tmpl"))
	orderTemplate = template.Must(template.ParseFS(templatesFS, "templates/order_paid.tmpl"))
)

type service struct {
	telegramClient client.TelegramClient
}

func NewTelegramService(telegramClient client.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) SendShipNotification(ctx context.Context, event model.ShipAssembled) error {
	message, err := s.buildShipMessage(event)
	if err != nil {
		return err
	}

	if err := s.telegramClient.SendMessage(ctx, chatId, message); err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatId), zap.String("message", message))
	return nil
}

func (s *service) SendOrderNotification(ctx context.Context, event model.OrderPaid) error {
	message, err := s.buildOrderMessage(event)
	if err != nil {
		return err
	}

	if err := s.telegramClient.SendMessage(ctx, chatId, message); err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatId), zap.String("message", message))
	return nil
}

func (s *service) buildShipMessage(event model.ShipAssembled) (string, error) {
	data := shipAssembledTemplate{
		EventUuid:    event.EventUuid,
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: event.BuildTimeSec,
	}

	var buf bytes.Buffer
	if err := shipTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) buildOrderMessage(event model.OrderPaid) (string, error) {
	data := orderPaidTemplate{
		EventUuid:       event.EventUuid,
		OrderUuid:       event.OrderUuid,
		UserUuid:        event.UserUuid,
		TransactionUuid: event.TransactionUuid,
	}

	var buf bytes.Buffer
	if err := orderTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
