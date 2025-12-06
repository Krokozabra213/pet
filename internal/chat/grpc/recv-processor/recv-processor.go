package recvprocessor

import (
	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IHandler interface {
	CanHandle(message *chat.ClientMessage) bool
	Handle(message *chat.ClientMessage) error
	GetType() string
}

type IHandlerFactory interface {
	GetHandler(message *chat.ClientMessage) (IHandler, error)
}

type MessageProcessor struct {
	// TODO: ADD LOGGER
	factory IHandlerFactory
}

func New(factory IHandlerFactory) *MessageProcessor {
	return &MessageProcessor{
		factory: factory,
	}
}

func (mp *MessageProcessor) Process(message *chat.ClientMessage) error {
	handler, err := mp.factory.GetHandler(message)
	if err != nil {
		// логировать ошибку и тип сообщения
		return status.Error(codes.InvalidArgument, err.Error())
	}

	// логировать тип handler

	if err := handler.Handle(message); err != nil {
		// логируем ошибку
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
