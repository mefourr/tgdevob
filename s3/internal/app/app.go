package app

import (
	"context"
	"github.com/mefourr/tgdevob/s3/config"
	"github.com/mefourr/tgdevob/s3/internal/app/cloud_storage_api"
	"github.com/mefourr/tgdevob/s3/pkg/grpcclient"
)

type App struct {
	S3  *cloud_storage_api.GrpcServer
	cfg *config.Config
}

func Init(ctx context.Context, cfg *config.Config, authClient *grpcclient.AuthClient) *App {
	return &App{
		S3: cloud_storage_api.New(
			ctx,
			cfg,
			authClient,
		),
		cfg: nil,
	}
}
