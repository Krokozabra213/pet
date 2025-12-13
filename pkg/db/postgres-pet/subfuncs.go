package postgrespet

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func isDuplicateKey(err *pgconn.PgError) bool {
	return err.Code == CodeUniqueViolation
}

func isValidation(err *pgconn.PgError) bool {
	validationErrorCodes := map[string]struct{}{
		CodeNotNullViolation:          {},
		CodeForeignKeyViolation:       {},
		CodeCheckViolation:            {},
		CodeStringTooLong:             {},
		CodeNumericOutOfRange:         {},
		CodeInvalidDatetimeFormat:     {},
		CodeDatetimeOverflow:          {},
		CodeDivisionByZero:            {},
		CodeInvalidTextRepresentation: {},
	}

	_, exists := validationErrorCodes[err.Code]
	return exists
}

func isNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func isRetryableTx(err *pgconn.PgError) bool {
	switch err.Code {
	case CodeSerializationFailure, CodeDeadlockDetected, CodeLockNotAvailable:
		return true
	default:
		return false
	}
}
