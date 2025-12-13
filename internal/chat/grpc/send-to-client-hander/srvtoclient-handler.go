package sendtoclienthander

import (
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/Krokozabra213/protos/gen/go/chat"
	"github.com/Krokozabra213/sso/internal/chat/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	SendTimeout = 3 * time.Second
)

type IMessageSender interface {
	SendToClient(message *chat.ServerMessage) error
}

type IMessager interface {
	GetClientMessage() (interface{}, error)
}

type SendToClientHandler struct {
	// TODO: ADD LOGGER
	Log      *slog.Logger
	Messager IMessager
	Sender   IMessageSender
}

func New(
	stream chat.Chat_ChatStreamServer, buffer <-chan interface{},
	shutdown <-chan struct{}, log *slog.Logger,
) *SendToClientHandler {
	op := "chat.Handler.SendToClient"
	log = log.With(
		slog.String("op", op),
	)

	messager := NewMessager(log, stream.Context(), buffer, shutdown)
	messageSender := NewMessageSender(log, SendTimeout, stream)

	return &SendToClientHandler{
		Log:      log,
		Messager: messager,
		Sender:   messageSender,
	}
}

func (handler *SendToClientHandler) Run() (err error) {

	defer func() {
		if r := recover(); r != nil {
			handler.Log.Error("panic recovered in SendToClientHandler",
				slog.String("error", fmt.Sprintf("%v", r)),
				slog.String("stack", string(debug.Stack())),
				slog.String("handler_type", fmt.Sprintf("%T", handler)),
			)
			err = fmt.Errorf("handler panic: %v", r)
		}
	}()

	for {
		// достаем сообщение из буфера
		message, err := handler.Messager.GetClientMessage()
		if err != nil {
			return err
		}

		var serverMessage *chat.ServerMessage
		if msg, ok := message.(chatdomain.IConvertServerMessage); ok {
			serverMessage = msg.ConvertToServerMessage()
		} else {
			return status.Error(codes.InvalidArgument, ErrUnknownMessageType.Error())
		}

		err = handler.Sender.SendToClient(serverMessage)
		if err != nil {
			return err
		}
	}
}
