package app

import (
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/proto/sso"
	"github.com/Krokozabra213/sso/configs/ssoconfig"
	appgrpc "github.com/Krokozabra213/sso/internal/auth/app/grpc"
	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	authgrpc "github.com/Krokozabra213/sso/internal/auth/grpc"
	keymanager "github.com/Krokozabra213/sso/internal/auth/lib/key-manager"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage/postgres"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage/redis"
	"github.com/Krokozabra213/sso/pkg/db"
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
func (builder *AuthAppBuilder) DBConn() *db.Db {
	return db.NewPGDb(builder.cfg.DB.DSN)
}

func (builder *AuthAppBuilder) NoSQLDBConn() *db.RDB {
	return db.NewRedisDB(
		builder.cfg.Redis.Addr, builder.cfg.Redis.Pass, builder.cfg.Redis.Cache,
	)
}

// repositories
func (builder *AuthAppBuilder) UserProvider(connect *db.Db) authBusiness.IUserProvider {
	return postgres.New(connect)
}

func (builder *AuthAppBuilder) AppProvider(connect *db.Db) authBusiness.IAppProvider {
	return postgres.New(connect)
}

func (builder *AuthAppBuilder) TokenProvider(connect *db.RDB) authBusiness.ITokenProvider {
	return redis.New(connect)
}

// libraries
func (builder *AuthAppBuilder) KeyManager() authBusiness.IKeyManager {
	return keymanager.New(builder.cfg.Security.PrivateKey)
}

// business-logic
func (builder *AuthAppBuilder) Business(
	userProvider authBusiness.IUserProvider, appProvider authBusiness.IAppProvider,
	tokenRepo authBusiness.ITokenProvider, keyManager authBusiness.IKeyManager,
) authgrpc.IBusiness {
	return authBusiness.New(
		builder.log, builder.cfg,
		userProvider, appProvider, // sqldb
		tokenRepo, //nosqldb
		keyManager,
	)
}

func (builder *AuthAppBuilder) Handler(business authgrpc.IBusiness) sso.AuthServer {
	return authgrpc.New(business)
}

func (builder *AuthAppBuilder) BuildGRPCApp(handler sso.AuthServer) *appgrpc.App {
	return appgrpc.New(builder.log, builder.cfg.Server.Host, builder.cfg.Server.Host, handler)
}
