package user

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/iam/internal/model"
	"github.com/ChopX4/raketka/iam/internal/repository/converter"
	repoModel "github.com/ChopX4/raketka/iam/internal/repository/model"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	"github.com/ChopX4/raketka/platform/pkg/pgxtx"
)

func (r *repository) Get(ctx context.Context, userUUID string) (model.User, error) {
	var user repoModel.User
	var rawNotificationMethods []byte

	sqlQuery := `
		SELECT uuid, login, email, password_hash, notification_methods FROM users
		WHERE uuid = $1
	`

	if err := pgxtx.GetQueryEngine(ctx, r.db).QueryRow(
		ctx,
		sqlQuery,
		userUUID,
	).Scan(
		&user.Uuid,
		&user.Login,
		&user.Email,
		&user.HashPassword,
		&rawNotificationMethods,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		logger.Error(ctx, "failed to get user from postgres", zap.String("user_uuid", userUUID), zap.Error(err))
		return model.User{}, err
	}

	if len(rawNotificationMethods) > 0 {
		if err := json.Unmarshal(rawNotificationMethods, &user.NotificationMethods); err != nil {
			logger.Error(ctx, "failed unmarshal notification methods", zap.Error(err))
			return model.User{}, err
		}
	}

	return converter.UserToModel(user), nil
}
