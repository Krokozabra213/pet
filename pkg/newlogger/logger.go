package newlogger

import (
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

// use lvl=-4 for local lvl=0 for production
func SetupLogger(lvl int) *slog.Logger {
	var log *slog.Logger
	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(lvl)}),
	)

	return log
}
