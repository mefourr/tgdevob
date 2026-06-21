package validator

import (
	"context"
	"github.com/mefourr/tgdevob/proto/voice-msg-validator/pb/v1"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"google.golang.org/grpc"
	"log/slog"
)

type GRPCProducer interface {
	Send(ctx context.Context, duration float32) error
}

type durationValidator struct {
	client pb.VoiceMessageDurationValidatorClient
}

func New(conn *grpc.ClientConn) GRPCProducer {
	return &durationValidator{client: pb.NewVoiceMessageDurationValidatorClient(conn)}
}

func (d *durationValidator) Send(ctx context.Context, duration float32) error {
	slog.DebugContext(ctx, "sending duration validation request", "duration", duration)
	_, err := d.client.Validate(ctx, &pb.VoiceMessageDataRq{
		Duration: duration,
	})
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "duration validation request failed", "duration", duration, "err", err)
		return err
	}
	slog.DebugContext(ctx, "duration validation passed", "duration", duration)
	return nil
}
