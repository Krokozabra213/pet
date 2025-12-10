package authusecases

import (
	"context"
	"errors"
	"log/slog"

	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
)

func (a *Auth) PublicKey(
	ctx context.Context, input *domain.PublicKeyInput,
) (*domain.PublicKeyOutput, error) {
	const op = "auth.PublicKey"

	appID := input.GetAppID()

	log := a.log.With(
		slog.String("op", op),
		slog.Int("app_id", appID),
	)
	log.Info("starting publickey process")

	_, err := a.appProvider.AppByID(ctx, appID)
	if err != nil {
		log.Error("failed get app by id", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return nil, authBusiness.BusinessError(domain.AppEntity, authBusiness.ErrTimeout)
		}
		if errors.Is(err, storage.ErrNotFound) {
			return nil, authBusiness.BusinessError(domain.AppEntity, authBusiness.ErrNotFound)
		}
		return nil, authBusiness.BusinessError(domain.AppEntity, authBusiness.ErrInternal)
	}

	log.Info("ending publickey process")
	publicKey, err := a.keyManager.GetPublicKeyPEM()
	uotput := domain.NewPublicKeyOutput(publicKey)
	return uotput, err
}
