package duration

import (
	"context"
	"errors"
	"github.com/mefourr/tgdevob/msg/voice/validator/internal/service/validation"
	"github.com/mefourr/tgdevob/proto/voice-msg-validator/pb/v1"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"log/slog"
)

type Validator interface {
	Validate(ctx context.Context, duration *durationpb.Duration, fileSize int64) (err error, res bool)
}

type ServerAPI struct {
	pb.UnimplementedValidateVMLengthServer
	val Validator
}

func Register(gRPC *grpc.Server) {
	// TODO: need pass working interface into &ServerAPI{}
	pb.RegisterValidateVMLengthServer(gRPC, &ServerAPI{val: validation.New()})
}

const (
	emptyValue = 0
)

func (s *ServerAPI) ValidateVMLength(ctx context.Context, in *pb.VoiceMessageDataRq) (*pb.ValidatedResultRs, error) {
	if in.GetDu() == nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, errors.New("duration cannot be nil")), "request fields validation failed")
		return nil, status.Error(codes.InvalidArgument, "duration cannot be nil")
	}

	if in.GetFileSize() == emptyValue {
		slog.ErrorContext(logging.ErrorCtx(ctx, errors.New("filesize cannot be zero")), "request fields validation failed")
		return nil, status.Error(codes.InvalidArgument, "filesize cannot be zero")
	}

	// TODO: implement voice message validation service
	err, res := s.val.Validate(ctx, in.GetDu(), in.GetFileSize())
	if err != nil {
		// TODO: error type handling
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &pb.ValidatedResultRs{
		IsValidated: res,
	}, nil
}
