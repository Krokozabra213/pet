package authusecases

import (
	"context"
	"time"

	"github.com/Krokozabra213/sso/internal/auth/domain"
	authgrpc "github.com/Krokozabra213/sso/internal/auth/grpc/auth-grpc"
	ssonewconfig "github.com/Krokozabra213/sso/newconfigs/sso"
	jwtv1 "github.com/Krokozabra213/sso/pkg/jwt-manager/v1"
)

//go:generate mockgen  -source=constructor.go -destination=mocks/mocks.go

type IJWTManager interface {
	GenerateTokens(data *jwtv1.Data) (string, string, error)
	ParseAccess(token string) (*jwtv1.AccessData, error)
	ParseRefresh(token string) (*jwtv1.RefreshData, error)
}

type IHasher interface {
	HashJWTTokenHMAC(token string) string
}

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

type Auth struct {
	cfg          *ssonewconfig.Config
	tokenRepo    ITokenProvider
	userProvider IUserProvider
	appProvider  IAppProvider
	jwtManager   IJWTManager
	hasher       IHasher
	publicKeyPEM string
}

func New(
	cfg *ssonewconfig.Config, userProvider IUserProvider, appProvider IAppProvider, tokenRepo ITokenProvider,
	jwtManager IJWTManager, hasher IHasher, publicKeyPEM string,
) authgrpc.IBusiness {
	return &Auth{
		cfg:          cfg,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenRepo:    tokenRepo,
		jwtManager:   jwtManager,
		hasher:       hasher,
		publicKeyPEM: publicKeyPEM,
	}
}
