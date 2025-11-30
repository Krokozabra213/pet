package app

import (
	"github.com/Krokozabra213/protos/gen/go/proto/chat"
	appgrpc "github.com/Krokozabra213/sso/internal/chat/app/grpc"
	chatBusiness "github.com/Krokozabra213/sso/internal/chat/business"
	chatgrpc "github.com/Krokozabra213/sso/internal/chat/grpc"
	custombroker "github.com/Krokozabra213/sso/pkg/custom-broker"
)

type IAppBuilder interface {

	// Connect
	BrokerConn() *custombroker.CBroker

	// Repository providers
	ClientProvider(brokerConn *custombroker.CBroker) chatBusiness.IClientRepo
	MessageProvider(brokerConn *custombroker.CBroker) chatBusiness.IMessageRepo

	// Business logic layer
	Business(
		clientProvider chatBusiness.IClientRepo,
		messageProvider chatBusiness.IMessageRepo,
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

	// Repositories
	ClientRepo := fact.builder.ClientProvider(brokerConn)
	MessageRepo := fact.builder.MessageProvider(brokerConn)

	// Business-Logic
	business := fact.builder.Business(
		ClientRepo, MessageRepo,
	)

	// Handler
	handler := fact.builder.Handler(business)

	// Application
	application := fact.builder.BuildGRPCApp(handler)

	return application
}
