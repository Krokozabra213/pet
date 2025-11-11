package chatapp

import (
	"log/slog"

	"github.com/Krokozabra213/sso/configs/chatconfig"
	chatappgrpc "github.com/Krokozabra213/sso/internal/chat/app/grpc"
	chatBusiness "github.com/Krokozabra213/sso/internal/chat/business"
	"github.com/Krokozabra213/sso/internal/chat/repository/broker"
)

type App struct {
	GRPCSrv *chatappgrpc.App
}

func New(
	log *slog.Logger,
	cfg *chatconfig.Config,
) *App {
	// connects
	//todo: add broker conn
	// repositories
	msgBroker := broker.New()
	//business
	business := chatBusiness.New(
		log, cfg, msgBroker, msgBroker,
	)
	grpcApp := chatappgrpc.New(log, cfg.Server.Port, business)

	return &App{
		GRPCSrv: grpcApp,
	}
}
