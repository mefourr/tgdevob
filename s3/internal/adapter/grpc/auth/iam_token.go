package auth

import (
	"context"
	"github.com/mefourr/tgdevob/proto/auth/pb/auth/v1"
	"github.com/mefourr/tgdevob/s3/internal/domain"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type YandexAuth interface {
	GetOrRequestAuthToken(context.Context) (*domain.Token, error)
}

type grpcYaTokenProvider struct {
	client auth.IamTokenGeneratorClient
}

func (g *grpcYaTokenProvider) GetOrRequestAuthToken(ctx context.Context) (*domain.Token, error) {
	slog.InfoContext(ctx, "ready to get token from IamTokenGeneratorClient")
	rs, err := g.client.GetIamToken(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	slog.InfoContext(ctx, "Expiration time of generated token", "expires_at", rs.GetExpiresAt().AsTime())
	return &domain.Token{
		IamToken:  rs.GetToken(),
		ExpiresAt: rs.GetExpiresAt().AsTime(),
	}, nil
}

func New(client auth.IamTokenGeneratorClient) YandexAuth {
	return &grpcYaTokenProvider{client: client}
}
