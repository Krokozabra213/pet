package domain

import "gorm.io/gorm"

const AppEntity = "App"

type App struct {
	gorm.Model
	Name string
}

func NewApp(name string) *App {
	return &App{
		Name: name,
	}
}

type PublicKeyInput struct {
	appID int
}

func NewPublicKeyInput(appID int) *PublicKeyInput {
	return &PublicKeyInput{
		appID: appID,
	}
}

func (input *PublicKeyInput) GetAppID() int {
	return input.appID
}

type PublicKeyOutput struct {
	publicKey string
}

func NewPublicKeyOutput(publicKey string) *PublicKeyOutput {
	return &PublicKeyOutput{
		publicKey: publicKey,
	}
}

func (input *PublicKeyOutput) GetPublicKey() string {
	return input.publicKey
}
