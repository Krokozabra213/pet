package postgrespet

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PGDB struct {
	Client *gorm.DB
}

func NewPGDB(dsn string) *PGDB {
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &PGDB{Client: client}
}
