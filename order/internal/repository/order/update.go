package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/ChopX4/raketka/order/internal/model"
	repo "github.com/ChopX4/raketka/order/internal/repository"
	"github.com/ChopX4/raketka/order/internal/repository/converter"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

func (r *repository) Update(ctx context.Context, order model.OrderByUUID) error {
	repoOrder := converter.OrderByUUIDToRepo(order)

	sqlQuery := `
		UPDATE orders
		SET user_uuid = $2,
			part_uuids = $3,
			total_price = $4,
			transaction_uuid = $5,
			payment_method = $6,
			status = $7
		WHERE order_uuid = $1
	`

	tag, err := repo.GetQueryEngine(ctx, r.db).Exec(
		ctx,
		sqlQuery,
		repoOrder.OrderUUID,
		repoOrder.UserUUID,
		repoOrder.PartUuids,
		repoOrder.TotalPrice,
		repoOrder.TransactionUUID,
		repoOrder.PaymentMethod,
		repoOrder.Status,
	)
	if err != nil {
		logger.Error(ctx, "failed to update order in postgres", zap.String("order_uuid", order.OrderUUID.String()), zap.Error(err))
		return err
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
