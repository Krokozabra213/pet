package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	DB *gorm.DB
}

func NewPGDb(dsn string) *Db {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Db{DB: db}
}
