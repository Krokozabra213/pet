package appgrpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/Krokozabra213/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCApp struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	host       string
	port       string
}

func New(log *slog.Logger, host string, port string, srv sso.AuthServer) *GRPCApp {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(ValidationUnaryInterceptor),
	)

	sso.RegisterAuthServer(gRPCServer, srv)

	reflection.Register(gRPCServer)

	return &GRPCApp{
		log:        log,
		gRPCServer: gRPCServer,
		host:       host,
		port:       port,
	}
}

func (a *GRPCApp) RunGRPC() error {
	const op = "authappgrpc.Run"

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

func (a *GRPCApp) MustRun() {
	if err := a.RunGRPC(); err != nil {
		panic(err)
	}
}

func (a *GRPCApp) Stop() {
	const op = "authgrpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.String("address", a.host+a.port))

	a.gRPCServer.GracefulStop()
}
