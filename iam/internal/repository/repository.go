package repository

import (
	"context"
	"time"

	"github.com/ChopX4/raketka/iam/internal/model"
)

type SessionRepository interface {
	Create(ctx context.Context, session model.Session, ttl time.Duration) error
	GetUserUUID(ctx context.Context, sessionUUID string) (string, error)
}

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	GetByLogin(ctx context.Context, login string) (model.User, error)
	Get(ctx context.Context, userUUID string) (model.User, error)
}
