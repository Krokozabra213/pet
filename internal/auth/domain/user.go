package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex:idx_username;not null"`
	Password string
}
