package app

import (
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/sso"
	"github.com/Krokozabra213/sso/configs/ssoconfig"
	appgrpc "github.com/Krokozabra213/sso/internal/auth/app/grpc"
	authusecases "github.com/Krokozabra213/sso/internal/auth/business/usecases"
	authgrpc "github.com/Krokozabra213/sso/internal/auth/grpc"
	keymanager "github.com/Krokozabra213/sso/internal/auth/lib/key-manager"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage/postgres"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage/redis"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
	redispet "github.com/Krokozabra213/sso/pkg/db/redis-pet"
)

type AuthAppBuilder struct {
	cfg *ssoconfig.Config
	log *slog.Logger
}

func NewAppBuilder(cfg *ssoconfig.Config, log *slog.Logger) *AuthAppBuilder {
	return &AuthAppBuilder{
		cfg: cfg,
		log: log,
	}
}

// connects
func (builder *AuthAppBuilder) DBConn() *postgrespet.PGDB {
	return postgrespet.NewPGDB(builder.cfg.DB.DSN)
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
func (builder *AuthAppBuilder) KeyManager() authusecases.IKeyManager {
	return keymanager.New(builder.cfg.Security.PrivateKey)
}

// business-logic
func (builder *AuthAppBuilder) Business(
	userProvider authusecases.IUserProvider, appProvider authusecases.IAppProvider,
	tokenRepo authusecases.ITokenProvider, keyManager authusecases.IKeyManager,
) authgrpc.IBusiness {
	return authusecases.New(
		builder.log, builder.cfg,
		userProvider, appProvider, // sqldb
		tokenRepo, //nosqldb
		keyManager,
	)
}

func (builder *AuthAppBuilder) Handler(business authgrpc.IBusiness) sso.AuthServer {
	return authgrpc.New(business)
}

func (builder *AuthAppBuilder) BuildGRPCApp(handler sso.AuthServer) *appgrpc.GRPCApp {
	return appgrpc.New(builder.log, builder.cfg.Server.Host, builder.cfg.Server.Port, handler)
}
