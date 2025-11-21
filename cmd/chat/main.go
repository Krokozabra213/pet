package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Krokozabra213/sso/configs/chatconfig"
	"github.com/Krokozabra213/sso/internal/chat/app"
	"github.com/Krokozabra213/sso/pkg/logger"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
	EnvDev   = "dev"
)

func main() {
	env := EnvLocal
	test := true

	log := logger.SetupLogger(env)
	cfg := chatconfig.Load(env, test)

	cfgJSON, _ := json.MarshalIndent(cfg, "", "  ")
	log.Info("loaded config", slog.String("config", string(cfgJSON)))

	appBuilder := app.NewAppBuilder(cfg, log)
	appFactory := app.NewAppFactory(appBuilder)
	application := appFactory.Create()
	go application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))

	application.Stop()

	log.Info("application stopped")
}
