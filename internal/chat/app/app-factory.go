package app

import (
	"github.com/Krokozabra213/protos/gen/go/chat"
	appgrpc "github.com/Krokozabra213/sso/internal/chat/app/grpc"
	chatusecases "github.com/Krokozabra213/sso/internal/chat/business/usecases"
	chatinterfaces "github.com/Krokozabra213/sso/internal/chat/grpc/interfaces"
)

type IAppBuilder interface {

	// Connect
	// BrokerConn() *custombroker.CBroker
	// PGConn() *postgrespet.PGDB

	// Repository providers
	ClientProvider() chatusecases.IClientRepo
	MessageProvider() chatusecases.IMessageRepo
	MessageSaver() chatusecases.IMessageSaver

	// Business logic layer
	Business(
		clientProvider chatusecases.IClientRepo,
		messageProvider chatusecases.IMessageRepo,
		MessageSaver chatusecases.IMessageSaver,
	) chatinterfaces.IBusiness

	// gRPC handler
	Handler(business chatinterfaces.IBusiness) chat.ChatServer

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

	// brokerConn := fact.builder.BrokerConn()
	// pgConn := fact.builder.PGConn()

	// Repositories
	clientRepo := fact.builder.ClientProvider()
	messageRepo := fact.builder.MessageProvider()
	messageSaver := fact.builder.MessageSaver()

	// Business-Logic
	business := fact.builder.Business(
		clientRepo, messageRepo, messageSaver,
	)

	// Handler
	handler := fact.builder.Handler(business)

	// Application
	application := fact.builder.BuildGRPCApp(handler)

	return application
}
