package authusecases

import (
	"context"
	"errors"
	"log/slog"
	"time"

	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	businessinput "github.com/Krokozabra213/sso/internal/auth/domain/business-input"
	businessoutput "github.com/Krokozabra213/sso/internal/auth/domain/business-output"
	"github.com/Krokozabra213/sso/internal/auth/lib/jwt"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
	"golang.org/x/crypto/bcrypt"
)

func (a *Auth) Login(
	ctx context.Context, input *businessinput.LoginInput,
) (*businessoutput.LoginOutput, error) {
	const op = "auth.Login"

	username := input.GetUsername()
	password := input.GetPassword()
	appID := input.GetAppID()

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
		slog.Int("app_id", appID),
	)
	log.Info("starting user logining process")

	app, err := a.appProvider.AppByID(ctx, appID)
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

	user, err := a.userProvider.User(ctx, username)
	if err != nil {
		log.Error("failed get user", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return nil, authBusiness.BusinessError(domain.UserEntity, authBusiness.ErrTimeout)
		}
		if errors.Is(err, storage.ErrNotFound) {
			return nil, authBusiness.BusinessError(domain.UserEntity, authBusiness.ErrNotFound)
		}
		return nil, authBusiness.BusinessError(domain.UserEntity, authBusiness.ErrInternal)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Error("failed compare passwords", slog.String("error", err.Error()))
		return nil, authBusiness.BusinessError(domain.UserEntity, authBusiness.ErrInvalidCredentials)
	}

	tokenGen := jwt.New(
		user, app, time.Duration(a.cfg.Security.AccessTokenTTL),
		time.Duration(a.cfg.Security.RefreshTokenTTL), a.keyManager,
	)

	tokenPair, err := tokenGen.GenerateTokenPair()
	if err != nil {
		log.Error("failed generate token", slog.String("error", err.Error()))
		return nil, authBusiness.BusinessError(domain.TokenEntity, authBusiness.ErrTokenGenerate)
	}

	log.Info("user successfully logining")
	uotput := businessoutput.NewLoginOutput(tokenPair.AccessToken, tokenPair.RefreshToken)
	return uotput, nil
}
