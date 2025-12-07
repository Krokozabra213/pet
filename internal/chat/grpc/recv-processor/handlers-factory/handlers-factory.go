package handlersfactory

import (
	"context"
	"errors"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/internal/chat/domain"
	recvprocessor "github.com/Krokozabra213/sso/internal/chat/grpc/recv-processor"
	"github.com/Krokozabra213/sso/internal/chat/grpc/recv-processor/handlers-factory/handlers"
)

var (
	ErrUnknownMessageType = errors.New("unknown message type")
)

// все методы бизнес логики всех хендлеров
type IHandlersBusiness interface {
	// добавлять по мере необходимости
	SendDefaultMessage(ctx context.Context, msg *domain.DefaultMessage) error
}

const (
	CountHandlers = 1
)

type MessageHandlerFactory struct {
	handlers []recvprocessor.IHandler
}

func New() *MessageHandlerFactory {
	return &MessageHandlerFactory{
		handlers: make([]recvprocessor.IHandler, 0, CountHandlers),
	}
}

func (factory *MessageHandlerFactory) GetHandler(message *chat.ClientMessage) (recvprocessor.IHandler, error) {
	for _, handler := range factory.handlers {
		if handler.CanHandle(message) {
			return handler, nil
		}
	}

	return nil, ErrUnknownMessageType
}

func (factory *MessageHandlerFactory) register(handler recvprocessor.IHandler) {
	factory.handlers = append(factory.handlers, handler)
}

func (factory *MessageHandlerFactory) InitHandlers(
	business IHandlersBusiness,
	ctx context.Context,
	userID int64,
	username string,
) {
	factory.register(handlers.NewSendMessageHandler(business, ctx, userID, username))
	// Добавляем другие обработчики по мере необходимости
}
