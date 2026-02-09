package usecase

import (
	"context"
	yandex_cloud "github.com/mefourr/tgdevob/s3/internal/adapter/grpc/yandex-clooud"
)

type yandexCloud struct {
	cloud yandex_cloud.YandexCloud
}

func New(cloud yandex_cloud.YandexCloud) YandexCloud {
	return &yandexCloud{cloud: cloud}
}

func (y *yandexCloud) CreateBucket(ctx context.Context) (string, error) {
	return y.cloud.CreateBucketInCLoud(ctx)
}
