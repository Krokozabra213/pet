package storage

import (
	"context"
	"errors"
)

func CtxError(err error) error {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return ErrCtxDeadline
	case errors.Is(err, context.Canceled):
		return ErrCtxCancelled
	default:
		return ErrUnknown
	}
}
