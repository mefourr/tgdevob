package vmlength

import (
	"context"
	"github.com/mefourr/tgdevob/proto/voice-msg-validator/pb/v1"
	"google.golang.org/grpc"
	"time"
)

type GRPCProducer interface {
	produce(ctx context.Context, duration time.Duration) (bool, error)
}

type durationValidator struct {
	client pb.ValidateVMLengthClient
}

func New(conn *grpc.ClientConn) GRPCProducer {
	return &durationValidator{client: pb.NewValidateVMLengthClient(conn)}
}

func (g *durationValidator) produce(_ context.Context, _ time.Duration) (bool, error) {
	//TODO implement me
	panic("implement me")
}
