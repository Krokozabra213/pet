package postgres

import (
	"errors"

	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
)

func ErrorFactory(err *postgrespet.CustomError) error {

	if err == nil {
		return nil
	}

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
	return ErrUnknown
}
