package redis

import (
	"context"
	"errors"

	redispet "github.com/Krokozabra213/sso/pkg/db/redis-pet"
)

func ErrorFactory(err error) error {

	switch {
	case errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled):
		return ErrContext
	}

	// ошибки redispet
	var customErr *redispet.CustomError
	if errors.As(err, &customErr) {
		switch {
		case errors.Is(customErr, redispet.ErrAuth):
			return ErrRedisAuth
		case errors.Is(customErr, redispet.ErrContext):
			return ErrRedisCtx
		case errors.Is(customErr, redispet.ErrConnection):
			return ErrRedisConnection
		case errors.Is(customErr, redispet.ErrOOM):
			return ErrRedisOOM
		case errors.Is(customErr, redispet.ErrInternal):
			return ErrRedisInternal
		}
	}

	return ErrUnknown
}
