package app

import (
	"github.com/mefourr/tgdevob/msg/voice/validator/config"
	"github.com/mefourr/tgdevob/msg/voice/validator/internal/app/grpcapp"
)

type App struct {
	GRPCSrv *grpcapp.App
	cfg     *config.Config
}

func New(cfg *config.Config) *App {
	return &App{
		GRPCSrv: grpcapp.New(cfg.Grpc.Port),
		cfg:     cfg,
	}
}
