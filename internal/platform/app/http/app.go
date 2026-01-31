package apphttp

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Krokozabra213/sso/internal/platform/business"
	ssoclient "github.com/Krokozabra213/sso/internal/platform/clients/sso"
	httpPlatform "github.com/Krokozabra213/sso/internal/platform/http"
	"github.com/Krokozabra213/sso/internal/platform/repository"
	platformconfig "github.com/Krokozabra213/sso/newconfigs/platform"
	"github.com/Krokozabra213/sso/pkg/db/mongodb"
	jwtv1 "github.com/Krokozabra213/sso/pkg/jwt-manager/v1"
	keymanagerv1 "github.com/Krokozabra213/sso/pkg/key-manager/v1"
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
	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		panic(err)
	}
	defer mongodb.Disconnect(mongoClient)

	db := mongoClient.Database(cfg.Mongo.Name)

	// TODO: ADD FILESTORAGE PROVIDER (MINIO)

	ssoClient, err := ssoclient.NewClient(cfg.SSOConfig.Timeout, cfg.SSOServiceAddress(), cfg.App.AppID)
	if err != nil {
		panic(err)
	}
	publickeyPEM, err := ssoClient.GetPublicKey(context.Background())
	if err != nil {
		panic(err)
	}
	publicManager, err := keymanagerv1.NewPublicManager(publickeyPEM)
	if err != nil {
		panic(err)
	}

	jwtValidator, err := jwtv1.NewValidator(publicManager.GetPublicKey())
	if err != nil {
		panic(err)
	}
	repositories := repository.NewRepositories(db)
	business := business.New(business.Deps{
		Config: cfg,
		Repos:  repositories,
	})
	handler := httpPlatform.NewHandler(business, jwtValidator)

	server := NewServer(cfg, handler.Init(cfg))
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
