package chatdomain

import (
	"time"

	"gorm.io/gorm"
)

type IMessage interface {
	GetUserID() int64
	GetUsername() string
	GetTimestamp() int64
}

type Message struct {
	id           uint64         `gorm:"primarykey;column:id"`
	userID       int64          `gorm:"column:user_id"`
	username     string         `gorm:"column:username"`
	type_message string         `gorm:"column:type"`
	timestamp    time.Time      `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP;column:timestamp"`
	deletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}

func NewMessage(username string, userID int64) *Message {
	return &Message{
		userID:   userID,
		username: username,
	}
}

func (m *Message) GetUserID() int64 {
	return m.userID
}

func (m *Message) GetUsername() string {
	return m.username
}

func (m *Message) GetType() string {
	return m.type_message
}

func (m *Message) GetTimestamp() time.Time {
	return m.timestamp
}
