package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ctxTimeout = 5 * time.Second
)

type IRedis interface {
	SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
}

type Redis struct {
	RDB IRedis
}

func New(RDB IRedis) *Redis {
	return &Redis{
		RDB: RDB,
	}
}
