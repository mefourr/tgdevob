package main

import (
	"context"
	"fmt"
	"github.com/mefourr/tgdevob/msg/voice/validator/pb/message"
	"github.com/mefourr/tgdevob/msg/voice/validator/pb/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	service.UnimplementedEchoServiceServer
}

func (s *server) Echo(ctx context.Context, in *message.StringMessage) (*message.StringMessage, error) {
	log.Printf("Received: %v", in.GetValue())
	return &message.StringMessage{Value: in.Value + " processed"}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5353))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	service.RegisterEchoServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
