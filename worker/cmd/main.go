package main

import (
	"context"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/handler/loaduser"
	"github.com/mefourr/tgdevob/worker/internal/handler/worker"
	"github.com/mefourr/tgdevob/worker/internal/kafka"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/cache"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/postgresql"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

const (
	brokers = "localhost:9092" // TODO: retrieve from sys env
	topic   = "tg_requests"
	group   = "example"
)

func main() {
	ctx := logging.Init()
	slog.InfoContext(ctx, "Logger for consumer is initialized")
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_ = rdb.FlushDB(ctx).Err()

	ping := rdb.Ping(context.Background())
	result, err := ping.Result()
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
	}
	slog.InfoContext(ctx, "after setting redis up result is ", result)

	pool, err := postgresql.NewClient(ctx, config.StorageConfig{
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
	defer pool.Close()

	consumer := kafka.NewConsumer(
		[]string{brokers},
		topic,
		group,
		worker.New(loaduser.New(
			cache.New(rdb),
			postgresql.New(pool),
		)),
	)
	if err := consumer.Consume(ctx); err != nil {
		panic(err)
	}
}
