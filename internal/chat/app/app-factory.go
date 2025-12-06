package app

import (
	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	appgrpc "github.com/Krokozabra213/sso/internal/chat/app/grpc"
	chatBusiness "github.com/Krokozabra213/sso/internal/chat/business"
	chatgrpc "github.com/Krokozabra213/sso/internal/chat/grpc"
	custombroker "github.com/Krokozabra213/sso/pkg/custom-broker"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
)

type IAppBuilder interface {

	// Connect
	BrokerConn() *custombroker.CBroker
	PGConn() *postgrespet.PGDB

	// Repository providers
	ClientProvider(brokerConn *custombroker.CBroker) chatBusiness.IClientRepo
	MessageProvider(brokerConn *custombroker.CBroker) chatBusiness.IMessageRepo
	DefaultMessageSaver(dbconn *postgrespet.PGDB) chatBusiness.IDefaultMessageSaver

	// Business logic layer
	Business(
		clientProvider chatBusiness.IClientRepo,
		messageProvider chatBusiness.IMessageRepo,
		defaultMessageSaver chatBusiness.IDefaultMessageSaver,
	) chatgrpc.IBusiness

	// gRPC handler
	Handler(business chatgrpc.IBusiness) chat.ChatServer

	// Application builder
	BuildGRPCApp(handler chat.ChatServer) *appgrpc.App
}

type AppFactory struct {
	builder IAppBuilder
}

func NewAppFactory(builder IAppBuilder) *AppFactory {
	return &AppFactory{
		builder: builder,
	}
}

func (fact *AppFactory) Create() *appgrpc.App {

	brokerConn := fact.builder.BrokerConn()
	pgConn := fact.builder.PGConn()

	// Repositories
	clientRepo := fact.builder.ClientProvider(brokerConn)
	messageRepo := fact.builder.MessageProvider(brokerConn)
	defaultMessageSaver := fact.builder.DefaultMessageSaver(pgConn)

	// Business-Logic
	business := fact.builder.Business(
		clientRepo, messageRepo, defaultMessageSaver,
	)

	// Handler
	handler := fact.builder.Handler(business)

	// Application
	application := fact.builder.BuildGRPCApp(handler)

	return application
}
