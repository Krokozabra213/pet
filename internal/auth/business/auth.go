package authBusiness

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Krokozabra213/sso/internal/auth/domain"
	"github.com/Krokozabra213/sso/internal/auth/lib/jwt"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
	"golang.org/x/crypto/bcrypt"
)

type IUserSaver interface {
	SaveUser(ctx context.Context, username string, passHash []byte) (uid int64, err error)
}

type IUserProvider interface {
	User(ctx context.Context, username string) (domain.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type IAppProvider interface {
	App(ctx context.Context, appID int) (domain.App, error)
}

type Auth struct {
	log          *slog.Logger
	tokenTTL     time.Duration
	userSaver    IUserSaver
	userProvider IUserProvider
	appProvider  IAppProvider
}

func New(
	log *slog.Logger, userSaver IUserSaver,
	userProvider IUserProvider, appProvider IAppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system
func (a *Auth) Login(
	ctx context.Context, username, password string, appID int,
) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)

	user, err := a.userProvider.User(ctx, username)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", "err", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	// if err != nil {
	// 	return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	// }

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return "", fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}

		log.Error("failed to get app", "err", err.Error())
		return "", fmt.Errorf("%s:%w", op, ErrUnknown)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", "err", err.Error())

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user loggining")

	return token, nil
}

// RegisterNewUser registers new user in the system and returns userID
func (a *Auth) RegisterNewUser(
	ctx context.Context, username, password string,
) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", "err", err)

		return 0, fmt.Errorf("%s: %w", op, ErrUnknown)
	}

	id, err := a.userSaver.SaveUser(ctx, username, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			return 0, fmt.Errorf("%s: %w", op, ErrUserExist)
		}

		log.Error("failed to save user", "err", err)
		return 0, fmt.Errorf("%s: %w", op, ErrUnknown)
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
		if errors.Is(err, storage.ErrUserNotFound) {
			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		if errors.Is(err, storage.ErrPermission) {
			return false, fmt.Errorf("%s: %w", op, ErrPermission)
		}
		log.Error("unknown isadmin error", "err", err)
		return false, fmt.Errorf("%s: %w", op, ErrUnknown)
	}

	return isAdmin, nil
}

// IsAdmin giving sault for app
func (a *Auth) AppSault(ctx context.Context, appID int) (string, error) {
	const op = "auth.AppSault"

	log := a.log.With(
		slog.String("op", op),
	)

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return "", fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}

		log.Error("unknown appsault error", "err", err)
		return "", fmt.Errorf("%s: %w", op, ErrUnknown)
	}

	return app.Sault, nil
}
