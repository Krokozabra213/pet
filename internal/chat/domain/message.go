package chatdomain

import (
	"time"

	"gorm.io/gorm"
)

// type IMessage interface {
// 	GetUserID() int64
// 	GetUsername() string
// 	GetCreatedAt() time.Time
// 	GetType() uint8
// 	GetID() uint64
// }

type Message struct {
	ID        uint64         `gorm:"column:id"`
	UserID    int64          `gorm:"column:user_id"`
	Username  string         `gorm:"-"`
	Type      uint8          `gorm:"column:type"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func NewMessage(username string, userID int64, msgType uint8) *Message {
	return &Message{
		UserID:   userID,
		Username: username,
		Type:     msgType,
	}
}

func (m *Message) TableName() string {
	return "messages"
}

func (m *Message) GetID() uint64 {
	return m.ID
}

func (m *Message) GetUserID() int64 {
	return m.UserID
}

func (m *Message) GetUsername() string {
	return m.Username
}

func (m *Message) GetType() uint8 {
	return m.Type
}

func (m *Message) GetCreatedAt() time.Time {
	return m.CreatedAt
}
