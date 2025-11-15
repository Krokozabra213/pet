package redis

import "errors"

var (
	ErrRedisAuth       = errors.New("auth error")
	ErrRedisCtx        = errors.New("context expired error")
	ErrRedisOOM        = errors.New("out of memory error")
	ErrRedisConnection = errors.New("connect error")
	ErrRedisInternal   = errors.New("internal error")
	ErrUnknown         = errors.New("unknown error")
)
