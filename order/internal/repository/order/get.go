package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/ChopX4/raketka/order/internal/repository/converter"
	repoModel "github.com/ChopX4/raketka/order/internal/repository/model"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	"github.com/ChopX4/raketka/platform/pkg/pgxtx"
)

func (r *repository) Get(ctx context.Context, orderUUID string) (model.OrderByUUID, error) {
	if _, err := uuid.Parse(orderUUID); err != nil {
		return model.OrderByUUID{}, err
	}

	dbQuery := "SELECT order_uuid, user_uuid, part_uuids, total_price, transaction_uuid, payment_method, status FROM orders WHERE order_uuid = $1"
	var order repoModel.OrderByUUID

	if err := pgxtx.GetQueryEngine(ctx, r.db).QueryRow(ctx, dbQuery, orderUUID).Scan(
		&order.OrderUUID,
		&order.UserUUID,
		&order.PartUuids,
		&order.TotalPrice,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.Status,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.OrderByUUID{}, model.ErrNotFound
		}
		logger.Error(ctx, "failed to get order from postgres", zap.String("order_uuid", orderUUID), zap.Error(err))
		return model.OrderByUUID{}, err
	}

	return converter.OrderByUUIDToModel(order), nil
}
