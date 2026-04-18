package session

import (
	"fmt"

	"github.com/ChopX4/raketka/platform/pkg/cache"
)

const (
	cacheKeyPrefix = "session:"
)

type repository struct {
	cache cache.RedisClient
}

func NewSessionRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

func (r *repository) getCacheKey(sessionUUID string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, sessionUUID)
}
