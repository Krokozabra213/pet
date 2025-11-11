package chatgrpcserver

import (
	"fmt"
	"log"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/internal/chat/domain"
)

func AcceptMessages(cli *domain.Client, stream chat.Chat_ChatStreamServer) error {

	for {
		select {
		case <-cli.Done:
			return nil
		case message, ok := <-cli.Buf:
			if !ok {
				// Канал закрыт - нормальное завершение
				return nil
			}

			var serverMsg *chat.ServerMessage
			switch msg := message.(type) {
			case *chat.UserJoined:
				serverMsg = &chat.ServerMessage{
					Type: &chat.ServerMessage_Joined{Joined: msg},
				}
			case *chat.ChatMessage:
				serverMsg = &chat.ServerMessage{
					Type: &chat.ServerMessage_SendMessage{SendMessage: msg},
				}
			case *chat.UserLeft:
				serverMsg = &chat.ServerMessage{
					Type: &chat.ServerMessage_Left{Left: msg},
				}
			default:
				log.Printf("Unknown message type: %T", message)
				continue // Пропускаем неизвестные сообщения
			}

			if err := stream.Send(serverMsg); err != nil {
				log.Printf("Failed to send message to client %v: %v", cli.ID, err)
				return fmt.Errorf("error send message to client: %w", err) // Возвращаем ошибку!
			}

		case <-stream.Context().Done():
			// Клиент отключился
			return stream.Context().Err()
		}
	}
}
