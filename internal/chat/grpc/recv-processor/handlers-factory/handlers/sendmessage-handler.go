package handlers

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/internal/chat/domain"
)

type ISendMessageBusiness interface {
	SendMessage(ctx context.Context, msg *domain.DefaultMessage) error
}

type SendMessageHandler struct {
	Business ISendMessageBusiness
	Ctx      context.Context
	UserID   int64
	Username string
}

func NewSendMessageHandler(
	business ISendMessageBusiness,
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
	return handler.Business.SendMessage(handler.Ctx, defaultMessage)
}
