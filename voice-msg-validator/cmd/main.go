package main

import (
	"github.com/mefourr/tgdevob/msg/voice/validator/config"
	"github.com/mefourr/tgdevob/msg/voice/validator/internal/app"
	"github.com/mefourr/tgdevob/msg/voice/validator/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoadConfig()
	ctx := logger.Init()

	slog.InfoContext(ctx, "starting validator", "config", cfg)

	application := app.New(cfg)
	go application.GRPCSrv.MustRun(ctx)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	sig := <-shutdown
	slog.InfoContext(ctx, "received shutdown signal", "signal", sig.String())

	application.GRPCSrv.Shutdown(ctx)
	slog.InfoContext(ctx, "app has been shutdown")
}
