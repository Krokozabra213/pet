package app

import (
	"log/slog"

	"github.com/Krokozabra213/sso/configs/ssoconfig"
	appgrpc "github.com/Krokozabra213/sso/internal/auth/app/grpc"
	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	keymanager "github.com/Krokozabra213/sso/internal/auth/lib/key-manager"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage/postgres"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage/redis"
	"github.com/Krokozabra213/sso/pkg/db"
)

type App struct {
	GRPCSrv *appgrpc.App
}

func New(
	log *slog.Logger,
	cfg *ssoconfig.Config,
) *App {
	// connects
	dbConn := db.NewPGDb(cfg.DB.DSN)
	redisConn := db.NewRedisDB(cfg.Redis.Addr, cfg.Redis.Pass, cfg.Redis.Cache)

	// repositories
	postgres := postgres.New(dbConn)
	redis := redis.New(redisConn)

	keyManager := keymanager.New(cfg.Security.PrivateKey)

	business := authBusiness.New(
		log, cfg, postgres, postgres, postgres, redis, keyManager,
	)
	grpcApp := appgrpc.New(log, cfg.Server.Port, business)

	return &App{
		GRPCSrv: grpcApp,
	}
}
