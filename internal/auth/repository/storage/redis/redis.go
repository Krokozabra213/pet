package redis

import (
	"context"
	"time"

	"github.com/Krokozabra213/sso/internal/auth/domain"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
	redispet "github.com/Krokozabra213/sso/pkg/db/redis-pet"
)

const (
	ctxTimeout = 5 * time.Second
)

type Redis struct {
	RDB *redispet.RDB
}

func New(RDB *redispet.RDB) *Redis {
	return &Redis{
		RDB: RDB,
	}
}

func (r *Redis) SaveToken(parentCtx context.Context, token string, expiresAt time.Time) error {

	ctx, cancel := storage.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return storage.CtxError(ctx.Err())
	}

	expiration := time.Until(expiresAt)
	if expiration <= 0 {
		return storage.ErrTokenExpired
	}

	err := r.RDB.Client.SetEx(ctx, token, "", expiration).Err()

	customErr := redispet.ErrorWrapper(err)
	if customErr != nil {
		err = ErrorFactory(domain.TokenEntity, customErr)
		return err
	}

	return nil
}

func (r *Redis) CheckToken(parentCtx context.Context, token string) (bool, error) {

	ctx, cancel := storage.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return false, storage.CtxError(ctx.Err())
	}

	exists, err := r.RDB.Client.Exists(ctx, token).Result()

	customErr := redispet.ErrorWrapper(err)
	if customErr != nil {
		err = ErrorFactory(domain.TokenEntity, customErr)
		return false, err
	}

	return exists == 1, nil
}
