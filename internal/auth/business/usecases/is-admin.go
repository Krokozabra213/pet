package authusecases

import (
	"context"
	"errors"
	"log/slog"

	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
)

func (a *Auth) IsAdmin(
	ctx context.Context, input *domain.IsAdminInput,
) (*domain.IsAdminOutput, error) {
	const op = "auth.IsAdmin"

	userID := input.GetUserID()

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)
	log.Info("starting isadmin process")

	_, err := a.userProvider.UserByID(ctx, userID)
	if err != nil {
		log.Error("failed to get user by id", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return nil, authBusiness.BusinessError(domain.UserEntity, authBusiness.ErrTimeout)
		}
		if errors.Is(err, storage.ErrNotFound) {
			return nil, authBusiness.BusinessError(domain.UserEntity, authBusiness.ErrNotFound)
		}
		return nil, authBusiness.BusinessError(domain.UserEntity, authBusiness.ErrInternal)
	}

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		log.Error("failed checking isadmin user", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return nil, authBusiness.BusinessError(domain.AdminEntity, authBusiness.ErrTimeout)
		}
		return nil, authBusiness.BusinessError(domain.AdminEntity, authBusiness.ErrInternal)
	}

	if !isAdmin {
		log.Info("err permission")
		return nil, authBusiness.BusinessError(domain.AdminEntity, authBusiness.ErrPermission)
	}

	log.Info("ending isadmin process")
	uotput := domain.NewIsAdminOutput(isAdmin)
	return uotput, nil
}
