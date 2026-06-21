package main

import (
	"github.com/mefourr/tgdevob/authentication/config"
	"github.com/mefourr/tgdevob/authentication/internal/app"
	"github.com/mefourr/tgdevob/authentication/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoadConfig()
	ctx := logger.Init()

	slog.InfoContext(ctx, "starting authentication service", "config", cfg)

	application := app.New(cfg)
	go application.Auth.MustRun(ctx)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	sig := <-shutdown
	slog.InfoContext(ctx, "received shutdown signal", "signal", sig.String())

	application.Auth.Shutdown(ctx)
	slog.InfoContext(ctx, "app has been shutdown")
}
