package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type IPG interface {
	WithContext(ctx context.Context) *gorm.DB
}

const (
	ctxTimeout = 5 * time.Second
)

type Postgres struct {
	DB IPG
}

func New(db IPG) *Postgres {
	return &Postgres{
		DB: db,
	}
}
