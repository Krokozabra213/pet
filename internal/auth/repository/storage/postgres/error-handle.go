package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	CodeUniqueViolation = "23505"
)

func duplicateKey(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == CodeUniqueViolation {
			return true
		}
	}
	return false
}

func notFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
