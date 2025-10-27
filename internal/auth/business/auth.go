package authBusiness

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Krokozabra213/sso/configs/ssoconfig"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	"github.com/Krokozabra213/sso/internal/auth/lib/hmac"
	"github.com/Krokozabra213/sso/internal/auth/lib/jwt"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	empty = ""
)

type IUserSaver interface {
	SaveUser(ctx context.Context, username string, pass string) (uid uint, err error)
}

type ITokenRepo interface {
	SaveToken(ctx context.Context, hashToken string, exp time.Time) (err error)
	CheckToken(ctx context.Context, hashToken string, exp time.Time) (err error)
}

type IUserProvider interface {
	User(ctx context.Context, username string) (*domain.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type IAppProvider interface {
	App(ctx context.Context, appID int) (*domain.App, error)
}

type IKeyManager interface {
	GetPrivateKey() *rsa.PrivateKey
	GetPublicKey() *rsa.PublicKey
	GetPublicKeyPEM() (string, error)
}

type Auth struct {
	log          *slog.Logger
	cfg          *ssoconfig.Config
	userSaver    IUserSaver
	tokenRepo    ITokenRepo
	userProvider IUserProvider
	appProvider  IAppProvider
	keyManager   IKeyManager
}

func New(
	log *slog.Logger, cfg *ssoconfig.Config, userSaver IUserSaver,
	userProvider IUserProvider, appProvider IAppProvider,
	tokenRepo ITokenRepo, keyManager IKeyManager,
) *Auth {
	return &Auth{
		log:          log,
		cfg:          cfg,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenRepo:    tokenRepo,
		keyManager:   keyManager,
	}
}

func (a *Auth) Logout(
	ctx context.Context, refreshToken string,
) (bool, error) {
	const op = "auth.Logout"

	log := a.log.With(
		slog.String("op", op),
	)

	jwtData, err := jwt.ParseRefresh(refreshToken, a.keyManager)
	if err != nil {
		log.Error("parse token error", "err", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	hashToken := hmac.HashJWTTokenHMAC(refreshToken, a.cfg.Security.Secret)

	err = a.tokenRepo.SaveToken(ctx, hashToken, jwtData.Exp)
	if err != nil {
		if errors.Is(err, storage.ErrTokenRevoked) {
			return false, fmt.Errorf("%s: %w", op, ErrTokenRevoked)
		}

		log.Error("unknown err revoke token", "err", err)
		return false, fmt.Errorf("%s: %w", op, ErrTokenUnknown)
	}

	log.Info("user logouting")
	return true, nil
}

// Login checks if user with given credentials exists in the system
func (a *Auth) Login(
	ctx context.Context, username, password string, appID int,
) (string, string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)

	user, err := a.userProvider.User(ctx, username)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return empty, empty, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("unknown err get user", "err", err.Error())
		return empty, empty, fmt.Errorf("%s: %w", op, ErrUserUnknown)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return empty, empty, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return empty, empty, fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}

		log.Error("unknown err get app", "err", err.Error())
		return empty, empty, fmt.Errorf("%s:%w", op, ErrAppUnknown)
	}

	tokenGen := jwt.New(
		user, app, time.Duration(a.cfg.Security.AccessTokenTTL),
		time.Duration(a.cfg.Security.RefreshTokenTTL), a.keyManager,
	)

	tokenPair, err := tokenGen.GenerateTokenPair()
	if err != nil {
		log.Error("failed to generate token", "err", err)
		return empty, empty, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user loggining")
	return tokenPair.AccessToken, tokenPair.RefreshToken, nil
}

// RegisterNewUser registers new user in the system and returns userID
func (a *Auth) RegisterNewUser(
	ctx context.Context, username, password string,
) (uint, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", "err", err)
		return 0, fmt.Errorf("%s: %w", op, ErrHashPassword)
	}

	id, err := a.userSaver.SaveUser(ctx, username, string(passHash))
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			return 0, fmt.Errorf("%s: %w", op, ErrUserExist)
		}

		log.Error("unknown err save user", "err", err)
		return 0, fmt.Errorf("%s: %w", op, ErrUserUnknown)
	}

	return id, nil
}

// IsAdmin checks if user is admin
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
	)

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		log.Error("unknown get isadmin error", "err", err)
		return false, fmt.Errorf("%s: %w", op, ErrUserUnknown)
	}

	return isAdmin, nil
}

// return publickey for parse jwt tokens
func (a *Auth) PublicKey(ctx context.Context, appID int) (string, error) {
	const op = "auth.PublicKey"

	log := a.log.With(
		slog.String("op", op),
	)

	_, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return empty, fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}

		log.Error("unknown get app error", "err", err)
		return empty, fmt.Errorf("%s: %w", op, ErrAppUnknown)
	}

	return a.keyManager.GetPublicKeyPEM()
}

func (a *Auth) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	const op = "auth.Refresh"

	log := a.log.With(
		slog.String("op", op),
	)

	claims, err := jwt.ParseRefresh(refreshToken, a.keyManager)
	if err != nil {
		log.Error("parse token err", "err", err)
		return empty, empty, fmt.Errorf("%s: %w", op, err)
	}

	hashToken := hmac.HashJWTTokenHMAC(refreshToken, a.cfg.Security.Secret)

	err = a.tokenRepo.SaveToken(ctx, hashToken, claims.Exp)
	if err != nil {
		if errors.Is(err, storage.ErrTokenRevoked) {
			return empty, empty, fmt.Errorf("%s: %w", op, ErrTokenRevoked)
		}

		log.Error("unknown err revoke token", "err", err)
		return empty, empty, fmt.Errorf("%s: %w", op, ErrTokenUnknown)
	}

	user, err := a.userProvider.User(ctx, claims.Username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return empty, empty, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("unknown err get user", "err", err.Error())
		return empty, empty, fmt.Errorf("%s: %w", op, ErrUserUnknown)
	}

	app, err := a.appProvider.App(ctx, claims.AppID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return empty, empty, fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}

		log.Error("unknown err get app", "err", err.Error())
		return empty, empty, fmt.Errorf("%s:%w", op, ErrAppUnknown)
	}

	tokenGen := jwt.New(
		user, app, time.Duration(a.cfg.Security.AccessTokenTTL),
		time.Duration(a.cfg.Security.RefreshTokenTTL), a.keyManager,
	)

	tokenPair, err := tokenGen.GenerateTokenPair()
	if err != nil {
		log.Error("failed to generate token", "err", err.Error())
		return empty, empty, fmt.Errorf("%s: %w", op, err)
	}

	return tokenPair.AccessToken, tokenPair.RefreshToken, nil
}
