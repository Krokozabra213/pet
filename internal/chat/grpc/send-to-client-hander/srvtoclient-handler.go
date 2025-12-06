package sendtoclienthander

import (
	"time"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
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
	Messager IMessager
	Sender   IMessageSender
}

func New(
	stream chat.Chat_ChatStreamServer, buffer <-chan interface{}, shutdown <-chan struct{},
) *SendToClientHandler {

	messager := NewMessager(stream.Context(), buffer, shutdown)
	messageSender := NewMessageSender(SendTimeout, stream)

	return &SendToClientHandler{
		Messager: messager,
		Sender:   messageSender,
	}
}

func (handler *SendToClientHandler) Run() error { // будем выкидывать в канал errch ошибку

	defer func() {
		if r := recover(); r != nil {
			// логируем панику
		}
	}()

	for {
		// достаем сообщение из буфера
		message, err := handler.Messager.GetClientMessage()
		if err != nil {
			return err
		}

		// конвертируем в серверное сообщение
		serverMessage, err := ConvertMessage(message)
		if err != nil {
			return err
		}

		// отправляем клиенту stream.Send()
		err = handler.Sender.SendToClient(serverMessage)
		if err != nil {
			return err
		}
	}
}
