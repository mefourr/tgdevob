package suite

import (
	"context"
	"github.com/mefourr/tgdevob/msg/voice/validator/config"
	"github.com/mefourr/tgdevob/proto/voice-msg-validator/pb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	VVMLClient pb.VoiceMessageDurationValidatorClient
}

const (
	host = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.LoadConfig("./")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Grpc.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	conn, err := grpc.NewClient(
		net.JoinHostPort(host, strconv.Itoa(cfg.Grpc.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Grpc server connection faild: %s", err)
		return nil, nil
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		VVMLClient: pb.NewVoiceMessageDurationValidatorClient(conn),
	}
}
