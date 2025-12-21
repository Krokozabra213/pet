package chatdomain

import (
	"github.com/Krokozabra213/protos/gen/go/chat"
)

type ITextMessage interface {
	// IServerMessage
	IConvertServerMessage
}

type TextMessage struct {
	message *Message
	text    *Text
}

func NewTextMessage(message *Message, text *Text) *TextMessage {
	return &TextMessage{
		message: message,
		text:    text,
	}
}

func (t *TextMessage) GetMessage() *Message {
	return t.message
}

func (t *TextMessage) GetText() *Text {
	return t.text
}

func (t *TextMessage) ConvertToServerMessage() *chat.ServerMessage {

	timestamp := t.message.GetCreatedAt()
	protoTimestamp := timeToProto(timestamp)

	return &chat.ServerMessage{
		Type: &chat.ServerMessage_SendMessage{SendMessage: &chat.ChatMessage{
			Message: &chat.Message{
				UserId:    t.message.GetUserID(),
				Timestamp: protoTimestamp,
			},
			Content: t.text.GetContent(),
		}},
	}
}
