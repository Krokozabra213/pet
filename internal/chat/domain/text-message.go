package chatdomain

import (
	"time"

	"github.com/Krokozabra213/protos/gen/go/chat"
)

type ITextMessage interface {
	GetTimestamp() time.Time
	IServerMessage
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

func (tm *TextMessage) GetMessage() *Message {
	return tm.message
}

func (tm *TextMessage) GetText() *Text {
	return tm.text
}

func (tm *TextMessage) GetUserID() int64 {
	return tm.message.GetUserID()
}

func (tm *TextMessage) GetUsername() string {
	return tm.message.GetUsername()
}

func (tm *TextMessage) GetTimestamp() time.Time {
	return tm.message.GetTimestamp()
}

func (tm *TextMessage) GetContent() string {
	return tm.text.GetContent()
}

func (tm *TextMessage) ConvertToServerMessage() *chat.ServerMessage {

	timestamp := tm.GetTimestamp()
	protoTimestamp := timeToProto(timestamp)

	return &chat.ServerMessage{
		Type: &chat.ServerMessage_SendMessage{SendMessage: &chat.ChatMessage{
			Message: &chat.Message{
				UserId:    tm.GetUserID(),
				Timestamp: protoTimestamp,
			},
			Content: tm.GetContent(),
		}},
	}
}
