package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/app"
	"github.com/mefourr/tgdevob/worker/internal/infra/repository/postgres"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/redis/go-redis/v9"
	"time"
)

// todo: main must load the app config, logger, initialize the worker app and do shutdown by signal
func main() {
	ctx := logging.Init()
	cfg := &config.Config{}

	rdb := connectToRedis(ctx, cfg)
	pool := connectToPostgres(ctx, cfg)
	defer pool.Close()

	application := app.New(cfg, rdb, pool)
	application.Worker.MustRun(ctx)

	//shutdown := make(chan os.Signal, 1)
	//signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	//
	//sig := <-shutdown
	//slog.InfoContext(ctx, "received shutdown signal", "signal", sig.String())
	//
	//application.Worker.Shutdown(ctx)
	//slog.InfoContext(ctx, "app has been shutdown")
}

func connectToPostgres(ctx context.Context, _ *config.Config) *pgxpool.Pool {
	pool, err := postgres.NewClient(ctx, config.Config{
		Username: "postgres",
		Password: "admin",
		Hostname: "localhost",
		Port:     "5432",
		Database: "postgres",
		RetryNum: 3,
		Delay:    time.Second,
	})
	if err != nil {
		panic(err)
	}
	return pool
}

func connectToRedis(ctx context.Context, _ *config.Config) *redis.Client {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	c.FlushDB(ctx)
	ping := c.Ping(ctx)
	if _, err := ping.Result(); err != nil {
		panic(err)
	}
	return c
}
