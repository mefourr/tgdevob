package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/adapter/grpc/validator"
	"github.com/mefourr/tgdevob/worker/internal/adapter/postgres"
	rediscache "github.com/mefourr/tgdevob/worker/internal/adapter/redis"
	"github.com/mefourr/tgdevob/worker/internal/app/worker"
	"github.com/mefourr/tgdevob/worker/internal/controller/kafka_consumer"
	processor_uc "github.com/mefourr/tgdevob/worker/internal/processor/usecase"
	user_uc "github.com/mefourr/tgdevob/worker/internal/user/usecase"
	validator_uc "github.com/mefourr/tgdevob/worker/internal/validator/usecase"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type App struct {
	Worker *worker.App
}

func New(
	cfg *config.Config,
	client *redis.Client,
	pool *pgxpool.Pool,
	conn *grpc.ClientConn,
) *App {
	return &App{
		Worker: worker.New(initConsumer(cfg, client, pool, conn)),
	}
}

func initConsumer(
	cfg *config.Config,
	client *redis.Client,
	pool *pgxpool.Pool,
	conn *grpc.ClientConn,
) *kafka_consumer.Consumer {
	var (
		grpcProducer = validator.New(conn)
		redisCache   = rediscache.New(client)
		postgresRepo = postgres.New(pool)
	)
	return kafka_consumer.New(kafka_consumer.NewEventHandler(
		processor_uc.New(
			user_uc.New(redisCache, postgresRepo),
			validator_uc.New(grpcProducer),
		),
	))
}
