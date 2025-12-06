package domain

import "github.com/Krokozabra213/protos/gen/go/proto/chat"

type IServerMessage interface {
	ConvertToServerMessage() *chat.ServerMessage
}

type IDefaultMessage interface {
	GetUserID() int64
	GetUsername() string
	GetContent() string
	GetTimestamp() int64
}

type DefaultMessage struct {
	UserID    int64
	Username  string
	Content   string
	Timestamp int64
}

func NewDefaultMessage(content string, username string, userID int64) *DefaultMessage {
	return &DefaultMessage{
		UserID:   userID,
		Username: username,
		Content:  content,
	}
}

func (m *DefaultMessage) ConvertToServerMessage() *chat.ServerMessage {
	return &chat.ServerMessage{
		Type: &chat.ServerMessage_SendMessage{SendMessage: &chat.ChatMessage{
			UserId:    m.GetUserID(),
			Username:  m.GetUsername(),
			Content:   m.GetContent(),
			Timestamp: m.GetTimestamp(),
		}},
	}
}

func (m *DefaultMessage) GetUserID() int64 {
	return m.UserID
}

func (m *DefaultMessage) GetUsername() string {
	return m.Username
}

func (m *DefaultMessage) GetContent() string {
	return m.Content
}

func (m *DefaultMessage) GetTimestamp() int64 {
	return m.Timestamp
}
