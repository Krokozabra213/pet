package sendtoclienthander

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrSendMessage = errors.New("error send message")
)

type MessageSender struct {
	log         *slog.Logger
	stream      chat.Chat_ChatStreamServer
	sendTimeout time.Duration
}

func NewMessageSender(log *slog.Logger, timeout time.Duration, stream chat.Chat_ChatStreamServer) *MessageSender {
	return &MessageSender{
		log:         log,
		sendTimeout: timeout,
		stream:      stream,
	}
}

func (ms *MessageSender) SendToClient(message *chat.ServerMessage) error {
	streamCtx := ms.stream.Context()
	if streamCtx.Err() != nil {
		ms.log.Debug(streamCtx.Err().Error())
		err := HandleStreamContextError(streamCtx)
		return status.Error(codes.Canceled, err.Error())
	}

	// Создаем контекст с таймаутом для отправки
	ctx, cancel := context.WithTimeout(streamCtx, ms.sendTimeout)
	defer cancel()

	// Используем канал для неблокирующей отправки
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)
		errCh <- ms.stream.Send(message)
	}()

	select {
	case err := <-errCh:
		if err != nil { // вернулась ошибка отправки сообщения
			ms.log.Error("stream send failed", "error", err)
			return status.Error(codes.Internal, ErrSendMessage.Error())
		}
		return nil // отправка сообщения прошла успешна, метод Send вернул nil
	case <-ctx.Done():
		ms.log.Debug(ctx.Err().Error())
		err := HandleStreamContextError(ctx)
		return status.Error(codes.Canceled, err.Error())
	}
}
