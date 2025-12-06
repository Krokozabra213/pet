package app

import (
	"log/slog"

	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	"github.com/Krokozabra213/sso/configs/chatconfig"
	appgrpc "github.com/Krokozabra213/sso/internal/chat/app/grpc"
	chatBusiness "github.com/Krokozabra213/sso/internal/chat/business"
	chatgrpc "github.com/Krokozabra213/sso/internal/chat/grpc"
	"github.com/Krokozabra213/sso/internal/chat/repository/broker"
	postgresrepo "github.com/Krokozabra213/sso/internal/chat/repository/storage/postgres-repo"
	custombroker "github.com/Krokozabra213/sso/pkg/custom-broker"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
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
	// TODO: Вынести в конфиг константы
	broker, err := custombroker.NewCBroker(2, 1000, 10)
	if err != nil {
		panic(err)
	}
	return broker
}

func (builder *ChatAppBuilder) PGConn() *postgrespet.PGDB {
	// TODO: раскомментировать
	// return postgrespet.NewPGDB(builder.cfg.DB.DSN)
	return nil
}

// repositories
// TODO:
func (builder *ChatAppBuilder) ClientProvider(brokerConn *custombroker.CBroker) chatBusiness.IClientRepo {
	return broker.New(brokerConn)
}

func (builder *ChatAppBuilder) MessageProvider(brokerConn *custombroker.CBroker) chatBusiness.IMessageRepo {
	return broker.New(brokerConn)
}

func (builder *ChatAppBuilder) DefaultMessageSaver(dbconn *postgrespet.PGDB) chatBusiness.IDefaultMessageSaver {
	return postgresrepo.New(dbconn)
}

// business-logic
func (builder *ChatAppBuilder) Business(
	clientProvider chatBusiness.IClientRepo,
	messageProvider chatBusiness.IMessageRepo,
	defaultMessageSaver chatBusiness.IDefaultMessageSaver,
) chatgrpc.IBusiness {
	return chatBusiness.New(
		builder.log, builder.cfg,
		clientProvider, messageProvider,
		defaultMessageSaver,
	)
}

func (builder *ChatAppBuilder) Handler(business chatgrpc.IBusiness) chat.ChatServer {
	return chatgrpc.New(builder.log, business)
}

func (builder *ChatAppBuilder) BuildGRPCApp(handler chat.ChatServer) *appgrpc.App {
	return appgrpc.New(builder.log, builder.cfg.Server.Host, builder.cfg.Server.Port, handler)
}
