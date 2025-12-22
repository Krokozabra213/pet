package app

import (
	"github.com/Krokozabra213/protos/gen/go/chat"
	appgrpc "github.com/Krokozabra213/sso/internal/chat/app/grpc"
	chatusecases "github.com/Krokozabra213/sso/internal/chat/business/usecases"
	chatgrpc "github.com/Krokozabra213/sso/internal/chat/grpc"
	chatinterfaces "github.com/Krokozabra213/sso/internal/chat/grpc/interfaces"
	"github.com/Krokozabra213/sso/internal/chat/repository/broker"
	postgresrepo "github.com/Krokozabra213/sso/internal/chat/repository/storage/postgres-repo"
	chatnewconfig "github.com/Krokozabra213/sso/newconfigs/chat"
	custombroker "github.com/Krokozabra213/sso/pkg/custom-broker"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
)

type ChatAppBuilder struct {
	cfg *chatnewconfig.Config
}

func NewAppBuilder(cfg *chatnewconfig.Config) *ChatAppBuilder {
	return &ChatAppBuilder{
		cfg: cfg,
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
	return postgrespet.NewPGDB(builder.cfg.PG.DSN)
}

// repositories
// TODO:
func (builder *ChatAppBuilder) ClientProvider(brokerConn *custombroker.CBroker) chatusecases.IClientRepo {
	return broker.New(brokerConn)
}

func (builder *ChatAppBuilder) MessageProvider(brokerConn *custombroker.CBroker) chatusecases.IMessageRepo {
	return broker.New(brokerConn)
}

func (builder *ChatAppBuilder) MessageSaver(dbconn *postgrespet.PGDB) chatusecases.IMessageSaver {
	return postgresrepo.New(dbconn)
}

// business-logic
func (builder *ChatAppBuilder) Business(
	clientProvider chatusecases.IClientRepo,
	messageProvider chatusecases.IMessageRepo,
	defaultMessageSaver chatusecases.IMessageSaver,
) chatinterfaces.IBusiness {
	return chatusecases.New(
		builder.cfg,
		clientProvider, messageProvider,
		defaultMessageSaver,
	)
}

func (builder *ChatAppBuilder) Handler(business chatinterfaces.IBusiness) chat.ChatServer {
	return chatgrpc.New(business)
}

func (builder *ChatAppBuilder) BuildGRPCApp(handler chat.ChatServer) *appgrpc.App {
	return appgrpc.New(builder.cfg.GRPC.Host, builder.cfg.GRPC.Port, handler)
}
