package usecase

import (
	"context"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"log/slog"
	"time"
)

func (v *durationValidator) Validate(ctx context.Context, duration time.Duration) error {
	slog.DebugContext(ctx, "validating voice message duration", "duration", duration)

	if duration == time.Duration(0) {
		slog.WarnContext(ctx, "voice message duration is zero, rejecting")
		return ErrDurationZero
	}

	if err := v.grpc.Send(ctx, float32(duration)); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "duration validation failed", "duration", duration, "err", err)
		return err
	}

	slog.DebugContext(ctx, "duration validation passed", "duration", duration)
	return nil
}
