package handlers

import (
	"context"

	"github.com/Krokozabra213/protos/gen/go/chat"
	chatdomain "github.com/Krokozabra213/sso/internal/chat/domain"
	chatinterfaces "github.com/Krokozabra213/sso/internal/chat/grpc/interfaces"
)

type ImageMessageHandler struct {
	Business chatinterfaces.ISendImageMessage
	Ctx      context.Context
	UserID   int64
	Username string
}

func NewImageMessageHandler(
	business chatinterfaces.ISendImageMessage,
	ctx context.Context,
	userID int64,
	username string,
) *ImageMessageHandler {
	return &ImageMessageHandler{
		Business: business,
		Ctx:      ctx,
		UserID:   userID,
		Username: username,
	}
}

func (handler *ImageMessageHandler) GetType() string {
	return "ImageMessage_Type"
}

func (handler *ImageMessageHandler) CanHandle(message *chat.ClientMessage) bool {
	_, ok := message.Type.(*chat.ClientMessage_ImageMessage)
	return ok
}

func (handler *ImageMessageHandler) Handle(msg *chat.ClientMessage) error {
	imgMsg := msg.Type.(*chat.ClientMessage_ImageMessage)
	message := chatdomain.NewMessage(handler.Username, handler.UserID)
	image := chatdomain.NewImage(imgMsg.ImageMessage.GetImageUrl())
	imageMessage := chatdomain.NewImageMessage(message, image)
	return handler.Business.SendImageMessage(handler.Ctx, imageMessage)
}
