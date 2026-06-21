package usecase

import (
	"context"
	yandex_cloud "github.com/mefourr/tgdevob/s3/internal/adapter/grpc/yandex-clooud"
	"github.com/mefourr/tgdevob/s3/pkg/logger"
	"log/slog"
)

type yandexCloud struct {
	cloud yandex_cloud.YandexCloud
}

func New(cloud yandex_cloud.YandexCloud) YandexCloud {
	return &yandexCloud{cloud: cloud}
}

func (y *yandexCloud) CreateBucket(ctx context.Context) (string, error) {
	slog.InfoContext(ctx, "creating bucket")
	id, err := y.cloud.CreateBucketInCLoud(ctx)
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to create bucket", "err", err)
		return "", err
	}
	slog.InfoContext(ctx, "bucket created", "bucket_id", id)
	return id, nil
}
