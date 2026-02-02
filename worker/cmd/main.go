package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/app"
	"github.com/mefourr/tgdevob/worker/pkg/grpcclient"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"github.com/mefourr/tgdevob/worker/pkg/postgres"
	redisclient "github.com/mefourr/tgdevob/worker/pkg/redis"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoadConfig()
	ctx := logger.Init()

	client := connectToRedis(ctx, cfg)
	slog.InfoContext(ctx, "Redis connection established")

	pool := connectToPostgres(ctx, cfg)
	slog.InfoContext(ctx, "Postgres connection established")

	conn := grpcClientConnection(ctx, cfg)
	slog.InfoContext(ctx, "GRPC server connection established")

	application := app.New(cfg, client, pool, conn)

	ctx, cancel := context.WithCancel(ctx)
	closed := make(chan error)
	go func() {
		closed <- application.Worker.MustRun(ctx, *cfg)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	cancel()

	application.Worker.Shutdown(ctx, closed)
	client.Close()
	pool.Close()
	conn.Close()
}

func grpcClientConnection(ctx context.Context, cfg *config.Config) *grpc.ClientConn {
	c, err := grpcclient.NewClient(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return c
}

func connectToPostgres(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	pool, err := postgres.NewClient(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return pool
}

func connectToRedis(ctx context.Context, cfg *config.Config) *redis.Client {
	client, err := redisclient.NewClient(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return client
}
