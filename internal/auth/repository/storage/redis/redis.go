package redis

import (
	"context"
	"time"

	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
	"github.com/Krokozabra213/sso/pkg/db"
)

type Redis struct {
	RDB *db.RDB
}

func New(RDB *db.RDB) *Redis {
	return &Redis{RDB: RDB}
}

func (r *Redis) SaveToken(ctx context.Context, token string, expiresAt time.Time) error {
	expiration := time.Until(expiresAt)
	if expiration <= 0 {
		return storage.ErrTokenExpired
	}

	err := r.RDB.DB.SetEx(ctx, token, "", expiration).Err()

	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) CheckToken(ctx context.Context, token string) (bool, error) {
	exists, err := r.RDB.DB.Exists(ctx, token).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
