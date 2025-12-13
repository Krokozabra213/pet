package chatdomain

import (
	"time"

	"github.com/Krokozabra213/protos/gen/go/chat"
)

type IImageMessage interface {
	GetImageUrl() string
	IServerMessage
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

func (im *ImageMessage) GetMessage() *Message {
	return im.message
}

func (im *ImageMessage) GetImage() *Image {
	return im.image
}

func (im *ImageMessage) GetUserID() int64 {
	return im.message.GetUserID()
}

func (im *ImageMessage) GetUsername() string {
	return im.message.GetUsername()
}

func (im *ImageMessage) GetTimestamp() time.Time {
	return im.message.GetTimestamp()
}

func (im *ImageMessage) GetImageUrl() string {
	return im.image.GetImageUrl()
}

func (im *ImageMessage) ConvertToServerMessage() *chat.ServerMessage {

	timestamp := im.GetTimestamp()
	protoTimestamp := timeToProto(timestamp)

	return &chat.ServerMessage{
		Type: &chat.ServerMessage_ImgMessage{ImgMessage: &chat.ImgMessage{
			Message: &chat.Message{
				UserId:    im.GetUserID(),
				Timestamp: protoTimestamp,
			},
			ImageUrl: im.GetImageUrl(),
		}},
	}
}
