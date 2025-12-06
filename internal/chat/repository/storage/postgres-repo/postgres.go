package postgresrepo

import (
	"context"
	"time"

	"github.com/Krokozabra213/sso/internal/chat/domain"
	"github.com/Krokozabra213/sso/internal/chat/repository/storage"
	contexthandler "github.com/Krokozabra213/sso/pkg/context-handler"
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

func (p *Postgres) SaveDefaultMessage(
	parentCtx context.Context, message *domain.DefaultMessage,
) (*domain.DefaultMessage, error) {

	ctx, cancel := contexthandler.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return nil, storage.CtxError(ctx.Err())
	}

	// TODO: ADD SAVE MESSAGE WITH TIMESTAMP RETURN

	return message, nil
}
