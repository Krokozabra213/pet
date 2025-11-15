package redis

import (
	"errors"

	redispet "github.com/Krokozabra213/sso/pkg/db/redis-pet"
)

func ErrorFactory(err *redispet.CustomError) error {

	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, redispet.ErrAuth):
		return ErrRedisAuth
	case errors.Is(err, redispet.ErrContext):
		return ErrRedisCtx
	case errors.Is(err, redispet.ErrConnection):
		return ErrRedisConnection
	case errors.Is(err, redispet.ErrOOM):
		return ErrRedisOOM
	case errors.Is(err, redispet.ErrInternal):
		return ErrRedisInternal
	default:
		return ErrUnknown
	}
}
