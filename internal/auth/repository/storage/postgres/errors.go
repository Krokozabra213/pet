package postgres

import "errors"

var (
	ErrPGValidation = errors.New("validation error")
	ErrPGDuplicate  = errors.New("duplicate key error")
	ErrPGNotFound   = errors.New("not found error")
	ErrPGInternal   = errors.New("internal error")
	ErrContext      = errors.New("context error")
	ErrUnknown      = errors.New("unknown error")
)
