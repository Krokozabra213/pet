package app

import (
	"log/slog"
	"time"

	appgrpc "github.com/Krokozabra213/sso/internal/auth/app/grpc"
	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage/postgres"
)

type App struct {
	GRPCSrv *appgrpc.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {

	postgres := postgres.New()
	business := authBusiness.New(log, postgres, postgres, postgres, tokenTTL)
	grpcApp := appgrpc.New(log, grpcPort, business)

	return &App{
		GRPCSrv: grpcApp,
	}
}
