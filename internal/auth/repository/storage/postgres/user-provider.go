package postgres

import (
	"context"

	"github.com/Krokozabra213/sso/internal/auth/domain"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
	contexthandler "github.com/Krokozabra213/sso/pkg/context-handler"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
)

func (p *Postgres) SaveUser(
	parentCtx context.Context, user *domain.User,
) (uid uint64, err error) {

	ctx, cancel := contexthandler.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return 0, storage.CtxError(ctx.Err())
	}

	result := p.DB.Client.WithContext(ctx).Create(user)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		err := ErrorFactory(domain.UserEntity, customErr)
		return 0, err
	}

	return user.ID, nil
}

func (p *Postgres) User(
	parentCtx context.Context, username string,
) (*domain.User, error) {

	ctx, cancel := contexthandler.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return nil, storage.CtxError(ctx.Err())
	}

	var user domain.User
	result := p.DB.Client.WithContext(ctx).First(&user, "username = ?", username)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		err := ErrorFactory(domain.UserEntity, customErr)
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) UserByID(
	parentCtx context.Context, id int64,
) (*domain.User, error) {

	ctx, cancel := contexthandler.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return nil, storage.CtxError(ctx.Err())
	}

	var user domain.User
	result := p.DB.Client.WithContext(ctx).First(&user, "id = ?", id)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		err := ErrorFactory(domain.UserEntity, customErr)
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) IsAdmin(
	parentCtx context.Context, userID int64,
) (bool, error) {

	ctx, cancel := contexthandler.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return false, storage.CtxError(ctx.Err())
	}

	var exists bool
	result := p.DB.Client.WithContext(ctx).Raw(
		"SELECT EXISTS(SELECT 1 FROM admins WHERE user_id = ?)",
		userID,
	).Scan(&exists)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		err := ErrorFactory(domain.AdminEntity, customErr)
		return false, err
	}
	return exists, nil
}
