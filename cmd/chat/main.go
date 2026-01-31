package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Krokozabra213/sso/internal/chat/app"
	chatnewconfig "github.com/Krokozabra213/sso/newconfigs/chat"
	custombroker "github.com/Krokozabra213/sso/pkg/custom-broker"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
	"github.com/Krokozabra213/sso/pkg/logger"
)

const (
	EnvLocal              = "local"
	EnvProd               = "prod"
	brokerShutdownTimeout = 20 * time.Second
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

	logger.Init(env)

	// TODO: Вынести в конфиг константы
	broker, err := custombroker.NewCBroker(2, 1000, 10)
	if err != nil {
		panic(err)
	}
	defer broker.GracefullShutdown(brokerShutdownTimeout)

	pg := postgrespet.NewPGDB(cfg.PG.DSN)

	appBuilder := app.NewAppBuilder(cfg, broker, pg)
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
