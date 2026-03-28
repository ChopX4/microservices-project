package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/ChopX4/raketka/order/internal/repository/converter"
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

	tag, err := r.db.Exec(
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
		return err
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
