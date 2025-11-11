package chatgrpcserver

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/internal/chat/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IBusiness interface {
	Subscribe(msg *domain.Client) error
	EntryPoint(msg *chat.ClientMessage) error
	Unsubscribe(userID int64) error
}

type serverAPI struct {
	chat.UnimplementedChatServer
	Business IBusiness
}

func Register(grpc *grpc.Server, business IBusiness) {
	chat.RegisterChatServer(grpc, &serverAPI{Business: business})
}

func (s *serverAPI) ChatStream(stream chat.Chat_ChatStreamServer) error {

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

	// создаем клиента
	client := &domain.Client{
		ID:   joinMsg.Join.GetUserId(),
		Name: joinMsg.Join.GetUsername(),
		Buf:  make(chan interface{}, 100), // для получения сообщений
		Done: make(chan struct{}),         // для graceful shutdown
	}

	// Подписываем клиента
	if err := s.Business.Subscribe(client); err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	// Гарантируем отписку при выходе
	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			s.Business.Unsubscribe(client.ID)
			close(client.Done)
		})
	}
	defer unsubscribe()

	// Канал для ошибок из горутин
	errCh := make(chan error, 2)

	// Запускаем отправку сообщений клиенту
	go func() {
		defer unsubscribe() // гарантируем cleanup
		if err := AcceptMessages(client, stream); err != nil {
			log.Printf("AcceptMessages error for user %d: %v", client.ID, err)
		}
	}()

	// Обрабатываем входящие сообщения от клиента
	for {
		select {
		// case <-client.Done:
		// 	return status.Error(codes.ResourceExhausted, "the server is full, try again later")
		case err := <-errCh:
			return err // ошибка из горутин
		case <-ctx.Done():
			return ctx.Err() // клиент отключился
		default:
			clientMsg, err := stream.Recv()
			if err != nil {
				if err != nil {
					if errors.Is(err, io.EOF) {
						return status.Error(codes.OK, "client disconnected")
					}
					return fmt.Errorf("failed to receive message: %w", err)
				}
			}

			if err := s.Business.EntryPoint(clientMsg); err != nil {
				log.Printf("Business.EntryPoint error: %v", err)
				return fmt.Errorf("failed to enter chat: %w", err)
			}
		}
	}
}
