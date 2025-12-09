package storage

import "errors"

var (
	ErrCtxCancelled = errors.New("context cancelled error")
	ErrCtxDeadline  = errors.New("context deadline error")
	ErrUnknown      = errors.New("unknown error")

	// postgres errors
	ErrValidation = errors.New("validation error")
	ErrDuplicate  = errors.New("duplicate key error")
	ErrNotFound   = errors.New("not found error")
	ErrInternal   = errors.New("internal error")
)
