package chatgrpc

import (
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IBusiness interface {
	Subscribe(username string) (chan interface{}, chan struct{}, uint64, error)
	EntryPoint(msg *chat.ClientMessage) error
	Unsubscribe(userID int64) error
}

type ServerAPI struct {
	chat.UnimplementedChatServer
	Business IBusiness
	Log      *slog.Logger
}

func New(log *slog.Logger, business IBusiness) *ServerAPI {
	return &ServerAPI{
		Business: business,
		Log:      log,
	}
}

func (s *ServerAPI) ChatStream(stream chat.Chat_ChatStreamServer) error {
	const op = "chat.ChatStreamHandler"
	log := s.Log.With(
		slog.String("op", op),
	)

	ctx := stream.Context()
	if ctx.Err() != nil {
		return status.Error(codes.Aborted, ctx.Err().Error())
	}

	// Получаем первое сообщение - Join
	req, err := stream.Recv()
	if err != nil {
		log.Error("failed recv join msg", slog.String("error", err.Error()))
		if errors.Is(err, io.EOF) {
			return status.Error(codes.Canceled, ErrDisconect.Error())
		}
		return status.Error(codes.Internal, ErrStream.Error())
	}

	joinMsg, ok := req.Type.(*chat.ClientMessage_Join)
	if !ok {
		log.Info("failed parse joinMsg")
		return status.Error(codes.InvalidArgument, ErrFirstMessage.Error())
	}

	userID := joinMsg.Join.GetUserId()
	username := joinMsg.Join.GetUsername()

	buffer, done, userUUID, err := s.Business.Subscribe(username)
	if err != nil {
		log.Error("failed subscribe", slog.String("error", err.Error()))
		return status.Error(codes.Internal, err.Error())
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
