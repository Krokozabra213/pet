package redis

import (
	"context"
	"log/slog"
	"time"

	redispet "github.com/Krokozabra213/sso/pkg/db/redis-pet"
)

const (
	ctxTimeout = 5 * time.Second
)

type Redis struct {
	RDB *redispet.RDB
	log *slog.Logger
}

func New(RDB *redispet.RDB, log *slog.Logger) *Redis {
	return &Redis{
		RDB: RDB,
		log: log,
	}
}

func (r *Redis) SaveToken(ctx context.Context, token string, expiresAt time.Time) error {

	const op = "redis.SaveToken"
	log := r.log.With(
		slog.String("op", op),
	)

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ctxTimeout)
		defer cancel()
	}

	if ctx.Err() != nil {
		log.Error("context error", "err", ctx.Err())
		return ErrContext
	}

	expiration := time.Until(expiresAt)
	if expiration <= 0 {
		return ErrTokenExpired
	}
	log.Error("repository error", "err", ErrTokenExpired)

	start := time.Now()
	err := r.RDB.Client.SetEx(ctx, token, "", expiration).Err()
	duration := time.Since(start)

	if duration > 100*time.Millisecond {
		log.Warn("slow redis operation", "duration", duration)
	}

	customErr := redispet.ErrorWrapper(err)
	if customErr != nil {
		log.Error("redis error", "err", customErr.Error())
		err = ErrorFactory(customErr)
		return err
	}

	return nil
}

func (r *Redis) CheckToken(ctx context.Context, token string) (bool, error) {

	const op = "redis.CheckToken"
	log := r.log.With(
		slog.String("op", op),
	)

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ctxTimeout)
		defer cancel()
	}

	if ctx.Err() != nil {
		log.Error("context error", "err", ctx.Err())
		return false, ErrContext
	}

	start := time.Now()
	exists, err := r.RDB.Client.Exists(ctx, token).Result()
	duration := time.Since(start)

	if duration > 100*time.Millisecond {
		log.Warn("slow redis operation", "duration", duration)
	}

	customErr := redispet.ErrorWrapper(err)
	if customErr != nil {
		log.Error("redis error", "err", customErr.Error())
		err = ErrorFactory(customErr)
		return false, err
	}

	return exists == 1, nil
}
