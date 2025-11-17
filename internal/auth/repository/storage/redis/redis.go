package redis

import (
	"context"
	"log/slog"
	"time"

	redispet "github.com/Krokozabra213/sso/pkg/db/redis-pet"
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

	expiration := time.Until(expiresAt)
	if expiration <= 0 {
		return ErrTokenExpired
	}
	log.Error("repository error", "err", ErrTokenExpired)

	err := r.RDB.Client.SetEx(ctx, token, "", expiration).Err()

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

	exists, err := r.RDB.Client.Exists(ctx, token).Result()
	customErr := redispet.ErrorWrapper(err)
	if customErr != nil {
		log.Error("redis error", "err", customErr.Error())
		err = ErrorFactory(customErr)
		return false, err
	}

	return exists == 1, nil
}
