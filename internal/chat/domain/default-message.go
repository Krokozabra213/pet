package domain

import (
	"time"

	"github.com/Krokozabra213/protos/gen/go/chat"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

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
	ID        uint64    `gorm:"primarykey"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP" json:"timestamp"`
	DeletedAt gorm.DeletedAt
}

func NewDefaultMessage(content string, username string, userID int64) *DefaultMessage {
	return &DefaultMessage{
		UserID:   userID,
		Username: username,
		Content:  content,
	}
}

func (m *DefaultMessage) ConvertToServerMessage() *chat.ServerMessage {

	timestamp := m.GetTimestamp()

	return &chat.ServerMessage{
		Type: &chat.ServerMessage_SendMessage{SendMessage: &chat.ChatMessage{
			UserId:    m.GetUserID(),
			Username:  m.GetUsername(),
			Content:   m.GetContent(),
			Timestamp: timeToProto(timestamp),
		}},
	}
}

func timeToProto(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
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

func (m *DefaultMessage) GetTimestamp() time.Time {
	return m.Timestamp
}
