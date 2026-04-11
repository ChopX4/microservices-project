package repository

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
	outboxModel "github.com/ChopX4/raketka/order/internal/repository/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.OrderByUUID) error
	Get(ctx context.Context, orderUUID string) (model.OrderByUUID, error)
	Update(ctx context.Context, order model.OrderByUUID) error
}

type OutboxRepository interface {
	Create(ctx context.Context, msg outboxModel.OutboxMessage) error
	ListPending(ctx context.Context, limit int) ([]outboxModel.OutboxMessage, error)
	MarkPublished(ctx context.Context, eventUUID string) error
}
