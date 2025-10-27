package app

import (
	"log/slog"

	"github.com/Krokozabra213/sso/configs/ssoconfig"
	appgrpc "github.com/Krokozabra213/sso/internal/auth/app/grpc"
	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	keymanager "github.com/Krokozabra213/sso/internal/auth/lib/key-manager"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage/postgres"
	"github.com/Krokozabra213/sso/pkg/db"
)

type App struct {
	GRPCSrv *appgrpc.App
}

func New(
	log *slog.Logger,
	cfg *ssoconfig.Config,
) *App {
	dbConn := db.NewPGDb(cfg.DB.DSN)
	postgres := postgres.New(dbConn)
	keyManager := keymanager.New(cfg.Security.PrivateKey)
	business := authBusiness.New(
		log, cfg, postgres, postgres, postgres, postgres, keyManager,
	)
	grpcApp := appgrpc.New(log, cfg.Server.Port, business)

	return &App{
		GRPCSrv: grpcApp,
	}
}
