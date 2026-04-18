package session

import (
	"context"
	"errors"

	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/iam/internal/model"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

func (r *repository) GetUserUUID(ctx context.Context, sessionUUID string) (string, error) {
	cacheKey := r.getCacheKey(sessionUUID)

	user, err := r.cache.Get(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return "", model.ErrSessionNotFound
		}

		logger.Error(
			ctx,
			"failed to get session from redis",
			zap.String("session_uuid", sessionUUID),
			zap.Error(err),
		)
		return "", err
	}

	return string(user), nil
}
