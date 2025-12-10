package domain

import "gorm.io/gorm"

const AdminEntity = "Admin"

type Admin struct {
	gorm.Model
	UserID int64 `gorm:"index"`
	User   User  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type IsAdminInput struct {
	userID int64
}

func NewIsAdminInput(userID int64) *IsAdminInput {
	return &IsAdminInput{
		userID: userID,
	}
}

func (input *IsAdminInput) GetUserID() int64 {
	return input.userID
}

type IsAdminOutput struct {
	access bool
}

func NewIsAdminOutput(access bool) *IsAdminOutput {
	return &IsAdminOutput{
		access: access,
	}
}

func (input *IsAdminOutput) GetAccess() bool {
	return input.access
}
