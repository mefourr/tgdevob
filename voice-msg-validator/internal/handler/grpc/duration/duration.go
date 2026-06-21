package duration

import (
	"context"
	"github.com/mefourr/tgdevob/msg/voice/validator/internal/service/validation"
	"github.com/mefourr/tgdevob/msg/voice/validator/pkg/logger"
	"github.com/mefourr/tgdevob/proto/voice-msg-validator/pb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type Validator interface {
	Generate(ctx context.Context, duration float32) error
}

type ServerAPI struct {
	pb.UnimplementedVoiceMessageDurationValidatorServer
	val Validator
}

func Register(gRPC *grpc.Server) {
	// TODO: need pass working interface into &ServerAPI{}
	pb.RegisterVoiceMessageDurationValidatorServer(
		gRPC,
		&ServerAPI{val: validation.New()},
	)
}

const succeed bool = true

func (s *ServerAPI) Validate(ctx context.Context, in *pb.VoiceMessageDataRq) (*pb.ValidationResultRs, error) {
	slog.InfoContext(ctx, "validate request received", "duration", in.GetDuration())

	if in.GetDuration() <= float32(0) {
		slog.WarnContext(ctx, "invalid duration", "duration", in.GetDuration())
		return nil, status.Error(codes.InvalidArgument, "duration cannot be less than or equal to 0")
	}

	if err := s.val.Generate(ctx, in.GetDuration()); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "validation failed", "duration", in.GetDuration(), "err", err)
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	slog.InfoContext(ctx, "validation passed", "duration", in.GetDuration())
	return &pb.ValidationResultRs{Ok: succeed}, nil
}
