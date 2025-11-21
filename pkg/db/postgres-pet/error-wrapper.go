package postgrespet

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	CodeUniqueViolation = "23505"

	// Коды ошибок валидации PostgreSQL
	CodeNotNullViolation          = "23502"
	CodeForeignKeyViolation       = "23503"
	CodeCheckViolation            = "23514"
	CodeStringTooLong             = "22001"
	CodeNumericOutOfRange         = "22003"
	CodeInvalidDatetimeFormat     = "22007"
	CodeDatetimeOverflow          = "22008"
	CodeDivisionByZero            = "22012"
	CodeInvalidTextRepresentation = "22P02"
)

var (
	ErrDuplicateKey = errors.New("DUPLICATE_KEY_ERROR")
	ErrNotFound     = errors.New("NOT_FOUND_ERROR")
	ErrValidation   = errors.New("VALIDATION_ERROR")
	ErrCtxCancelled = errors.New("CTX_CANCELLED_ERROR")
	ErrCtxDeadline  = errors.New("CTX_DEADLINE_ERROR")
	ErrInternal     = errors.New("INTERNAL_ERROR")
)

func ErrorWrapper(err error) *CustomError {
	if err == nil {
		return nil
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {

		switch {
		case isValidation(pgError):
			return &CustomError{
				Message: ErrValidation.Error(),
				Err:     err,
			}
		case isDuplicateKey(pgError):
			return &CustomError{
				Message: ErrDuplicateKey.Error(),
				Err:     err,
			}
		default:
			return &CustomError{
				Message: ErrInternal.Error(),
				Err:     err,
			}
		}
	}

	switch {
	case isNotFound(err):
		return &CustomError{
			Message: ErrNotFound.Error(),
			Err:     err,
		}
	case ctxCancelled(err):
		return &CustomError{
			Message: ErrCtxCancelled.Error(),
			Err:     err,
		}
	case ctxDeadline(err):
		return &CustomError{
			Message: ErrCtxDeadline.Error(),
			Err:     err,
		}
	default:
		return &CustomError{
			Message: ErrInternal.Error(),
			Err:     err,
		}
	}
}

func ctxCancelled(err error) bool {
	return errors.Is(err, context.Canceled)
}

func ctxDeadline(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
}

// Вспомогательные функции
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
