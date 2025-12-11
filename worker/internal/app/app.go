package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/app/worker"
	"github.com/mefourr/tgdevob/worker/internal/infra/grpc/vmlength"
	"github.com/mefourr/tgdevob/worker/internal/infra/kafka"
	"github.com/mefourr/tgdevob/worker/internal/infra/repository/postgres"
	rediscache "github.com/mefourr/tgdevob/worker/internal/infra/repository/redis"
	ikafka "github.com/mefourr/tgdevob/worker/internal/interfaces/kafka"
	"github.com/mefourr/tgdevob/worker/internal/service"
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
) *kafka.Consumer {
	var (
		grpcProducer = vmlength.New(conn)
		cache        = rediscache.New(client)
		repo         = postgres.New(pool)
	)
	return kafka.NewConsumer(
		cfg,
		ikafka.NewEventHandler(
			ikafka.NewEventProcessor(
				service.NewParserSvc(),
				service.NewLoadUserSvc(cache, repo),
				service.NewSaveUserSvc(cache),
				service.NewChecker(),
				service.NewVoiceDurationValidator(grpcProducer),
			),
		),
	)
}
