package app

import (
	"github.com/mefourr/tgdevob/authentication/config"
	"github.com/mefourr/tgdevob/authentication/internal/app/auth"
)

type App struct {
	Auth *auth.App
	cfg  *config.Config
}

func New(cfg *config.Config) *App {
	return &App{
		Auth: auth.New(cfg.Grpc.Port),
		cfg:  cfg,
	}
}
