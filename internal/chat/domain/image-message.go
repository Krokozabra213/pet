package chatdomain

import (
	"github.com/Krokozabra213/protos/gen/go/chat"
)

type IImageMessage interface {
	// GetImageUrl() string
	// IServerMessage
	IConvertServerMessage
}

type ImageMessage struct {
	message *Message
	image   *Image
}

func NewImageMessage(message *Message, image *Image) *ImageMessage {
	return &ImageMessage{
		message: message,
		image:   image,
	}
}

func (i *ImageMessage) GetMessage() *Message {
	return i.message
}

func (i *ImageMessage) GetImage() *Image {
	return i.image
}

func (i *ImageMessage) ConvertToServerMessage() *chat.ServerMessage {

	timestamp := i.message.GetCreatedAt()
	protoTimestamp := timeToProto(timestamp)

	return &chat.ServerMessage{
		Type: &chat.ServerMessage_ImgMessage{ImgMessage: &chat.ImgMessage{
			Message: &chat.Message{
				UserId:    i.message.GetUserID(),
				Timestamp: protoTimestamp,
			},
			ImageUrl: i.image.GetImageUrl(),
		}},
	}
}
