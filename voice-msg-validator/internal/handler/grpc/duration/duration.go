package duration

import (
	"context"
	"github.com/mefourr/tgdevob/msg/voice/validator/internal/service/validation"
	"github.com/mefourr/tgdevob/proto/voice-msg-validator/pb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Validator interface {
	Validate(ctx context.Context, duration float32) error
}

type ServerAPI struct {
	pb.UnimplementedVoiceMessageDurationValidatorServer
	val Validator
}

func Register(gRPC *grpc.Server) {
	// TODO: need pass working interface into &ServerAPI{}
	pb.RegisterVoiceMessageDurationValidatorServer(gRPC, &ServerAPI{val: validation.New()})
}

const succeed bool = true

func (s *ServerAPI) Validate(ctx context.Context, in *pb.VoiceMessageDataRq) (*pb.ValidationResultRs, error) {
	if in.GetDuration() <= float32(0) {
		return nil, status.Error(codes.InvalidArgument, "duration cannot be less than 0")
	}

	if err := s.val.Validate(ctx, in.GetDuration()); err != nil {
		// TODO: error type handling
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &pb.ValidationResultRs{Ok: succeed}, nil
}
