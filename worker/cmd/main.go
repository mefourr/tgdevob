package main

import (
	"gihub.com/mefourr/tgdevob/worker/internal/handler"
	"gihub.com/mefourr/tgdevob/worker/internal/kafka"
	"gihub.com/mefourr/tgdevob/worker/internal/utils"
	"log/slog"
)

const (
	brokers = "localhost:9092" // TODO: retrieve from sys env
	topic   = "tg_requests"
	group   = "example"
)

func main() {
	ctx := utils.Init()
	slog.InfoContext(ctx, "Logger for consumer is initialized")
	consumer := kafka.NewConsumer(
		[]string{brokers},
		topic,
		group,
		handler.NewHandler(),
	)
	if err := consumer.Consume(ctx); err != nil {
		panic(err)
	}
}
