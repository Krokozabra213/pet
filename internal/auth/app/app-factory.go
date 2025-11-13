package app

import (
	"github.com/Krokozabra213/protos/gen/go/proto/sso"
	appgrpc "github.com/Krokozabra213/sso/internal/auth/app/grpc"
	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	authgrpc "github.com/Krokozabra213/sso/internal/auth/grpc"
	"github.com/Krokozabra213/sso/pkg/db"
)

type IAppBuilder interface {
	// Connects
	DBConn() *db.Db
	NoSQLDBConn() *db.RDB

	// Repositories
	UserProvider(connect *db.Db) authBusiness.IUserProvider
	AppProvider(connect *db.Db) authBusiness.IAppProvider
	TokenProvider(connect *db.RDB) authBusiness.ITokenProvider

	// Libraries
	KeyManager() authBusiness.IKeyManager

	// Business-Logic
	Business(
		userProvider authBusiness.IUserProvider,
		appProvider authBusiness.IAppProvider,
		tokenRepo authBusiness.ITokenProvider,
		keyManager authBusiness.IKeyManager,
	) authgrpc.IBusiness

	// Handler
	Handler(business authgrpc.IBusiness) sso.AuthServer

	// Application
	BuildGRPCApp(handler sso.AuthServer) *appgrpc.App
}

type AppFactory struct {
	builder IAppBuilder
}

func NewAppFactory(builder IAppBuilder) *AppFactory {
	return &AppFactory{
		builder: builder,
	}
}

func (fack *AppFactory) Create() *appgrpc.App {
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
