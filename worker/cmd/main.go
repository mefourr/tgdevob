package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/app"
	"github.com/mefourr/tgdevob/worker/internal/infra/repository/postgres"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
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

	conn := grpcClientConnection(cfg)
	defer conn.Close()
	slog.InfoContext(ctx, "GRPC server connection established")

	application := app.New(cfg, client, pool, conn)

	ctx, cancel := context.WithCancel(ctx)
	go application.Worker.MustRun(ctx)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	cancel()
}

func grpcClientConnection(cfg *config.Config) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		net.JoinHostPort(cfg.Grpc.Host, cfg.Grpc.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Grpc server connection failed")
	}
	return conn
}

func connectToPostgres(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	pool, err := postgres.NewClient(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return pool
}

func connectToRedis(ctx context.Context, cfg *config.Config) *redis.Client {
	// todo: wrap conn creation
	c := redis.NewClient(&redis.Options{
		//Addr:     "localhost:6379",
		Addr:               net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port),
		Password:           cfg.Redis.Password, // no password set
		Username:           cfg.Redis.Username,
		DB:                 cfg.Redis.Database, // use default DB
		DialerRetries:      cfg.Redis.DialerRetries,
		DialerRetryTimeout: cfg.Redis.DialerRetryTimeout,
	})
	c.FlushDB(ctx)
	ping := c.Ping(ctx)
	if _, err := ping.Result(); err != nil {
		panic(err)
	}
	return c
}
