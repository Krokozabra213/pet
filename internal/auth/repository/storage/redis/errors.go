package redis

import "errors"

var (
	ErrRedisAuth       = errors.New("auth error")
	ErrRedisCtx        = errors.New("context expired error")
	ErrRedisOOM        = errors.New("out of memory error")
	ErrRedisConnection = errors.New("connect error")
	ErrRedisInternal   = errors.New("internal error")
	ErrTokenExpired    = errors.New("token already expired")
	ErrContext         = errors.New("context error")
	ErrUnknown         = errors.New("unknown error")
)
