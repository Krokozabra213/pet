package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Krokozabra213/sso/configs/ssoconfig"
	"github.com/Krokozabra213/sso/internal/auth/app"
	apphttp "github.com/Krokozabra213/sso/internal/auth/app/http"
	"github.com/Krokozabra213/sso/pkg/logger"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
	EnvDev   = "dev"
)

func main() {
	test := true
	env := EnvLocal

	log := logger.SetupLogger(env)
	cfg := ssoconfig.Load(env, test)

	cfgJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Info("loaded config", slog.String("config", string(cfgJSON)))

	builder := app.NewAppBuilder(cfg, log)
	appFactory := app.NewAppFactory(builder)
	grpcApplication := appFactory.Create()
	go grpcApplication.MustRun()

	httpapp := apphttp.New(log, cfg.Server.HttpHost, cfg.Server.HttpPort, cfg.Server.Host, cfg.Server.Port)
	// go httpapp.MustRun()
	go httpapp.RunHTTP()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))

	httpapp.Stop()
	grpcApplication.Stop()
	log.Info("application stopped")
}
