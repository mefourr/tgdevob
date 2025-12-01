package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/app/worker"
	"github.com/mefourr/tgdevob/worker/internal/infra/kafka"
	"github.com/mefourr/tgdevob/worker/internal/infra/repository/postgres"
	rediscache "github.com/mefourr/tgdevob/worker/internal/infra/repository/redis"
	ikafka "github.com/mefourr/tgdevob/worker/internal/interfaces/kafka"
	"github.com/mefourr/tgdevob/worker/internal/service"
	"github.com/redis/go-redis/v9"
)

const (
	brokers = "localhost:9092" // TODO: retrieve from sys env
	topic   = "tg_requests"
	group   = "example"
)

type App struct {
	cfg    *config.Config
	client *redis.Client
	pool   *pgxpool.Pool
	Worker *worker.App
}

func New(cfg *config.Config, client *redis.Client, pool *pgxpool.Pool) *App {
	return &App{
		cfg:    cfg,
		client: client,
		pool:   pool,
		Worker: worker.New(createConsumer(client, pool)),
	}
}

func createConsumer(rdb *redis.Client, pool *pgxpool.Pool) *kafka.Consumer {
	return kafka.NewConsumer(
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
}
