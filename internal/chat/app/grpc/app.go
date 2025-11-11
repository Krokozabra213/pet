package chatappgrpc

import (
	"fmt"
	"log/slog"
	"net"

	chatgrpcserver "github.com/Krokozabra213/sso/internal/chat/grpc"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(log *slog.Logger, port string, auth chatgrpcserver.IBusiness) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(ValidationUnaryInterceptor),
		grpc.MaxConcurrentStreams(10_000), // Максимум одновременных стримов
		grpc.MaxRecvMsgSize(2*1024*1024),  // Максимальный размер сообщения 2MB
	)
	chatgrpcserver.Register(gRPCServer, auth)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
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
