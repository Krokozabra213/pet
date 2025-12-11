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
