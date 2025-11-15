package redispet

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const ctxTimeout = 5 * time.Second

type RDB struct {
	Client *redis.Client
}

func NewRedisDB(addr, password string, db int) *RDB {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Redis connection failed: %v", err))
	}

	return &RDB{
		Client: client,
	}
}
