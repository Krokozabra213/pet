package handlers

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/chat"
	chatdomain "github.com/Krokozabra213/sso/internal/chat/domain"
	chatinterfaces "github.com/Krokozabra213/sso/internal/chat/grpc/interfaces"
)

type TextMessageHandler struct {
	Business chatinterfaces.ISendTextMessage
	Ctx      context.Context
	UserID   int64
	Username string
}

func NewTextMessageHandler(
	business chatinterfaces.ISendTextMessage,
	ctx context.Context,
	userID int64,
	username string,
) *TextMessageHandler {
	return &TextMessageHandler{
		Business: business,
		Ctx:      ctx,
		UserID:   userID,
		Username: username,
	}
}

func (handler *TextMessageHandler) GetType() string {
	return "TextMessage_Type"
}

func (handler *TextMessageHandler) CanHandle(message *chat.ClientMessage) bool {
	_, ok := message.Type.(*chat.ClientMessage_SendMessage)
	return ok
}

func (handler *TextMessageHandler) Handle(msg *chat.ClientMessage) error {
	textMsg := msg.Type.(*chat.ClientMessage_SendMessage)
	message := chatdomain.NewMessage(handler.Username, handler.UserID, TextType)
	text := chatdomain.NewText(textMsg.SendMessage.GetContent())
	textMessage := chatdomain.NewTextMessage(message, text)
	return handler.Business.SendTextMessage(handler.Ctx, textMessage)
}
