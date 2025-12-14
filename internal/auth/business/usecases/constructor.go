package authusecases

import (
	"context"
	"crypto/rsa"
	"log/slog"
	"time"

	"github.com/Krokozabra213/sso/configs/ssoconfig"
	"github.com/Krokozabra213/sso/internal/auth/domain"
	authgrpc "github.com/Krokozabra213/sso/internal/auth/grpc/auth-grpc"
)

type ITokenProvider interface {
	SaveToken(ctx context.Context, token string, expiresAt time.Time) error
	CheckToken(ctx context.Context, token string) (bool, error)
}

type IUserProvider interface {
	SaveUser(ctx context.Context, user *domain.User) (uid uint64, err error)
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
