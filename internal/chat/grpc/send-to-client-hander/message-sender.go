package sendtoclienthander

import (
	"context"
	"errors"
	"time"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	// chatgrpc "github.com/Krokozabra213/sso/internal/chat/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrSendMessage = errors.New("error send message")
)

type MessageSender struct {
	stream      chat.Chat_ChatStreamServer
	sendTimeout time.Duration
}

func NewMessageSender(timeout time.Duration, stream chat.Chat_ChatStreamServer) *MessageSender {
	return &MessageSender{
		sendTimeout: timeout,
		stream:      stream,
	}
}

func (ms *MessageSender) SendToClient(message *chat.ServerMessage) error {
	// Создаем контекст с таймаутом для отправки
	ctx, cancel := context.WithTimeout(ms.stream.Context(), ms.sendTimeout)
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
			return status.Error(codes.Internal, ErrSendMessage.Error())
		}
		return nil // отправка сообщения прошла успешна, метод Send вернул nil
	case <-ctx.Done():
		err := HandleStreamContextError(ctx)
		return status.Error(codes.Canceled, err.Error())
	}
}
