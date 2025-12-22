package authusecases

import (
	"context"
	"errors"
	"log/slog"

	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	businessinput "github.com/Krokozabra213/sso/internal/auth/domain/business-input"
	businessoutput "github.com/Krokozabra213/sso/internal/auth/domain/business-output"
	"github.com/Krokozabra213/sso/internal/auth/lib/hmac"
	"github.com/Krokozabra213/sso/internal/auth/lib/jwt"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
)

func (a *Auth) Refresh(
	ctx context.Context, input *businessinput.RefreshInput,
) (*businessoutput.RefreshOutput, error) {
	const op = "auth.Refresh"

	refreshToken := input.GetRefreshToken()

	log := slog.With(
		slog.String("op", op),
		slog.String("token", refreshToken),
	)
	log.Info("starting refresh token process")

	// хешируем токен для безопасного хранения в базе
	hashToken := hmac.HashJWTTokenHMAC(refreshToken, a.cfg.Auth.AppSecretKey)

	exist, err := a.tokenRepo.CheckToken(ctx, hashToken)
	if err != nil {
		log.Error("failed to check token", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrTimeout)
		}
		return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrInternal)
	}
	if exist {
		return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrTokenRevoked)
	}

	// достаём claims из refresh токена
	claims, err := jwt.ParseRefresh(refreshToken, a.keyManager)
	if err != nil {
		log.Error("failed to parse token", slog.String("error", err.Error()))
		return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrParse)
	}

	// проверяем наличие пользователя с таким id
	user, err := a.userProvider.UserByID(ctx, int64(claims.UserID))
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

	// проверяем наличие приложения с таким id
	app, err := a.appProvider.AppByID(ctx, claims.AppID)
	if err != nil {
		log.Error("failed to get app by id", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return nil, authBusiness.BusinessError(domain.AppEntity, authBusiness.ErrTimeout)
		}
		if errors.Is(err, storage.ErrNotFound) {
			return nil, authBusiness.BusinessError(domain.AppEntity, authBusiness.ErrNotFound)
		}
		return nil, authBusiness.BusinessError(domain.AppEntity, authBusiness.ErrInternal)
	}

	// генерируем новую пару токенов
	tokenGen := jwt.New(
		user, app, a.cfg.Auth.JWT.AccessTokenTTL,
		a.cfg.Auth.JWT.RefreshTokenTTL, a.keyManager,
	)

	tokenPair, err := tokenGen.GenerateTokenPair()
	if err != nil {
		log.Error("failed to generate token", "err", err.Error())
		return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrTokenGenerate)
	}

	err = a.tokenRepo.SaveToken(ctx, hashToken, claims.Exp)
	if err != nil {
		log.Error("failed to save revoking token", "err", err.Error())
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrTimeout)
		}
		if errors.Is(err, storage.ErrTokenExpired) {
			return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrTokenExpired)
		}
		return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrInternal)
	}
	log.Info("user refreshed token", "tokens", tokenPair.AccessToken+", "+tokenPair.RefreshToken)
	uotput := businessoutput.NewRefreshOutput(tokenPair.AccessToken, tokenPair.RefreshToken)
	return uotput, nil
}
