package postgresrepo

import (
	"errors"

	"github.com/Krokozabra213/sso/internal/chat/repository/storage"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
)

func ErrorFactory(err error) error {

	//ошибки postgres-pet
	var customErr *postgrespet.CustomError
	if errors.As(err, &customErr) {
		switch {
		case errors.Is(err, postgrespet.ErrTransaction):
			return storage.ErrTransaction
		case errors.Is(err, postgrespet.ErrCtxCancelled):
			return storage.ErrCtxCancelled
		case errors.Is(err, postgrespet.ErrCtxDeadline):
			return storage.ErrCtxDeadline
		case errors.Is(err, postgrespet.ErrValidation):
			return storage.ErrValidation
		case errors.Is(err, postgrespet.ErrDuplicateKey):
			return storage.ErrDuplicate
		case errors.Is(err, postgrespet.ErrNotFound):
			return storage.ErrNotFound
		case errors.Is(err, postgrespet.ErrInternal):
			return storage.ErrInternal
		}
	}

	return storage.ErrUnknown
}
