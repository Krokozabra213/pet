package chatgrpc

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IBusiness interface {
	Subscribe(userID int64, username string) (chan interface{}, chan struct{}, error)
	EntryPoint(msg *chat.ClientMessage) error
	Unsubscribe(userID int64) error
}

type ServerAPI struct {
	chat.UnimplementedChatServer
	Business IBusiness
}

func New(business IBusiness) *ServerAPI {
	return &ServerAPI{
		Business: business,
	}
}

func (s *ServerAPI) ChatStream(stream chat.Chat_ChatStreamServer) error {

	fmt.Println("пришел запрос")
	ctx := stream.Context()

	// Получаем первое сообщение - Join
	req, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("failed to receive initial message: %w", err)
	}

	joinMsg, ok := req.Type.(*chat.ClientMessage_Join)
	if !ok {
		return status.Error(codes.InvalidArgument, "first message must be Join")
	}

	userID := joinMsg.Join.GetUserId()
	buffer, done, err := s.Business.Subscribe(userID, joinMsg.Join.GetUsername())
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}
	defer close(done)
	defer close(buffer)
	defer s.Business.Unsubscribe(joinMsg.Join.GetUserId())

	// Канал для ошибок из горутин
	errCh := make(chan error, 2)

	// Запускаем отправку сообщений server -> client
	go func() {
		if err := SendFromSrvToClient(buffer, done, joinMsg.Join.GetUserId(), stream); err != nil {
			select {
			case errCh <- err:
				// ошибка отправления
			default:
				// если errCh полон, логируем и продолжаем
				log.Printf("Failed to send error to channel for user %d: %v", userID, err)
			}
		}
	}()

	// Обрабатываем входящие сообщения от клиента
	for {
		select {
		case err := <-errCh:
			if err != nil {
				log.Printf("Send goroutine error for user %d: %v", userID, err)
				return fmt.Errorf("send operation failed: %w", err)
			}
		case <-ctx.Done():
			return ctx.Err() // клиент отключился
		default:
			clientMsg, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return status.Error(codes.OK, "client disconnected")
				}
				return fmt.Errorf("failed to receive message: %w", err)
			}

			if err := s.Business.EntryPoint(clientMsg); err != nil {
				log.Printf("Business.EntryPoint error: %v", err)
				return fmt.Errorf("failed to enter chat: %w", err)
			}
		}
	}
}
