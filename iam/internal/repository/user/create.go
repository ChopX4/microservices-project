package user

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/iam/internal/model"
	"github.com/ChopX4/raketka/iam/internal/repository/converter"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	"github.com/ChopX4/raketka/platform/pkg/pgxtx"
)

func (r *repository) Create(ctx context.Context, user model.User) error {
	repoUser := converter.UserToRepo(user)

	sqlQuery := `
		INSERT INTO users (
			uuid,
			login,
			email,
			password_hash,
			notification_methods
		) VALUES ($1, $2, $3, $4, $5)
	`
	notificationMethods, err := json.Marshal(repoUser.NotificationMethods)
	if err != nil {
		logger.Error(ctx, "failed to marshal notification methods", zap.Error(err))
		return err
	}

	_, err = pgxtx.GetQueryEngine(ctx, r.db).Exec(
		ctx,
		sqlQuery,
		repoUser.Uuid,
		repoUser.Login,
		repoUser.Email,
		repoUser.HashPassword,
		notificationMethods,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.ErrUserAlreadyExists
		}

		logger.Error(ctx, "failed to create user in postgres", zap.String("user_uuid", user.Uuid), zap.Error(err))
		return err
	}

	return nil
}
