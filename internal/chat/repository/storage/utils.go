package storage

import (
	"context"
	"errors"
)

var (
	ErrCtxCancelled = errors.New("context cancelled error")
	ErrCtxDeadline  = errors.New("context deadline error")
	ErrUnknown      = errors.New("unknown error")
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
