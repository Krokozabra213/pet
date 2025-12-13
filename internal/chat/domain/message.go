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
	ID        uint64         `gorm:"primarykey;column:id"`
	UserID    int64          `gorm:"column:user_id"`
	Username  string         `gorm:"-"`
	Type      string         `gorm:"column:type"`
	Timestamp time.Time      `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP;column:timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func NewMessage(username string, userID int64) *Message {
	return &Message{
		UserID:   userID,
		Username: username,
	}
}

func (m *Message) GetUserID() int64 {
	return m.UserID
}

func (m *Message) GetUsername() string {
	return m.Username
}

func (m *Message) GetType() string {
	return m.Type
}

func (m *Message) GetTimestamp() time.Time {
	return m.Timestamp
}
