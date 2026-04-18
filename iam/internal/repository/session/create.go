package session

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/ChopX4/raketka/iam/internal/model"
	"github.com/ChopX4/raketka/iam/internal/repository/converter"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

func (r *repository) Create(ctx context.Context, session model.Session, ttl time.Duration) error {
	repoSession := converter.SessionToRepo(session)
	cacheKey := r.getCacheKey(repoSession.SessionUUID)

	if err := r.cache.SetWithTTL(ctx, cacheKey, repoSession.UserUUID, ttl); err != nil {
		logger.Error(
			ctx,
			"failed to create session in redis",
			zap.String("session_uuid", repoSession.SessionUUID),
			zap.String("user_uuid", repoSession.UserUUID),
			zap.Error(err),
		)
		return err
	}

	return nil
}
