package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/ChopX4/raketka/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order model.OrderByUUID) error {
	repoOrder := converter.OrderByUUIDToRepo(order)

	sqlQuery := `
		INSERT INTO orders (
			order_uuid,
			user_uuid,
			part_uuids,
			total_price,
			transaction_uuid,
			payment_method,
			status
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.ErrAlreadyExists
		}

		return err
	}

	return nil
}
