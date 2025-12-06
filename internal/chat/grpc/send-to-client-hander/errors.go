package sendtoclienthander

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnknownMessageType = errors.New("unknown message type")
	ErrGracefulShutdown   = errors.New("graceful shutdown")

	ErrContextCancelled = errors.New("stream cancelled")
	ErrContextDeadline  = errors.New("stream deadline exceeded")
	ErrContextUnknown   = errors.New("stream context unknown error")
)

func HandleStreamContextError(ctx context.Context) error {

	// Определяем тип ошибки контекста
	switch {
	case errors.Is(ctx.Err(), context.Canceled):
		// Клиент отменил запрос или соединение закрыто
		return status.Error(codes.Canceled, ErrContextCancelled.Error())

	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		// Превышен таймаут
		return status.Error(codes.DeadlineExceeded, ErrContextDeadline.Error())

	default:
		// Любая другая ошибка контекста
		return status.Error(codes.Aborted, ErrContextUnknown.Error())
	}
}
