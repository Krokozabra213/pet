package chatgrpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
)

var (
	ErrSendMessage        = errors.New("error send message")
	ErrUnknownMessageType = errors.New("unknown message type")
)

const (
	SendTimeout = 10 * time.Second
)

func SendFromSrvToClient(buf chan interface{}, done chan struct{}, userID int64, stream chat.Chat_ChatStreamServer) error {

	for {
		select {
		case <-done:
			return nil
		case message, ok := <-buf:
			if !ok {
				// Канал закрыт - нормальное завершение
				return nil
			}

			serverMsg, err := convertToServerMessage(message)
			if err != nil {
				log.Println(err)
				continue
			}

			if err := sendWithTimeout(stream, serverMsg, userID); err != nil {
				return fmt.Errorf("failed to send message to client %d: %w", userID, err)
			}

		case <-stream.Context().Done():
			// Клиент отключился
			ctxErr := stream.Context().Err()
			log.Printf("Stream context done for user %d: %v", userID, ctxErr)
			return ctxErr
		}
	}
}

func convertToServerMessage(message interface{}) (*chat.ServerMessage, error) {
	switch msg := message.(type) {
	case *chat.UserJoined:
		return &chat.ServerMessage{
			Type: &chat.ServerMessage_Joined{Joined: msg},
		}, nil
	case *chat.ChatMessage:
		return &chat.ServerMessage{
			Type: &chat.ServerMessage_SendMessage{SendMessage: msg},
		}, nil
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
		if err != nil {
			return fmt.Errorf("send failed for user %d: %w", userID, err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("send timeout for user %d: %w", userID, ctx.Err())
	}
}
