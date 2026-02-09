package suite

import (
	"context"
	"github.com/mefourr/tgdevob/proto/s3_storage/pb/s3_storage/v1"
	"github.com/mefourr/tgdevob/s3/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg             *config.Config
	S3StorageClient s3_storage.S3StorageClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.LoadConfig("./config")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Grpc.Server.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	conn, err := grpc.NewClient(
		net.JoinHostPort(cfg.Grpc.Server.Host, strconv.Itoa(cfg.Grpc.Server.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Grpc server connection faild: %s", err)
		return nil, nil
	}

	return ctx, &Suite{
		T:               t,
		Cfg:             cfg,
		S3StorageClient: s3_storage.NewS3StorageClient(conn),
	}
}
