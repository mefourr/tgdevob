package generator

import (
	"context"
	"github.com/mefourr/tgdevob/authentication/internal/domain"
	"github.com/mefourr/tgdevob/authentication/pkg/logger"
	"github.com/mefourr/tgdevob/proto/auth/pb/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
)

type Generator interface {
	GetOrCreateToken(ctx context.Context) (*domain.Token, error)
}

type ServerAPI struct {
	auth.UnimplementedIamTokenGeneratorServer
	gen Generator
}

func Register(gRPC *grpc.Server, gen Generator) {
	auth.RegisterIamTokenGeneratorServer(gRPC, &ServerAPI{
		auth.UnimplementedIamTokenGeneratorServer{},
		gen,
	})
}

func (s *ServerAPI) GetIamToken(ctx context.Context, _ *emptypb.Empty) (*auth.GeneratorRs, error) {
	t, err := s.gen.GetOrCreateToken(ctx)
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "error getting token", "err", err)
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &auth.GeneratorRs{
		Token:     t.IamToken,
		ExpiresAt: timestamppb.New(t.ExpiresAt),
	}, nil
}
