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

	err = p.DB.WithContext(ctx).Create(token).Error
	if err != nil {
		customErr := postgrespet.ErrorWrapper(err)
		repoErr := ErrorFactory(domain.TokenEntity, customErr)
		return repoErr
	}

	return nil
}
