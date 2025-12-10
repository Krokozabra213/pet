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

type RegisterInput struct {
	username string
	password string
}

func NewRegisterInput(username string, password string) *RegisterInput {
	return &RegisterInput{
		username: username,
		password: password,
	}
}

func (input *RegisterInput) GetUsername() string {
	return input.username
}

func (input *RegisterInput) GetPassword() string {
	return input.password
}

type LoginInput struct {
	username string
	password string
	appID    int
}

func NewLoginInput(username string, password string, appID int) *LoginInput {
	return &LoginInput{
		username: username,
		password: password,
		appID:    appID,
	}
}

func (input *LoginInput) GetUsername() string {
	return input.username
}

func (input *LoginInput) GetPassword() string {
	return input.password
}

func (input *LoginInput) GetAppID() int {
	return input.appID
}

type RegisterOutput struct {
	userID uint64
}

func NewRegisterOutput(userID uint64) *RegisterOutput {
	return &RegisterOutput{
		userID: userID,
	}
}

func (input *RegisterOutput) GetUserID() uint64 {
	return input.userID
}

type LogoutOutput struct {
	success bool
}

func NewLogoutOutput(success bool) *LogoutOutput {
	return &LogoutOutput{
		success: success,
	}
}

func (input *LogoutOutput) GetSuccess() bool {
	return input.success
}
