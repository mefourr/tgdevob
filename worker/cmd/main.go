package main

import (
	"context"
	"gihub.com/mefourr/tgdevob/worker/internal/handler/worker"
	"gihub.com/mefourr/tgdevob/worker/internal/kafka"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/redis/go-redis/v9"
	"log/slog"
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

	consumer := kafka.NewConsumer(
		[]string{brokers},
		topic,
		group,
		worker.New(rdb),
	)
	if err := consumer.Consume(ctx); err != nil {
		panic(err)
	}
}
