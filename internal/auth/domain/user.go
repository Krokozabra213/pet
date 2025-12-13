package domain

import "gorm.io/gorm"

const UserEntity = "User"

type User struct {
	gorm.Model
	ID       uint64 `gorm:"primarykey"`
	Username string `gorm:"uniqueIndex:idx_username;not null"`
	Password string
}

func NewUser(username string, password string) *User {
	return &User{
		Username: username,
		Password: password,
	}
}
