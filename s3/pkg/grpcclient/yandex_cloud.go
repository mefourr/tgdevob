package grpcclient

import (
	"context"
	"github.com/mefourr/tgdevob/s3/config"
	"github.com/mefourr/tgdevob/s3/internal/adapter/grpc/auth"
	"github.com/mefourr/tgdevob/s3/pkg/logger"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/storage/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"net"
	"strconv"
)

type YandexCloudStorageClient struct {
	Conn   *grpc.ClientConn
	Client storage.BucketServiceClient
}

// NewYandexCloudStorageClient TODO: think of how to dial Conn when token gets from Auth service
func NewYandexCloudStorageClient(ctx context.Context, cfg *config.Config, auth auth.YandexAuth) (*YandexCloudStorageClient, error) {
	conn, err := grpc.NewClient(
		net.JoinHostPort(cfg.Grpc.YandexCloud.Client.Host, strconv.Itoa(cfg.Grpc.YandexCloud.Client.Port)),
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
		grpc.WithUnaryInterceptor(authIntercept(auth)),
	)
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Failed to create grpc authClient")
		return nil, err
	}

	return &YandexCloudStorageClient{
		Conn:   conn,
		Client: storage.NewBucketServiceClient(conn),
	}, nil
}

func authIntercept(auth auth.YandexAuth) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		token, err := auth.GetOrRequestAuthToken(ctx)
		if err != nil {
			slog.ErrorContext(logger.ErrorCtx(ctx, err), "Failed to get or request auth token")
			return err
		}

		md := metadata.Pairs("authorization", "Bearer "+token.IamToken)
		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

//
//type ycTokenCredentials struct {
//	token string
//}
//
//func (c *ycTokenCredentials) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
//	return map[string]string{
//		"authorization": "Bearer " + c.token,
//	}, nil
//}
//
//func (c *ycTokenCredentials) RequireTransportSecurity() bool {
//	return true
//}
