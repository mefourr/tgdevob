package yandex_clooud

import (
	"context"
	"github.com/mefourr/tgdevob/s3/pkg/grpcclient"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/storage/v1"
)

type YandexCloud interface {
	CreateBucketInCLoud(context.Context) (string, error)
}

type grpcToYandexCloud struct {
	yaclient *grpcclient.YandexCloudStorageClient
}

func (g grpcToYandexCloud) CreateBucketInCLoud(ctx context.Context) (string, error) {
	//TODO implement me
	op, err := g.yaclient.Client.Create(ctx, &storage.CreateBucketRequest{
		Name:     "yandex-cloud-storage-test",
		FolderId: "b1gqapgvrc4qove680qd",
		MaxSize:  1024,
	})
	if err != nil {
		return "", err
	}

	return op.Id, nil
}

func New(yaclient *grpcclient.YandexCloudStorageClient) YandexCloud {
	return &grpcToYandexCloud{yaclient: yaclient}
}
