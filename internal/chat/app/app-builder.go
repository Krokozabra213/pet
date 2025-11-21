package app

import (
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/configs/chatconfig"
	appgrpc "github.com/Krokozabra213/sso/internal/chat/app/grpc"
	chatBusiness "github.com/Krokozabra213/sso/internal/chat/business"
	chatgrpc "github.com/Krokozabra213/sso/internal/chat/grpc"
	"github.com/Krokozabra213/sso/internal/chat/repository/broker"
)

type ChatAppBuilder struct {
	cfg *chatconfig.Config
	log *slog.Logger
}

func NewAppBuilder(cfg *chatconfig.Config, log *slog.Logger) *ChatAppBuilder {
	return &ChatAppBuilder{
		cfg: cfg,
		log: log,
	}
}

// repositories
// TODO:
func (builder *ChatAppBuilder) ClientProvider() chatBusiness.IClientRepo {
	return broker.New()
}

func (builder *ChatAppBuilder) MessageProvider() chatBusiness.IMessageRepo {
	return broker.New()
}

// business-logic
func (builder *ChatAppBuilder) Business(
	clientProvider chatBusiness.IClientRepo,
	messageProvider chatBusiness.IMessageRepo,
) chatgrpc.IBusiness {
	return chatBusiness.New(
		builder.log, builder.cfg,
		clientProvider, messageProvider,
	)
}

func (builder *ChatAppBuilder) Handler(business chatgrpc.IBusiness) chat.ChatServer {
	return chatgrpc.New(business)
}

func (builder *ChatAppBuilder) BuildGRPCApp(handler chat.ChatServer) *appgrpc.App {
	return appgrpc.New(builder.log, builder.cfg.Server.Host, builder.cfg.Server.Port, handler)
}
