package postgres

import (
	"context"

	"github.com/Krokozabra213/sso/internal/auth/domain"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
	contexthandler "github.com/Krokozabra213/sso/pkg/context-handler"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
)

func (p *Postgres) SaveToken(
	parentCtx context.Context, token *domain.BlackToken,
) (err error) {

	ctx, cancel := contexthandler.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return storage.CtxError(ctx.Err())
	}

	result := p.DB.Client.WithContext(ctx).Create(token)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		err := ErrorFactory(domain.TokenEntity, customErr)
		return err
	}
	return nil
}
