package grpcclient

import (
	"context"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
)

func NewClient(ctx context.Context, cfg *config.Config) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		net.JoinHostPort(cfg.Grpc.Host, cfg.Grpc.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Failed to create grpc client")
		return nil, err
	}
	return conn, nil
}
