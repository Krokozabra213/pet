package authBusiness

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"time"

	"github.com/Krokozabra213/sso/configs/ssoconfig"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	authgrpc "github.com/Krokozabra213/sso/internal/auth/grpc"
	"github.com/Krokozabra213/sso/internal/auth/lib/hmac"
	"github.com/Krokozabra213/sso/internal/auth/lib/jwt"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage/postgres"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage/redis"
	"golang.org/x/crypto/bcrypt"
)

type ITokenProvider interface {
	SaveToken(ctx context.Context, token string, expiresAt time.Time) error
	CheckToken(ctx context.Context, token string) (bool, error)
}

type IUserProvider interface {
	SaveUser(ctx context.Context, username string, pass string) (uid uint, err error)
	User(ctx context.Context, username string) (*domain.User, error)
	UserByID(ctx context.Context, userID int) (*domain.User, error)
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
	)
	log.Info("user register")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", "err", err)
		return 0, fmt.Errorf("%s: %w", op, ErrHashPassword)
	}

	errorMap := map[error]error{
		postgres.ErrPGDuplicate: ErrUserExist,
		postgres.ErrContext:     ErrTimeout,
	}

	id, err := a.userProvider.SaveUser(ctx, username, string(passHash))
	if err != nil {
		return 0, ErrorGateway(op, log, err, errorMap, ErrUserUnknown)
	}

	return id, nil
}

//-------------------------LOGIN-LOGIC-------------------------------------------------//

func (a *Auth) Login(
	ctx context.Context, username, password string, appID int,
) (string, string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("user login")

	errorMap := map[error]error{
		postgres.ErrPGNotFound: ErrInvalidAppId,
		postgres.ErrContext:    ErrTimeout,
	}

	app, err := a.appProvider.AppByID(ctx, appID)
	if err != nil {
		return "", "", ErrorGateway(op, log, err, errorMap, ErrAppUnknown)
	}

	errorMap = map[error]error{
		postgres.ErrPGNotFound: ErrInvalidCredentials,
		postgres.ErrContext:    ErrTimeout,
	}

	user, err := a.userProvider.User(ctx, username)
	if err != nil {
		return "", "", ErrorGateway(op, log, err, errorMap, ErrUserUnknown)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	tokenGen := jwt.New(
		user, app, time.Duration(a.cfg.Security.AccessTokenTTL),
		time.Duration(a.cfg.Security.RefreshTokenTTL), a.keyManager,
	)

	tokenPair, err := tokenGen.GenerateTokenPair()
	if err != nil {
		log.Error("failed to generate token", "err", err)
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user loggining")
	return tokenPair.AccessToken, tokenPair.RefreshToken, nil
}

//-------------------------LOGOUT-LOGIC-------------------------------------------------//

func (a *Auth) Logout(
	ctx context.Context, refreshToken string,
) (bool, error) {
	const op = "auth.Logout"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("user logout")

	claims, err := jwt.ParseRefresh(refreshToken, a.keyManager)
	if err != nil {
		log.Error("failed parse token", "err", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	hashToken := hmac.HashJWTTokenHMAC(refreshToken, a.cfg.Security.Secret)

	errorMap := map[error]error{
		redis.ErrTokenExpired: ErrTokenExpired,
		postgres.ErrContext:   ErrTimeout,
	}

	err = a.tokenRepo.SaveToken(ctx, hashToken, claims.Exp)
	if err != nil {
		return false, ErrorGateway(op, log, err, errorMap, ErrTokenUnknown)
	}
	log.Info("user logouting")

	return true, nil
}

//-------------------------ISADMIN-LOGIC-------------------------------------------------//

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
	)

	errorMap := map[error]error{
		postgres.ErrContext: ErrTimeout,
	}

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, ErrorGateway(op, log, err, errorMap, ErrUserUnknown)
	}

	if !isAdmin {
		return false, fmt.Errorf("%s: %w", op, ErrPermission)
	}

	return isAdmin, nil
}

//-------------------------PUBLICKEY-LOGIC-------------------------------------------------//

func (a *Auth) PublicKey(ctx context.Context, appID int) (string, error) {
	const op = "auth.PublicKey"

	log := a.log.With(
		slog.String("op", op),
	)

	errorMap := map[error]error{
		postgres.ErrPGNotFound: ErrInvalidAppId,
		postgres.ErrContext:    ErrTimeout,
	}

	_, err := a.appProvider.AppByID(ctx, appID)
	if err != nil {
		return "", ErrorGateway(op, log, err, errorMap, ErrAppUnknown)
	}

	return a.keyManager.GetPublicKeyPEM()
}

//-------------------------REFRESH-LOGIC-------------------------------------------------//

func (a *Auth) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	const op = "auth.Refresh"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("user try to refresh token")

	// хешируем токен для безопасного хранения в базе
	hashToken := hmac.HashJWTTokenHMAC(refreshToken, a.cfg.Security.Secret)

	exist, err := a.tokenRepo.CheckToken(ctx, hashToken)
	if err != nil {
		log.Error("check token err", "err", err)
		return "", "", fmt.Errorf("%s: %w", op, ErrTokenUnknown)
	}
	if exist {
		return "", "", ErrTokenRevoked
	}

	// достаём claims из refresh токена
	claims, err := jwt.ParseRefresh(refreshToken, a.keyManager)
	if err != nil {
		log.Error("parse token err", "err", err)
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	errorMap := map[error]error{
		postgres.ErrPGNotFound: ErrInvalidCredentials,
		postgres.ErrContext:    ErrTimeout,
	}

	// проверяем наличие пользователя с таким id
	user, err := a.userProvider.UserByID(ctx, claims.UserID)
	if err != nil {
		return "", "", ErrorGateway(op, log, err, errorMap, ErrUserUnknown)
	}

	errorMap = map[error]error{
		postgres.ErrPGNotFound: ErrInvalidAppId,
		postgres.ErrContext:    ErrTimeout,
	}

	// проверяем наличие приложения с таким id
	app, err := a.appProvider.AppByID(ctx, claims.AppID)
	if err != nil {
		return "", "", ErrorGateway(op, log, err, errorMap, ErrAppUnknown)
	}

	// генерируем новую пару токенов
	tokenGen := jwt.New(
		user, app, time.Duration(a.cfg.Security.AccessTokenTTL),
		time.Duration(a.cfg.Security.RefreshTokenTTL), a.keyManager,
	)

	tokenPair, err := tokenGen.GenerateTokenPair()
	if err != nil {
		log.Error("failed to generate token", "err", err.Error())
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	errorMap = map[error]error{
		redis.ErrTokenExpired: ErrTokenExpired,
		postgres.ErrContext:   ErrTimeout,
	}

	err = a.tokenRepo.SaveToken(ctx, hashToken, claims.Exp)
	if err != nil {
		return "", "", ErrorGateway(op, log, err, errorMap, ErrTokenUnknown)
	}
	log.Info("user refreshed token", "tokens", tokenPair.AccessToken+", "+tokenPair.RefreshToken)

	return tokenPair.AccessToken, tokenPair.RefreshToken, nil
}
