package usecase

import (
	"context"
	"errors"
	"github.com/mefourr/tgdevob/worker/internal/validator/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_voiceDurationValidator_Execute_Grpc_Error(t *testing.T) {
	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "duration_limit_exceeded",
			args: args{time.Duration(99999)},
		},
	}
	//mock.Anything
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grpc := mocks.NewMockGRPCProducer(t)
			voiceDurationValidation := New(grpc)

			grpc.On("Send", mock.Anything, float32(tt.args.duration)).
				Return(errors.New("some errors")).
				Once()

			err := voiceDurationValidation.Validate(context.Background(), tt.args.duration)
			assert.NotNil(t, err)
		})
	}
}

func Test_voiceDurationValidator_Execute_Duration_Zero(t *testing.T) {
	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "duration_zero",
			args: args{time.Duration(0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			voiceDurationValidation := New(nil)
			err := voiceDurationValidation.Validate(context.Background(), tt.args.duration)
			assert.NotNil(t, err)
			assert.Error(t, err, ErrDurationZero.Error())
		})
	}
}

func Test_voiceDurationValidator_Execute_Happy_Path(t *testing.T) {
	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "base case",
			args: args{time.Duration(1)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grpc := mocks.NewMockGRPCProducer(t)
			voiceDurationValidation := New(grpc)

			grpc.On("Send", mock.Anything, float32(tt.args.duration)).
				Return(nil).
				Once()

			err := voiceDurationValidation.Validate(context.Background(), tt.args.duration)
			assert.Nil(t, err)
		})
	}
}
