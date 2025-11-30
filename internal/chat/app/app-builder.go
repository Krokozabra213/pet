package app

import (
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/configs/chatconfig"
	appgrpc "github.com/Krokozabra213/sso/internal/chat/app/grpc"
	chatBusiness "github.com/Krokozabra213/sso/internal/chat/business"
	chatgrpc "github.com/Krokozabra213/sso/internal/chat/grpc"
	"github.com/Krokozabra213/sso/internal/chat/repository/broker"
	custombroker "github.com/Krokozabra213/sso/pkg/custom-broker"
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

// connects
// TODO:
func (builder *ChatAppBuilder) BrokerConn() *custombroker.CBroker {
	broker, err := custombroker.NewCBroker(4, 1000, 10)
	if err != nil {
		panic(err)
	}
	return broker
}

// repositories
// TODO:
func (builder *ChatAppBuilder) ClientProvider(brokerConn *custombroker.CBroker) chatBusiness.IClientRepo {
	return broker.New(brokerConn)
}

func (builder *ChatAppBuilder) MessageProvider(brokerConn *custombroker.CBroker) chatBusiness.IMessageRepo {
	return broker.New(brokerConn)
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
	return chatgrpc.New(builder.log, business)
}

func (builder *ChatAppBuilder) BuildGRPCApp(handler chat.ChatServer) *appgrpc.App {
	return appgrpc.New(builder.log, builder.cfg.Server.Host, builder.cfg.Server.Port, handler)
}
