package main

import (
	"github.com/mefourr/tgdevob/msg/voice/validator/config"
	"github.com/mefourr/tgdevob/msg/voice/validator/internal/app"
	"github.com/mefourr/tgdevob/msg/voice/validator/pkg/logging"
	"log/slog"
)

func main() {
	cfg := config.LoadConfig()
	ctx := logging.Init()

	slog.InfoContext(ctx, "starting validator", "config", cfg)

	application := app.New(cfg)
	application.GRPCSrv.MustRun(ctx)
}
