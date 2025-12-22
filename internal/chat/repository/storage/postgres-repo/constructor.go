package postgresrepo

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type IPG interface {
	WithContext(ctx context.Context) *gorm.DB
}

const (
	ctxTimeout = 3 * time.Second
	MaxRetries = 3
)

type Postgres struct {
	DB IPG
}

func New(db IPG) *Postgres {
	return &Postgres{
		DB: db,
	}
}
