package handlers

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/chat"
	"github.com/Krokozabra213/sso/internal/chat/domain"
)

type ISendDefaultMessage interface {
	SendDefaultMessage(ctx context.Context, msg *domain.DefaultMessage) error
}

type SendMessageHandler struct {
	Business ISendDefaultMessage
	Ctx      context.Context
	UserID   int64
	Username string
}

func NewSendMessageHandler(
	business ISendDefaultMessage,
	ctx context.Context,
	userID int64,
	username string,
) *SendMessageHandler {
	return &SendMessageHandler{
		Business: business,
		Ctx:      ctx,
		UserID:   userID,
		Username: username,
	}
}

func (handler *SendMessageHandler) GetType() string {
	return "SendMessageType"
}

func (handler *SendMessageHandler) CanHandle(message *chat.ClientMessage) bool {
	_, ok := message.Type.(*chat.ClientMessage_SendMessage)
	return ok
}

func (handler *SendMessageHandler) Handle(message *chat.ClientMessage) error {
	sendMsg := message.Type.(*chat.ClientMessage_SendMessage)
	defaultMessage := domain.NewDefaultMessage(
		sendMsg.SendMessage.GetContent(),
		handler.Username,
		handler.UserID,
	)
	return handler.Business.SendDefaultMessage(handler.Ctx, defaultMessage)
}
