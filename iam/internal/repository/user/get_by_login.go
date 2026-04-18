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

func (r *repository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	sqlQuery := `
		SELECT uuid, login, email, password_hash, notification_methods FROM users
		WHERE login = $1
	`
	var user repoModel.User
	var notificationMethods []byte

	if err := pgxtx.GetQueryEngine(ctx, r.db).QueryRow(
		ctx,
		sqlQuery,
		login,
	).Scan(
		&user.Uuid,
		&user.Login,
		&user.Email,
		&user.HashPassword,
		&notificationMethods,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		logger.Error(ctx, "failed to get user by login from postgres", zap.String("user_login", login), zap.Error(err))
		return model.User{}, err
	}

	if len(notificationMethods) > 0 {
		if err := json.Unmarshal(notificationMethods, &user.NotificationMethods); err != nil {
			logger.Error(ctx, "failed unmarshal notification methods", zap.Error(err))
			return model.User{}, err
		}
	}

	return converter.UserToModel(user), nil
}
