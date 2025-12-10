package authusecases

import (
	"context"
	"errors"
	"log/slog"

	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	"github.com/Krokozabra213/sso/internal/auth/lib/hmac"
	"github.com/Krokozabra213/sso/internal/auth/lib/jwt"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
)

func (a *Auth) Logout(
	ctx context.Context, input *domain.LogoutInput,
) (*domain.LogoutOutput, error) {
	const op = "auth.Logout"

	refreshToken := input.GetRefreshToken()

	log := a.log.With(
		slog.String("op", op),
		slog.String("token", refreshToken),
	)
	log.Info("starting user logouting process")

	claims, err := jwt.ParseRefresh(refreshToken, a.keyManager)
	if err != nil {
		log.Error("failed parsing token", slog.String("error", err.Error()))
		return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrParse)
	}

	hashToken := hmac.HashJWTTokenHMAC(refreshToken, a.cfg.Security.Secret)

	err = a.tokenRepo.SaveToken(ctx, hashToken, claims.Exp)
	if err != nil {
		log.Error("failed revoking token", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrTimeout)
		}
		if errors.Is(err, storage.ErrTokenExpired) {
			return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrTokenExpired)
		}
		return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrInternal)
	}
	log.Info("user successfully logouting")
	uotput := domain.NewLogoutOutput(true)
	return uotput, nil
}
