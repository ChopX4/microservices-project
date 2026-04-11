package outbox

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/order/internal/model"
	repo "github.com/ChopX4/raketka/order/internal/repository"
	repoModel "github.com/ChopX4/raketka/order/internal/repository/model"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

func (r *repository) Create(ctx context.Context, msg repoModel.OutboxMessage) error {
	sqlQuery := `INSERT INTO events (
    	event_uuid,
    	topic,
    	key,
    	payload
	) VALUES ($1, $2, $3, $4)
	`

	_, err := repo.GetQueryEngine(ctx, r.db).Exec(
		ctx,
		sqlQuery,
		msg.EventUUID,
		msg.Topic,
		msg.Key,
		msg.Payload,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.ErrAlreadyExists
		}

		logger.Error(ctx, "failed to create event in postgres", zap.String("event_uuid", msg.EventUUID), zap.Error(err))
		return err
	}

	return nil
}
