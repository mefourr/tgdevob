package usecase

import (
	"context"
	"errors"
	"github.com/mefourr/tgdevob/worker/internal/adapter/grpc/validator"
	"time"
)

var ErrDurationZero = errors.New("generator cannot be zero")

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type GRPCProducer interface {
	Send(ctx context.Context, duration float32) error
}

type Validator interface {
	Validate(ctx context.Context, duration time.Duration) error
}

type durationValidator struct {
	grpc validator.GRPCProducer
}

func New(grpc validator.GRPCProducer) Validator {
	return &durationValidator{
		grpc: grpc,
	}
}
