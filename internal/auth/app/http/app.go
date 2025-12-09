package apphttp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Krokozabra213/protos/gen/go/sso"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ReadHeaderTimeout = 10 * time.Second
	ShutDownTimeout   = 5 * time.Second
)

type HTTPApp struct {
	log      *slog.Logger
	gwServer *http.Server
	host     string
	port     string
	grpcHost string
	grpcPort string
}

func New(
	log *slog.Logger, host string, port string,
	grpcHost string, grpcPort string,
) *HTTPApp {
	return &HTTPApp{
		log:      log,
		host:     host,
		port:     port,
		grpcHost: grpcHost,
		grpcPort: grpcPort,
	}
}

func (a *HTTPApp) RunHTTP() error {
	const op = "authapphttp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("host", a.host),
		slog.String("port", a.port),
		slog.String("grpcHost", a.grpcHost),
		slog.String("grpcPort", a.grpcPort),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := sso.RegisterAuthHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("%s:%s", a.grpcHost, a.grpcPort),
		opts,
	)
	if err != nil {
		log.Error("failed to register gateway", "err", err.Error())
	}

	a.gwServer = &http.Server{
		Addr:              fmt.Sprintf("%s:%s", a.host, a.port),
		Handler:           mux,
		ReadHeaderTimeout: ReadHeaderTimeout,
	}

	log.Info("http server with grpc-gateway listening...", slog.String("addr", a.gwServer.Addr))
	err = a.gwServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("failed to serve HTTP", "err", err.Error())
		return err
	}

	return nil
}

func (a *HTTPApp) MustRun() {
	if err := a.RunHTTP(); err != nil {
		panic(err)
	}
}

func (a *HTTPApp) Stop() {
	const op = "authapphttp.Stop"

	log := a.log.With(
		slog.String("op", op),
		slog.String("host", a.host),
		slog.String("port", a.port),
	)

	if a.gwServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutDownTimeout)
		defer cancel()
		if err := a.gwServer.Shutdown(shutdownCtx); err != nil {
			log.Error("failed shutdown http server", "err", err.Error())
			return
		}
		log.Info("http server stopped")
	}
}
