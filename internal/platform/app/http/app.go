package apphttp

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	platformconfig "github.com/Krokozabra213/sso/newconfigs/platform"
	"github.com/Krokozabra213/sso/pkg/logger"
)

func Run(configfile, envfile string) {
	cfg, err := platformconfig.Init(configfile, envfile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)

	logger.Init(cfg.App.Environment)

	// TODO: START MONGO CLIENT
	// TODO: CONNECT MONGODB

	// TODO: ADD FILESTORAGE PROVIDER (MINIO)

	// TODO: ADD REPOSITORIES CONSTRUCTOR
	// TODO: ADD SERVICES CONSTRUCTOR

	// TODO: ADD HANDLER

	server := NewServer(cfg, nil)
	go server.MustRun()
	fmt.Println("fdfd")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	slog.Info("stopping application", slog.String("signal", sign.String()))

	if err := server.StopHTTP(); err != nil {
		slog.Error("failed to stop server", "err", err)
		return
	}

	// TODO: DISCONNECT MONGO клиента

	slog.Info("application stopped")
}
