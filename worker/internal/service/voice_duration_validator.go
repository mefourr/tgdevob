package service

import (
	"context"
	"github.com/mefourr/tgdevob/worker/internal/infra/grpc/vmlength"
	"time"
)

type VoiceDurationValidator interface {
	Execute(ctx context.Context, duration time.Duration) (bool, error)
}

type voiceDurationValidator struct {
	producer vmlength.GRPCProducer
}

func NewVoiceDurationValidator(producer vmlength.GRPCProducer) VoiceDurationValidator {
	return &voiceDurationValidator{
		producer: producer,
	}
}

func (v *voiceDurationValidator) Execute(_ context.Context, _ time.Duration) (bool, error) {
	//TODO implement me
	panic("implement me")
}
