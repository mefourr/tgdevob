package main

import (
	"context"
	"fmt"
	"github.com/mefourr/tgdevob/proto/voice-msg-validator/v1/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedValidateVMLengthServer
}

func (s *server) ValidateVMLength(ctx context.Context, in *pb.VoiceMessageDataRq) (*pb.ValidatedResultRs, error) {
	log.Printf("Received: %d, %d", in.GetDuration(), in.GetFileSize())
	return &pb.ValidatedResultRs{
		IsValidated: false,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5353))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterValidateVMLengthServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
