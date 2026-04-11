package outbox

import (
	"context"

	repo "github.com/ChopX4/raketka/order/internal/repository"
	repoModel "github.com/ChopX4/raketka/order/internal/repository/model"
)

func (r *repository) ListPending(ctx context.Context, limit int) ([]repoModel.OutboxMessage, error) {
	if limit <= 0 {
		limit = 100
	}

	const sqlQuery = `
		SELECT event_uuid, topic, key, payload, status, created_at
		FROM events
		WHERE status = 'PENDING'
		ORDER BY created_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`
	rows, err := repo.GetQueryEngine(ctx, r.db).Query(ctx, sqlQuery, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]repoModel.OutboxMessage, 0)
	for rows.Next() {
		var msg repoModel.OutboxMessage
		if err := rows.Scan(
			&msg.EventUUID,
			&msg.Topic,
			&msg.Key,
			&msg.Payload,
			&msg.Status,
			&msg.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
