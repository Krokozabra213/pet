package redis

import (
	"errors"

	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
	redispet "github.com/Krokozabra213/sso/pkg/db/redis-pet"
)

func ErrorFactory(entity string, err error) error {

	// ошибки redis-pet
	var customErr *redispet.CustomError
	if errors.As(err, &customErr) {
		switch {
		case errors.Is(customErr, redispet.ErrAuth):
			return storage.ErrAuth
		case errors.Is(customErr, redispet.ErrCtxCancelled):
			return storage.ErrCtxCancelled
		case errors.Is(customErr, redispet.ErrCtxDeadline):
			return storage.ErrCtxDeadline
		case errors.Is(customErr, redispet.ErrConnection):
			return storage.ErrConnection
		case errors.Is(customErr, redispet.ErrOOM):
			return storage.ErrOOM
		case errors.Is(customErr, redispet.ErrInternal):
			return storage.ErrInternal
		}
	}

	return storage.ErrUnknown
}
