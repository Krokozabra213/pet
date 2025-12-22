package postgrespet

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PGDB struct {
	*gorm.DB
}

func NewPGDB(dsn string) *PGDB {
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(dsn)
		panic(err)
	}
	return &PGDB{client}
}
