package recvprocessor

import (
	"fmt"
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/chat"
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
	log     *slog.Logger
	factory IHandlerFactory
}

func New(factory IHandlerFactory) *MessageProcessor {
	op := "chat.Handler.MessageProcessor"
	log := slog.With(
		slog.String("op", op),
	)

	return &MessageProcessor{
		log:     log,
		factory: factory,
	}
}

func (mp *MessageProcessor) Process(message *chat.ClientMessage) error {
	handler, err := mp.factory.GetHandler(message)
	if err != nil {
		mp.log.Error("failed to get handler",
			"err", err,
			"message type", fmt.Sprintf("%T", message))
		return status.Error(codes.InvalidArgument, err.Error())
	}

	// логировать тип handler

	if err := handler.Handle(message); err != nil {
		mp.log.Debug(err.Error())
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
