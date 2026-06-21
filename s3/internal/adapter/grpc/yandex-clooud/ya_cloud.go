package yandex_clooud

import (
	"context"
	"github.com/mefourr/tgdevob/s3/pkg/grpcclient"
	"github.com/mefourr/tgdevob/s3/pkg/logger"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/storage/v1"
	"log/slog"
)

type YandexCloud interface {
	CreateBucketInCLoud(context.Context) (string, error)
}

type grpcToYandexCloud struct {
	yaclient *grpcclient.YandexCloudStorageClient
}

func (g grpcToYandexCloud) CreateBucketInCLoud(ctx context.Context) (string, error) {
	slog.InfoContext(ctx, "calling yandex cloud create bucket")
	op, err := g.yaclient.Client.Create(ctx, &storage.CreateBucketRequest{
		Name:     "yandex-cloud-storage-test",
		FolderId: "b1gqapgvrc4qove680qd",
		MaxSize:  1024,
	})
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "yandex cloud create bucket failed", "err", err)
		return "", err
	}
	slog.InfoContext(ctx, "yandex cloud operation created", "operation_id", op.Id)
	return op.Id, nil
}

func New(yaclient *grpcclient.YandexCloudStorageClient) YandexCloud {
	return &grpcToYandexCloud{yaclient: yaclient}
}
