package cloud_storage_api

import (
	"context"
	"fmt"
	"github.com/mefourr/tgdevob/proto/auth/pb/auth/v1"
	"github.com/mefourr/tgdevob/s3/config"
	auth_adapter "github.com/mefourr/tgdevob/s3/internal/adapter/grpc/auth"
	"github.com/mefourr/tgdevob/s3/internal/controller"
	"github.com/mefourr/tgdevob/s3/pkg/grpcclient"
	"github.com/mefourr/tgdevob/s3/pkg/logger"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type GrpcServer struct {
	srv          *grpc.Server
	port         int
	auth         *grpcclient.AuthClient
	cloudStorage *grpcclient.YandexCloudStorageClient
}

func New(
	ctx context.Context,
	cfg *config.Config,
	authClient *grpcclient.AuthClient,
) *GrpcServer {
	cloudStorage, err := grpcclient.NewYandexCloudStorageClient(
		ctx,
		cfg,
		auth_adapter.New(
			auth.NewIamTokenGeneratorClient(authClient.Conn),
		),
	)
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to create cloud storage client")
		panic(err)
	}
	slog.InfoContext(ctx, "cloud storage grpc client connection established")

	srv := grpc.NewServer()
	controller.Register(srv, cloudStorage)

	return &GrpcServer{
		srv:          srv,
		port:         cfg.Grpc.Server.Port,
		auth:         authClient,
		cloudStorage: cloudStorage,
	}
}

func (g *GrpcServer) Run(ctx context.Context) error {
	// TODO: add controller logger attribute to slog

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to listen:", "err", err)
		return err
	}

	slog.InfoContext(ctx, "server listening at", "addr", lis.Addr())

	if err := g.srv.Serve(lis); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to serve:", "err", err)
		return err
	}

	return nil
}

func (g *GrpcServer) Shutdown(ctx context.Context) {
	if err := g.auth.Conn.Close(); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to shutdown:", "err", err)
		return
	}
	if err := g.cloudStorage.Conn.Close(); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to shutdown:", "err", err)
		return
	}
}
