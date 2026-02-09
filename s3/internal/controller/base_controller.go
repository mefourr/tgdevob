package controller

import (
	"context"
	"github.com/mefourr/tgdevob/proto/s3_storage/pb/s3_storage/v1"
	yandex_cloud "github.com/mefourr/tgdevob/s3/internal/adapter/grpc/yandex-clooud"
	"github.com/mefourr/tgdevob/s3/internal/usecase"
	"github.com/mefourr/tgdevob/s3/pkg/grpcclient"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"google.golang.org/grpc"
	"log/slog"
)

type CloudAPI interface {
	CreateBucket(context.Context) error
}

type ServerAPI struct {
	s3_storage.UnimplementedS3StorageServer
	cloud usecase.YandexCloud
}

func Register(gRPC *grpc.Server, client *grpcclient.YandexCloudStorageClient) {
	s3_storage.RegisterS3StorageServer(
		gRPC,
		&ServerAPI{
			cloud: usecase.New(yandex_cloud.New(client)),
		})
}

func (c *ServerAPI) StoreAudio(ctx context.Context, in *s3_storage.TempRq) (*s3_storage.TempRs, error) {
	slog.InfoContext(ctx, "StoreAudio got a call", "value", in.GetSmth())
	bucketId, err := c.cloud.CreateBucket(ctx)
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Error occurs while creating bucket", "err", err.Error())
		return nil, err
	}
	slog.InfoContext(ctx, "Bucket created", "bucketId", bucketId)
	return &s3_storage.TempRs{
		Smth: "Bucket has been created",
	}, nil
}
