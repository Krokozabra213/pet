package postgres

import (
	"context"

	"github.com/Krokozabra213/sso/internal/auth/domain"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
	contexthandler "github.com/Krokozabra213/sso/pkg/context-handler"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
)

func (p *Postgres) AppByID(
	parentCtx context.Context, appID int,
) (*domain.App, error) {

	ctx, cancel := contexthandler.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return nil, storage.CtxError(ctx.Err())
	}

	var app domain.App
	result := p.DB.Client.WithContext(ctx).First(&app, "id = ?", appID)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		err := ErrorFactory(domain.AppEntity, customErr)
		return nil, err
	}
	return &app, nil
}
