package chatgrpc

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/internal/chat/domain"
	recvprocessor "github.com/Krokozabra213/sso/internal/chat/grpc/recv-processor"
	handlersfactory "github.com/Krokozabra213/sso/internal/chat/grpc/recv-processor/handlers-factory"
	sendtoclienthander "github.com/Krokozabra213/sso/internal/chat/grpc/send-to-client-hander"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ErrChanCap = 2
)

type IChatClient interface {
	GetUUID() uint64
	GetName() string
	GetBuffer() chan interface{}
	GetDone() chan struct{}
}

type IBusiness interface {
	Subscribe(ctx context.Context, username string) (IChatClient, error)
	SendMessage(ctx context.Context, msg *domain.DefaultMessage) error
	Unsubscribe(userID uint64)
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
	const op = "chat.Handler"
	log := s.Log.With(
		slog.String("op", op),
	)

	ctx := stream.Context()
	if ctx.Err() != nil {
		return HandleStreamContextError(ctx)
	}

	// Получаем первое сообщение - Join
	req, err := stream.Recv()
	if err != nil {
		err = ValidateStreamRecvErrors(err)
		return err
	}

	joinMessage, err := s.ValidateJoinMessage(req)
	if err != nil {
		return err
	}
	username := joinMessage.Join.GetUsername()
	userID := joinMessage.Join.GetUserId()

	// логируем сообщение
	msg, _ := json.MarshalIndent(joinMessage, "", "  ")
	log.Debug("debug message", slog.String("message", string(msg)))

	client, err := s.Business.Subscribe(ctx, username)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer s.Business.Unsubscribe(client.GetUUID())

	// Канал для ошибок из горутин
	errCh := make(chan error, ErrChanCap)

	// Запускаем отправку сообщений server -> client
	go func() {
		messageHandler := sendtoclienthander.New(stream, client.GetBuffer(), client.GetDone(), s.Log)
		if err := messageHandler.Run(); err != nil {
			select {
			case errCh <- err:
				// ошибка отправления
			default:
				return // если канал полон, то перестаём пытаться отправить сообщения
			}
		}
	}()

	factory := handlersfactory.New()
	factory.InitHandlers(s.Business, ctx, userID, username)
	processor := recvprocessor.New(s.Log, factory)

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
				s.Log.Error("recv message failed", "error", err)
				err = ValidateStreamRecvErrors(err)
				return err
			}

			if err := processor.Process(clientMsg); err != nil {
				return err
			}
		}
	}
}
