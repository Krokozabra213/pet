package appgrpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/Krokozabra213/protos/gen/go/chat"
	"google.golang.org/grpc"
)

const (
	MaxParallelStreams = 10_000
	MaxRecvMsgSize     = 2 * 1024 * 1024 // 2MB
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	host       string
	port       string
}

func New(log *slog.Logger, host string, port string, srv chat.ChatServer) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(ValidationUnaryInterceptor),
		grpc.MaxConcurrentStreams(MaxParallelStreams),
		grpc.MaxRecvMsgSize(MaxRecvMsgSize),
	)

	chat.RegisterChatServer(gRPCServer, srv)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		host:       host,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("host", a.host),
		slog.String("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.host, a.port))
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.String("port", a.port))

	a.gRPCServer.GracefulStop()
}
