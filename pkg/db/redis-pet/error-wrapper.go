package redispet

import (
	"context"
	"errors"
	"strings"

	"github.com/redis/go-redis/v9"
)

var (
	ErrAuth         = errors.New("REDIS_AUTH_ERROR")
	ErrOOM          = errors.New("REDIS_OUT_OF_MEMORY")
	ErrConnection   = errors.New("REDIS_CONNECT_ERROR")
	ErrCtxCancelled = errors.New("CTX_CANCELLED_ERROR")
	ErrCtxDeadline  = errors.New("CTX_DEADLINE_ERROR")
	ErrInternal     = errors.New("INTERNAL_ERROR")
)

func ErrorWrapper(err error) *CustomError {
	if err == nil {
		return nil
	}

	var redisErr *redis.Error
	if errors.As(err, redisErr) {

		switch {
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
		case isAuth(err):
			return &CustomError{
				Message: ErrAuth.Error(),
				Err:     err,
			}
		case isOOM(err):
			return &CustomError{
				Message: ErrOOM.Error(),
				Err:     err,
			}
		case isConnect(err):
			return &CustomError{
				Message: ErrConnection.Error(),
				Err:     err,
			}
		}
	}

	return &CustomError{
		Message: ErrInternal.Error(),
		Err:     err,
	}
}

// Вспомогательные функции
func isAuth(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToUpper(err.Error())

	// Различные варианты ошибок аутентификации в Redis
	authPatterns := []string{
		"AUTH",
		"AUTHENTICATION",
		"WRONGPASS",
		"NOAUTH",
		"NO AUTH",
		"INVALID PASSWORD",
		"AUTH FAILED",
		"AUTHENTICATION FAILED",
		"AUTH REQUIRE",
		"REQUIRES AUTH",
		"NOAUTH AUTH",
	}

	for _, pattern := range authPatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}

func isOOM(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToUpper(err.Error())

	// Различные варианты ошибок аутентификации в Redis
	oomPatterns := []string{
		"OOM",
		"OUT OF MEMORY",
		"MAXMEMORY",
		"MEMORY LIMIT",
		"MEMORY EXCEEDED",
		"MEMORY USAGE",
	}

	for _, pattern := range oomPatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}

func isConnect(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToUpper(err.Error())

	// Различные варианты ошибок аутентификации в Redis
	connectPatterns := []string{
		"CONNECT",
		"CONNECTION",
		"NETWORK",
		"TIMEOUT",
		"REFUSED",
		"DIAL",
		"SOCKET",
		"HOST",
		"PORT",
		"UNREACHABLE",
		"NO ROUTE TO HOST",
		"RESET BY PEER",
		"BROKEN PIPE",
		"CLOSED",
		"EOF",
		"IO TIMEOUT",
		"CONTEXT DEADLINE EXCEEDED",
		"CONNECTION POOL",
		"POOL TIMEOUT",
		"UNABLE TO CONNECT",
		"FAILED TO CONNECT",
		"CANNOT CONNECT",
		"LOST CONNECTION",
	}

	for _, pattern := range connectPatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}

func ctxCancelled(err error) bool {
	return errors.Is(err, context.Canceled)
}

func ctxDeadline(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
}
