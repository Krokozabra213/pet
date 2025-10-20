package domain

import "gorm.io/gorm"

type App struct {
	gorm.Model
	Name  string
	Sault string
}
