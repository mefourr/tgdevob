package service

import (
	"context"
	"errors"
	"github.com/mefourr/tgdevob/worker/internal/infra/grpc/vmlength"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"log/slog"
	"time"
)

var ErrDurationZero = errors.New("duration cannot be zero")

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type GRPCProducer interface {
	Send(ctx context.Context, duration float32) error
}

type VoiceDurationValidator interface {
	Execute(ctx context.Context, duration time.Duration) error
}

type voiceDurationValidator struct {
	grpc vmlength.GRPCProducer
}

func NewVoiceDurationValidator(grpc vmlength.GRPCProducer) VoiceDurationValidator {
	return &voiceDurationValidator{
		grpc: grpc,
	}
}

func (v *voiceDurationValidator) Execute(ctx context.Context, duration time.Duration) error {
	if duration == time.Duration(0) {
		return ErrDurationZero
	}

	if err := v.grpc.Send(ctx, float32(duration)); err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "error while validating duration", "duration", duration)
		return err
	}
	return nil
}
