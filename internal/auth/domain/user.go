package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"index"`
	Password string
}
