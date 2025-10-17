package app

import (
	"log/slog"
	"time"

	authgrpc "github.com/Krokozabra213/sso/internal/auth/app/grpc"
)

type App struct {
	GRPCSrv *authgrpc.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: инициализировать хранилище (storage)

	// TODO: init business

	grpcApp := authgrpc.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
