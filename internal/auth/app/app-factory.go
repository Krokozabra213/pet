package app

import (
	"github.com/Krokozabra213/protos/gen/go/sso"
	appgrpc "github.com/Krokozabra213/sso/internal/auth/app/grpc"
	authusecases "github.com/Krokozabra213/sso/internal/auth/business/usecases"
	authgrpc "github.com/Krokozabra213/sso/internal/auth/grpc/auth-grpc"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
	redispet "github.com/Krokozabra213/sso/pkg/db/redis-pet"
)

type IAppBuilder interface {
	// Connects
	DBConn() *postgrespet.PGDB
	NoSQLDBConn() *redispet.RDB

	// Repositories
	UserProvider(connect *postgrespet.PGDB) authusecases.IUserProvider
	AppProvider(connect *postgrespet.PGDB) authusecases.IAppProvider
	TokenProvider(connect *redispet.RDB) authusecases.ITokenProvider

	// Libraries
	KeyManager() authusecases.IKeyManager

	// Business-Logic
	Business(
		userProvider authusecases.IUserProvider,
		appProvider authusecases.IAppProvider,
		tokenRepo authusecases.ITokenProvider,
		keyManager authusecases.IKeyManager,
	) authgrpc.IBusiness

	// Handler
	Handler(business authgrpc.IBusiness) sso.AuthServer

	// Application
	BuildGRPCApp(handler sso.AuthServer) *appgrpc.GRPCApp
}

type AppFactory struct {
	builder IAppBuilder
}

func NewAppFactory(builder IAppBuilder) *AppFactory {
	return &AppFactory{
		builder: builder,
	}
}

func (fack *AppFactory) Create() *appgrpc.GRPCApp {
	// Connects
	dbConn := fack.builder.DBConn()
	redisConn := fack.builder.NoSQLDBConn()

	// Repositories
	userRepo := fack.builder.UserProvider(dbConn)
	appRepo := fack.builder.AppProvider(dbConn)
	tokenRepo := fack.builder.TokenProvider(redisConn)

	//libraries
	keyManager := fack.builder.KeyManager()

	// Business-Logic
	business := fack.builder.Business(
		userRepo, appRepo, tokenRepo, keyManager,
	)

	// Handler
	handler := fack.builder.Handler(business)

	// Application
	application := fack.builder.BuildGRPCApp(handler)

	return application
}
