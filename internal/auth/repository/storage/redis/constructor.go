package redis

import (
	"time"

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
