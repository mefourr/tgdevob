package duration

import (
	"context"
	"github.com/mefourr/tgdevob/proto/voice-msg-validator/pb/v1"
	"google.golang.org/grpc"
	"log"
)

type ServerAPI struct {
	pb.UnimplementedValidateVMLengthServer
}

func Register(gRPC *grpc.Server) {
	pb.RegisterValidateVMLengthServer(gRPC, &ServerAPI{})
}

func (s *ServerAPI) ValidateVMLength(_ context.Context, in *pb.VoiceMessageDataRq) (*pb.ValidatedResultRs, error) {
	log.Printf("Received: %d, %d", in.GetDuration(), in.GetFileSize())
	return &pb.ValidatedResultRs{
		IsValidated: false,
	}, nil
}
