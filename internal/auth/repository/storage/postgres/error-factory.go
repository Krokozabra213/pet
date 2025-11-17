package postgres

import (
	"context"
	"errors"

	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
)

func ErrorFactory(err error) error {

	switch {
	case errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled):
		return ErrContext
	}

	//ошибки postgres-pet
	var customErr *postgrespet.CustomError
	if errors.As(err, &customErr) {
		switch {
		case errors.Is(err, postgrespet.ErrValidation):
			return ErrPGValidation
		case errors.Is(err, postgrespet.ErrDuplicateKey):
			return ErrPGDuplicate
		case errors.Is(err, postgrespet.ErrNotFound):
			return ErrPGNotFound
		case errors.Is(err, postgrespet.ErrInternal):
			return ErrPGInternal
		}
	}

	return ErrUnknown
}
