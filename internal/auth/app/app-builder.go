package app

import (
	"crypto/rsa"

	"github.com/Krokozabra213/protos/gen/go/sso"
	appgrpc "github.com/Krokozabra213/sso/internal/auth/app/grpc"
	authusecases "github.com/Krokozabra213/sso/internal/auth/business/usecases"
	authgrpc "github.com/Krokozabra213/sso/internal/auth/grpc/auth-grpc"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage/postgres"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage/redis"
	ssonewconfig "github.com/Krokozabra213/sso/newconfigs/sso"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
	redispet "github.com/Krokozabra213/sso/pkg/db/redis-pet"
	hmacv1 "github.com/Krokozabra213/sso/pkg/hmac/v1"
	jwtv1 "github.com/Krokozabra213/sso/pkg/jwt-manager/v1"
	keymanagerv1 "github.com/Krokozabra213/sso/pkg/key-manager/v1"
)

type IKeyManager interface {
	GetPrivateKey() *rsa.PrivateKey
	GetPublicKey() *rsa.PublicKey
	GetPublicKeyPEM() string
}

type AuthAppBuilder struct {
	cfg *ssonewconfig.Config
}

func NewAppBuilder(cfg *ssonewconfig.Config) *AuthAppBuilder {
	return &AuthAppBuilder{
		cfg: cfg,
	}
}

// connects
func (builder *AuthAppBuilder) DBConn() *postgrespet.PGDB {
	return postgrespet.NewPGDB(builder.cfg.PG.DSN)
}

func (builder *AuthAppBuilder) NoSQLDBConn() *redispet.RDB {
	return redispet.NewRedisDB(
		builder.cfg.Redis.Addr, builder.cfg.Redis.Pass, builder.cfg.Redis.Cache,
	)
}

// repositories
func (builder *AuthAppBuilder) UserProvider(connect *postgrespet.PGDB) authusecases.IUserProvider {
	return postgres.New(connect)
}

func (builder *AuthAppBuilder) AppProvider(connect *postgrespet.PGDB) authusecases.IAppProvider {
	return postgres.New(connect)
}

func (builder *AuthAppBuilder) TokenProvider(connect *redispet.RDB) authusecases.ITokenProvider {
	return redis.New(connect)
}

// libraries
func (builder *AuthAppBuilder) KeyManager() IKeyManager {
	manager, err := keymanagerv1.New(builder.cfg.Auth.JWT.PrivateKey)
	if err != nil {
		panic(err)
	}
	return manager
}

func (builder *AuthAppBuilder) Hasher() authusecases.IHasher {
	return hmacv1.New(builder.cfg.Auth.AppSecretKey)
}

func (builder *AuthAppBuilder) JWTManager(public *rsa.PublicKey, private *rsa.PrivateKey) authusecases.IJWTManager {
	manager, err := jwtv1.New(public, private)
	if err != nil {
		panic(err)
	}
	return manager
}

// business-logic
func (builder *AuthAppBuilder) Business(
	userProvider authusecases.IUserProvider, appProvider authusecases.IAppProvider,
	tokenRepo authusecases.ITokenProvider, jwtManager authusecases.IJWTManager, hasher authusecases.IHasher,
	publicKeyPEM string,
) authgrpc.IBusiness {
	return authusecases.New(
		builder.cfg, userProvider, appProvider, tokenRepo, jwtManager, hasher, publicKeyPEM,
	)
}

func (builder *AuthAppBuilder) Handler(business authgrpc.IBusiness) sso.AuthServer {
	return authgrpc.New(business)
}

func (builder *AuthAppBuilder) BuildGRPCApp(handler sso.AuthServer) *appgrpc.GRPCApp {
	return appgrpc.New(builder.cfg.GRPC.Host, builder.cfg.GRPC.Port, handler)
}
