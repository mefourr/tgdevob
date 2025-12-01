package main

import (
	"context"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/infra/kafka"
	"github.com/mefourr/tgdevob/worker/internal/infra/repository/postgres"
	"github.com/mefourr/tgdevob/worker/internal/infra/repository/redis"
	ikafka "github.com/mefourr/tgdevob/worker/internal/interfaces/kafka"
	"github.com/mefourr/tgdevob/worker/internal/service"
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
	slog.InfoContext(ctx, "after setting redis up", "res", result)

	pool, err := postgres.NewClient(ctx, config.StorageConfig{
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
		ikafka.NewEventHandler(
			ikafka.NewEventProcessor(
				service.NewParserSvc(),
				service.NewLoadUserSvc(rediscache.New(rdb), postgres.New(pool)),
				service.NewSaveUserSvc(rediscache.New(rdb)),
				service.NewChecker(),
			),
		),
	)
	if err := consumer.Consume(ctx); err != nil {
		panic(err)
	}
}
