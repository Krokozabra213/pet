package chatgrpc

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ErrChanCap = 2
)

type IBusiness interface {
	Subscribe(ctx context.Context, username string) (<-chan interface{}, chan struct{}, uint64, error)
	SendMessage(ctx context.Context, cliSendMsg *chat.ClientMessage) error
	// EntryPoint(msg *chat.ClientMessage) error
	Unsubscribe(userID uint64) error
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

	msg, _ := json.MarshalIndent(joinMsg, "", "  ")
	log.Info("message", slog.String("message", string(msg)))

	username := joinMsg.Join.GetUsername()

	buffer, done, userUUID, err := s.Business.Subscribe(ctx, username)
	if err != nil {
		log.Error("failed subscribe", slog.String("error", err.Error()))
		return status.Error(codes.Internal, err.Error())
	}
	defer s.Business.Unsubscribe(userUUID)

	// Канал для ошибок из горутин
	errCh := make(chan error, ErrChanCap)

	// Запускаем отправку сообщений server -> client
	go func() {
		if err := SendFromSrvToClient(buffer, done, joinMsg.Join.GetUserId(), stream); err != nil {
			log.Error("failed to send message to client", "error", err)
			select {
			case errCh <- err:
				// ошибка отправления
			default:
				return // если канал полон, то перестаём пытаться отправить сообщения
			}
		}
	}()

	// Обрабатываем входящие сообщения от клиента
	for {
		select {
		case err := <-errCh: // возникла ошибка отправки сообщений клиенту
			return err
		case <-ctx.Done():
			return ctx.Err() // клиент отключился
		default:
			clientMsg, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return status.Error(codes.OK, ErrDisconect.Error())
				}
				log.Error("failed to receive message", "err", err)
				return ErrRecvMessage
			}

			if err := s.Business.SendMessage(ctx, clientMsg); err != nil {
				log.Error("failed send message", "err", err)
				return err
			}
		}
	}
}
