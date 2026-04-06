package service

import (
	"context"

	"github.com/ChopX4/raketka/assembly/internal/model"
)

type OrderConsumer interface {
	RunOrderConsumer(ctx context.Context) error
}

type ShipProducer interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error
}
