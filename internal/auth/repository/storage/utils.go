package storage

import (
	"context"
	"errors"
	"time"
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

func EnsureCtxTimeout(ctx context.Context, defaultTimeout time.Duration) (context.Context, context.CancelFunc) {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		return context.WithTimeout(ctx, defaultTimeout)
	}
	return ctx, func() {} // пустая функция, если deadline уже есть
}
