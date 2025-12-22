package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Krokozabra213/sso/internal/auth/app"
	apphttp "github.com/Krokozabra213/sso/internal/auth/app/http"
	ssonewconfig "github.com/Krokozabra213/sso/newconfigs/sso"
	"github.com/Krokozabra213/sso/pkg/logger"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

func main() {
	env := EnvLocal
	var configfile string

	switch env {
	case EnvLocal:
		configfile = "settings/sso_main.yml"
	case EnvProd:
		configfile = "settings/sso_prod.yml"
	}

	cfg, err := ssonewconfig.Init(configfile, "sso.env")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", cfg)

	// log := logger.SetupLogger(env)
	logger.Init(env)

	builder := app.NewAppBuilder(cfg)
	appFactory := app.NewAppFactory(builder)
	grpcApplication := appFactory.Create()
	go grpcApplication.MustRun()

	httpapp := apphttp.New(cfg.HTTP.Host, cfg.HTTP.Port, cfg.GRPC.Host, cfg.GRPC.Port)
	go httpapp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	slog.Info("stopping application", slog.String("signal", sign.String()))

	httpapp.Stop()
	grpcApplication.Stop()
	slog.Info("application stopped")
}
