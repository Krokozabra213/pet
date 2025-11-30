package chatgrpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
)

const (
	SendTimeout = 3 * time.Second
)

func SendFromSrvToClient(buf <-chan interface{}, done chan struct{}, userID int64, stream chat.Chat_ChatStreamServer) error {

	for {
		select {
		case <-done:
			// Канал закрыт - нормальное завершение
			// можно добавить ошибку, чат отключен
			return nil
		case message, ok := <-buf:
			if !ok {
				// Канал закрыт - нормальное завершение
				// можно добавить ошибку, чат отключен
				return nil
			}

			// валидируем сообщение
			serverMsg, err := convertToServerMessage(message)
			if err != nil {
				// пропускаем сообщение если не проходит валидацию
				log.Println(err) // добавить логирование ошибки через slog.logger handler слой (вынести функцию в функцию хендлера)
				continue
			}

			// отправляем сообщение клиенту
			if err := sendWithTimeout(stream, serverMsg, userID); err != nil {
				// логировать ошибку через slog.logger
				return fmt.Errorf("failed to send message to client %d: %w", userID, err) // вынести ошибку в errors.go
			}

		case <-stream.Context().Done():
			// Клиент отключился
			ctxErr := stream.Context().Err()
			log.Printf("Stream context done for user %d: %v", userID, ctxErr) // логировать ошибку через slog.logger
			return ctxErr
		}
	}
}

func convertToServerMessage(message interface{}) (*chat.ServerMessage, error) {
	switch msg := message.(type) {
	// какой-то пользователь подключился
	case *chat.UserJoined:

		return &chat.ServerMessage{
			Type: &chat.ServerMessage_Joined{Joined: msg},
		}, nil
	// пришло обычное сообщение
	case *chat.ChatMessage:
		return &chat.ServerMessage{
			Type: &chat.ServerMessage_SendMessage{SendMessage: msg},
		}, nil
	// п
	case *chat.UserLeft:
		return &chat.ServerMessage{
			Type: &chat.ServerMessage_Left{Left: msg},
		}, nil
	default:
		return nil, fmt.Errorf("%w: %T", ErrUnknownMessageType, message)
	}
}

// Отправка с таймаутом для предотвращения блокировки
func sendWithTimeout(stream chat.Chat_ChatStreamServer, msg *chat.ServerMessage, userID int64) error {
	// Создаем контекст с таймаутом для отправки
	ctx, cancel := context.WithTimeout(stream.Context(), SendTimeout)
	defer cancel()

	// Используем канал для неблокирующей отправки
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)
		errCh <- stream.Send(msg)
	}()

	select {
	case err := <-errCh:
		if err != nil { // вернулась ошибка отправки сообщения
			return ErrSendMessage
			//логируем ошибку через log.slog (err, userID, ошибку send)
		}
		return nil // отправка сообщения прошла успешна, метод Send вернул nil
	case <-ctx.Done():
		// определить ошибку ctx (deadline или cancel) и вернуть свою
		// логируем ошибку через log.slog
		return fmt.Errorf("send timeout for user %d: %w", userID, ctx.Err())
	}
}
