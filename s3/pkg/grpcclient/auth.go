package grpcclient

import (
	"context"
	"github.com/mefourr/tgdevob/proto/auth/pb/auth/v1"
	"github.com/mefourr/tgdevob/s3/config"
	"github.com/mefourr/tgdevob/s3/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"strconv"
)

type AuthClient struct {
	Conn   *grpc.ClientConn
	Client auth.IamTokenGeneratorClient
}

func NewAuthClient(ctx context.Context, cfg *config.Config) (*AuthClient, error) {
	conn, err := grpc.NewClient(
		net.JoinHostPort(cfg.Grpc.Auth.Client.Host, strconv.Itoa(cfg.Grpc.Auth.Client.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Failed to create grpc authClient")
		return nil, err
	}

	return &AuthClient{
		Conn:   conn,
		Client: auth.NewIamTokenGeneratorClient(conn),
	}, nil
}
