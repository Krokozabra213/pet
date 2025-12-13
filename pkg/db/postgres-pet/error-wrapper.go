package postgrespet

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func ErrorWrapper(err error) error {
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
		case isRetryableTx(pgError):
			return &CustomError{
				Message: ErrTransaction.Error(),
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
