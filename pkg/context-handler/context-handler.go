package contexthandler

import (
	"context"
	"time"
)

func EnsureCtxTimeout(ctx context.Context, defaultTimeout time.Duration) (context.Context, context.CancelFunc) {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		return context.WithTimeout(ctx, defaultTimeout)
	}
	return ctx, func() {} // пустая функция, если deadline уже есть
}
