package postgresrepo

import (
	"time"

	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
)

const (
	ctxTimeout = 3 * time.Second
)

type Postgres struct {
	DB *postgrespet.PGDB
}

func New(db *postgrespet.PGDB) *Postgres {
	return &Postgres{
		DB: db,
	}
}
