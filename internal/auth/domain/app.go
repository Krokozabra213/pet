package domain

import "gorm.io/gorm"

type App struct {
	gorm.Model
	Name string
}

func NewApp(name string) *App {
	return &App{
		Name: name,
	}
}

const AppEntity = "App"
