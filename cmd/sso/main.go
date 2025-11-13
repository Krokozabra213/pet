package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Krokozabra213/sso/configs/ssoconfig"
	"github.com/Krokozabra213/sso/internal/auth/app"
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
	log.Info("config", slog.String("ssoconfig", fmt.Sprintf("%#v", cfg)))

	builder := app.NewAppBuilder(cfg, log)
	appFactory := app.NewAppFactory(builder)
	application := appFactory.Create()
	go application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))

	application.Stop()
	log.Info("application stopped")
}
