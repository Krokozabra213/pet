package chatgrpc

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Krokozabra213/protos/gen/go/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateMessage(message *chat.ClientMessage) (interface{}, error) {
	switch msg := message.Type.(type) {
	// проверяем тип сообщения
	case *chat.ClientMessage_SendMessage:
		return msg.SendMessage, nil
	default:
		return nil, fmt.Errorf("%w: %T", ErrUnknownMessageType, message)
	}
}

func ValidateStreamRecvErrors(err error) error {
	if errors.Is(err, io.EOF) {
		return status.Error(codes.Canceled, ErrDisconect.Error())
	}
	return status.Error(codes.Internal, ErrStream.Error())
}

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
