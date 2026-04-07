package service

import (
	"context"

	"github.com/ChopX4/raketka/notification/internal/model"
)

type TelegramService interface {
	SendShipNotification(ctx context.Context, event model.ShipAssembled) error
	SendOrderNotification(ctx context.Context, event model.OrderPaid) error
}

type OrderConsumer interface {
	RunOrderConsumer(ctx context.Context) error
}

type AssembledConsumer interface {
	RunAssembledConsumer(ctx context.Context) error
}
