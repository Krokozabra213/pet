package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Krokozabra213/sso/configs/ssoconfig"
	"github.com/Krokozabra213/sso/internal/auth/app"
	"github.com/Krokozabra213/sso/pkg/logger"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

var (
	host     = flag.String("host", "", "host for connect db")
	user     = flag.String("user", "", "username for connect db")
	password = flag.String("password", "", "password for connect db")
	dbname   = flag.String("dbname", "", "dbname for connect db")
	port     = flag.String("port", "", "port for connect db")
	sslmode  = flag.String("sslmode", "", "sslmode for connect db")
)

func main() {
	flag.Parse()
	// dsn := fmt.Sprintf(
	// 	"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
	// 	*host, *user, *password, *dbname, *port, *sslmode,
	// )
	cfg := ssoconfig.Load(envLocal)
	log := logger.SetupLogger(envLocal)

	application := app.New(log, cfg)
	// application := app.New(log, 44044, "", 15*time.Minute, 10_000*time.Minute, dsn, "testSecret")

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()

	log.Info("application stopped")
}
