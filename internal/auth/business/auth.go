package authBusiness

import (
	"context"
	"crypto/rsa"
	"errors"
	"log/slog"
	"time"

	"github.com/Krokozabra213/sso/configs/ssoconfig"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	authgrpc "github.com/Krokozabra213/sso/internal/auth/grpc"
	"github.com/Krokozabra213/sso/internal/auth/lib/hmac"
	"github.com/Krokozabra213/sso/internal/auth/lib/jwt"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
	"golang.org/x/crypto/bcrypt"
)

type ITokenProvider interface {
	SaveToken(ctx context.Context, token string, expiresAt time.Time) error
	CheckToken(ctx context.Context, token string) (bool, error)
}

type IUserProvider interface {
	SaveUser(ctx context.Context, user *domain.User) (uid uint, err error)
	User(ctx context.Context, username string) (*domain.User, error)
	UserByID(ctx context.Context, userID int64) (*domain.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type IAppProvider interface {
	AppByID(ctx context.Context, appID int) (*domain.App, error)
}

type IKeyManager interface {
	GetPrivateKey() *rsa.PrivateKey
	GetPublicKey() *rsa.PublicKey
	GetPublicKeyPEM() (string, error)
}

type Auth struct {
	log          *slog.Logger
	cfg          *ssoconfig.Config
	tokenRepo    ITokenProvider
	userProvider IUserProvider
	appProvider  IAppProvider
	keyManager   IKeyManager
}

func New(
	log *slog.Logger, cfg *ssoconfig.Config,
	userProvider IUserProvider, appProvider IAppProvider,
	tokenRepo ITokenProvider, keyManager IKeyManager,
) authgrpc.IBusiness {
	return &Auth{
		log:          log,
		cfg:          cfg,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenRepo:    tokenRepo,
		keyManager:   keyManager,
	}
}

//-------------------------REGISTER-LOGIC-------------------------------------------------//

func (a *Auth) RegisterNewUser(
	ctx context.Context, username, password string,
) (uint, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)
	log.Info("starting user registration process")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed password hashing", slog.String("error", err.Error()))
		return 0, BusinessError(domain.UserEntity, ErrHashPassword)
	}

	user := domain.NewUser(username, string(passHash))
	userID, err := a.userProvider.SaveUser(ctx, user)
	if err != nil {
		log.Error("failed save new user", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return 0, BusinessError(domain.UserEntity, ErrTimeout)
		}
		if errors.Is(err, storage.ErrDuplicate) {
			return 0, BusinessError(domain.UserEntity, ErrExists)
		}
		return 0, BusinessError(domain.UserEntity, ErrInternal)
	}

	log.Info("user successfully registered",
		slog.Uint64("user_id", uint64(userID)))

	return userID, nil
}

//-------------------------LOGIN-LOGIC-------------------------------------------------//

func (a *Auth) Login(
	ctx context.Context, username, password string, appID int,
) (string, string, error) {
	const op = "auth.Login"

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
			return "", "", BusinessError(domain.AppEntity, ErrTimeout)
		}
		if errors.Is(err, storage.ErrNotFound) {
			return "", "", BusinessError(domain.AppEntity, ErrNotFound)
		}
		return "", "", BusinessError(domain.AppEntity, ErrInternal)
	}

	user, err := a.userProvider.User(ctx, username)
	if err != nil {
		log.Error("failed get user", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return "", "", BusinessError(domain.UserEntity, ErrTimeout)
		}
		if errors.Is(err, storage.ErrNotFound) {
			return "", "", BusinessError(domain.UserEntity, ErrNotFound)
		}
		return "", "", BusinessError(domain.UserEntity, ErrInternal)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Error("failed compare passwords", slog.String("error", err.Error()))
		return "", "", BusinessError(domain.UserEntity, ErrInvalidCredentials)
	}

	tokenGen := jwt.New(
		user, app, time.Duration(a.cfg.Security.AccessTokenTTL),
		time.Duration(a.cfg.Security.RefreshTokenTTL), a.keyManager,
	)

	tokenPair, err := tokenGen.GenerateTokenPair()
	if err != nil {
		log.Error("failed generate token", slog.String("error", err.Error()))
		return "", "", BusinessError(domain.TokenEntity, ErrTokenGenerate)
	}

	log.Info("user successfully logining")
	return tokenPair.AccessToken, tokenPair.RefreshToken, nil
}

//-------------------------LOGOUT-LOGIC-------------------------------------------------//

func (a *Auth) Logout(
	ctx context.Context, refreshToken string,
) (bool, error) {
	const op = "auth.Logout"

	log := a.log.With(
		slog.String("op", op),
		slog.String("token", refreshToken),
	)
	log.Info("starting user logouting process")

	claims, err := jwt.ParseRefresh(refreshToken, a.keyManager)
	if err != nil {
		log.Error("failed parsing token", slog.String("error", err.Error()))
		return false, BusinessError(domain.TokenEntity, ErrParse)
	}

	hashToken := hmac.HashJWTTokenHMAC(refreshToken, a.cfg.Security.Secret)

	err = a.tokenRepo.SaveToken(ctx, hashToken, claims.Exp)
	if err != nil {
		log.Error("failed revoking token", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return false, BusinessError(domain.TokenEntity, ErrTimeout)
		}
		if errors.Is(err, storage.ErrTokenExpired) {
			return false, BusinessError(domain.TokenEntity, ErrTokenExpired)
		}
		return false, BusinessError(domain.TokenEntity, ErrInternal)
	}
	log.Info("user successfully logouting")

	return true, nil
}

//-------------------------ISADMIN-LOGIC-------------------------------------------------//

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)
	log.Info("starting isadmin process")

	_, err := a.userProvider.UserByID(ctx, userID)
	if err != nil {
		log.Error("failed to get user by id", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return false, BusinessError(domain.UserEntity, ErrTimeout)
		}
		if errors.Is(err, storage.ErrNotFound) {
			return false, BusinessError(domain.UserEntity, ErrNotFound)
		}
		return false, BusinessError(domain.UserEntity, ErrInternal)
	}

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		log.Error("failed checking isadmin user", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return false, BusinessError(domain.AdminEntity, ErrTimeout)
		}
		return false, BusinessError(domain.AdminEntity, ErrInternal)
	}

	if !isAdmin {
		log.Info("err permission")
		return false, BusinessError(domain.AdminEntity, ErrPermission)
	}

	log.Info("ending isadmin process")

	return isAdmin, nil
}

//-------------------------PUBLICKEY-LOGIC-------------------------------------------------//

func (a *Auth) PublicKey(ctx context.Context, appID int) (string, error) {
	const op = "auth.PublicKey"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("app_id", appID),
	)
	log.Info("starting publickey process")

	_, err := a.appProvider.AppByID(ctx, appID)
	if err != nil {
		log.Error("failed get app by id", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return "", BusinessError(domain.AppEntity, ErrTimeout)
		}
		if errors.Is(err, storage.ErrNotFound) {
			return "", BusinessError(domain.AppEntity, ErrNotFound)
		}
		return "", BusinessError(domain.AppEntity, ErrInternal)
	}

	log.Info("ending publickey process")

	return a.keyManager.GetPublicKeyPEM()
}

//-------------------------REFRESH-LOGIC-------------------------------------------------//

func (a *Auth) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	const op = "auth.Refresh"

	log := a.log.With(
		slog.String("op", op),
		slog.String("token", refreshToken),
	)
	log.Info("starting refresh token process")

	// хешируем токен для безопасного хранения в базе
	hashToken := hmac.HashJWTTokenHMAC(refreshToken, a.cfg.Security.Secret)

	exist, err := a.tokenRepo.CheckToken(ctx, hashToken)
	if err != nil {
		log.Error("failed to check token", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return "", "", BusinessError(domain.TokenEntity, ErrTimeout)
		}
		return "", "", BusinessError(domain.TokenEntity, ErrInternal)
	}
	if exist {
		return "", "", BusinessError(domain.TokenEntity, ErrTokenRevoked)
	}

	// достаём claims из refresh токена
	claims, err := jwt.ParseRefresh(refreshToken, a.keyManager)
	if err != nil {
		log.Error("failed to parse token", slog.String("error", err.Error()))
		return "", "", BusinessError(domain.TokenEntity, ErrParse)
	}

	// проверяем наличие пользователя с таким id
	user, err := a.userProvider.UserByID(ctx, int64(claims.UserID))
	if err != nil {
		log.Error("failed to get user by id", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return "", "", BusinessError(domain.UserEntity, ErrTimeout)
		}
		if errors.Is(err, storage.ErrNotFound) {
			return "", "", BusinessError(domain.UserEntity, ErrNotFound)
		}
		return "", "", BusinessError(domain.UserEntity, ErrInternal)
	}

	// проверяем наличие приложения с таким id
	app, err := a.appProvider.AppByID(ctx, claims.AppID)
	if err != nil {
		log.Error("failed to get app by id", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return "", "", BusinessError(domain.AppEntity, ErrTimeout)
		}
		if errors.Is(err, storage.ErrNotFound) {
			return "", "", BusinessError(domain.AppEntity, ErrNotFound)
		}
		return "", "", BusinessError(domain.AppEntity, ErrInternal)
	}

	// генерируем новую пару токенов
	tokenGen := jwt.New(
		user, app, time.Duration(a.cfg.Security.AccessTokenTTL),
		time.Duration(a.cfg.Security.RefreshTokenTTL), a.keyManager,
	)

	tokenPair, err := tokenGen.GenerateTokenPair()
	if err != nil {
		log.Error("failed to generate token", "err", err.Error())
		return "", "", BusinessError(domain.TokenEntity, ErrTokenGenerate)
	}

	err = a.tokenRepo.SaveToken(ctx, hashToken, claims.Exp)
	if err != nil {
		log.Error("failed to save revoking token", "err", err.Error())
		if errors.Is(err, storage.ErrCtxCancelled) || errors.Is(err, storage.ErrCtxDeadline) {
			return "", "", BusinessError(domain.TokenEntity, ErrTimeout)
		}
		if errors.Is(err, storage.ErrTokenExpired) {
			return "", "", BusinessError(domain.TokenEntity, ErrTokenExpired)
		}
		return "", "", BusinessError(domain.TokenEntity, ErrInternal)
	}
	log.Info("user refreshed token", "tokens", tokenPair.AccessToken+", "+tokenPair.RefreshToken)

	return tokenPair.AccessToken, tokenPair.RefreshToken, nil
}
