package main

import (
	"github.com/mefourr/tgdevob/s3/config"
	"github.com/mefourr/tgdevob/s3/internal/app"
	"github.com/mefourr/tgdevob/s3/pkg/grpcclient"
	"github.com/mefourr/tgdevob/s3/pkg/logger"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

/*
each usecase that uses external API must cache the authT-token
authT-usecase has only one grpc method to call GetIamToken() returns struct{token, expiredDate}
*/

//type IamTokenGenerator struct {
//	client auth.IamTokenGeneratorClient
//}

func main() {
	cfg := config.MustLoadConfig()
	ctx := logger.Init()

	authClient, err := grpcclient.NewAuthClient(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.InfoContext(ctx, "auth grpc client connection established")

	application := app.Init(ctx, cfg, authClient)

	// TODO: think of sync and graceful shutdown
	go func() {
		if err := application.S3.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	application.S3.Shutdown(ctx)
	//ctx := context.Background()
	////fixme(ctx)
	//authClient, err := grpc.NewAuthClient("localhost:5552", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	panic(err)
	//}
	//defer authClient.Close()
	//
	//i := IamTokenGenerator{
	//	client: authT.NewIamTokenGeneratorClient(authClient),
	//}
	//_, err = i.client.GetIamToken(ctx, &emptypb.Empty{})
	//if err != nil {
	//	panic(err)
	//}
}

//
//type S3BucketCreater struct {
//	client              storage.BucketServiceClient
//	folderServiceClient resourcemanager.FolderServiceClient
//}
//
//func fixme(ctx context.Context) {
//	conn, err := grpc.NewClient(
//		"storage.cloud_storage_api.cloud.yandex.net",
//		//"resource-manager.cloud_storage_api.cloud.yandex.net",
//		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
//		grpc.WithPerRPCCredentials(&ycTokenCredentials{token: "iamToken"}),
//	)
//	if err != nil {
//		panic(err)
//	}
//	defer conn.Close()
//
//	s3 := S3BucketCreater{
//		client:              storage.NewBucketServiceClient(conn),
//		folderServiceClient: resourcemanager.NewFolderServiceClient(conn),
//	}
//	op, err := s3.client.Create(ctx, &storage.CreateBucketRequest{
//		Name:     "yandex-cloud-storage-test",
//		FolderId: "b1gqapgvrc4qove680qd",
//		MaxSize:  1024,
//	})
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(op)
//}
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
