package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Krokozabra213/sso/internal/chat/app"
	chatnewconfig "github.com/Krokozabra213/sso/newconfigs/chat"
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
		configfile = "settings/chat_main.yml"
	case EnvProd:
		configfile = "settings/chat_prod.yml"
	}

	cfg, err := chatnewconfig.Init(configfile, "chat.env")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", cfg)

	// log := logger.SetupLogger(env)
	logger.Init(env)

	appBuilder := app.NewAppBuilder(cfg)
	appFactory := app.NewAppFactory(appBuilder)
	application := appFactory.Create()
	go application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	slog.Info("stopping application", slog.String("signal", sign.String()))

	application.Stop()

	slog.Info("application stopped")
}
