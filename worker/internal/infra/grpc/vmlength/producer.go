package vmlength

import (
	"context"
	"github.com/mefourr/tgdevob/proto/voice-msg-validator/pb/v1"
	"google.golang.org/grpc"
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
	_, err := d.client.Validate(ctx, &pb.VoiceMessageDataRq{
		Duration: duration,
	})
	if err != nil {
		return err
	}

	return nil
}
