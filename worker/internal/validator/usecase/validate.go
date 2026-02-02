package usecase

import (
	"context"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"log/slog"
	"time"
)

func (v *durationValidator) Validate(ctx context.Context, duration time.Duration) error {
	if duration == time.Duration(0) {
		return ErrDurationZero
	}

	if err := v.grpc.Send(ctx, float32(duration)); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "error while validating generator", "generator", duration)
		return err
	}
	return nil
}
