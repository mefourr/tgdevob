package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/app"
	grpccon "github.com/mefourr/tgdevob/worker/internal/infra/grpc"
	"github.com/mefourr/tgdevob/worker/internal/infra/repository/postgres"
	rediscon "github.com/mefourr/tgdevob/worker/internal/infra/repository/redis"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoadConfig()
	ctx := logging.Init()

	client := connectToRedis(ctx, cfg)
	defer client.Close()
	slog.InfoContext(ctx, "Redis connection established")

	pool := connectToPostgres(ctx, cfg)
	defer pool.Close()
	slog.InfoContext(ctx, "Postgres connection established")

	conn := grpcClientConnection(ctx, cfg)
	defer conn.Close()
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
}

func grpcClientConnection(ctx context.Context, cfg *config.Config) *grpc.ClientConn {
	c, err := grpccon.NewClient(ctx, cfg)
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
	client, err := rediscon.NewClient(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return client
}
