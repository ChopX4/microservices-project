package outbox

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
	repo "github.com/ChopX4/raketka/order/internal/repository"
)

func (r *repository) MarkPublished(ctx context.Context, eventUUID string) error {
	const sqlQuery = `
		UPDATE events
		SET status = 'PUBLISHED'
		WHERE event_uuid = $1
	`

	tag, err := repo.GetQueryEngine(ctx, r.db).Exec(ctx, sqlQuery, eventUUID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
