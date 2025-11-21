package domain

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	UserID int64 `gorm:"index"`
	User   User  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

const AdminEntity = "Admin"
