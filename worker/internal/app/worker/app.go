package worker

import (
	"context"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/infra/kafka"
)

type App struct {
	con *kafka.Consumer
}

func New(con *kafka.Consumer) *App {
	return &App{con: con}
}

// MustRun must be wrapped by canceled context
func (a App) MustRun(ctx context.Context, cfg config.Config) {
	if err := a.run(ctx, cfg); err != nil {
		panic(err)
	}
}

func (a App) run(ctx context.Context, cfg config.Config) error {
	group, err := kafka.NewConsumerGroup(cfg)
	if err != nil {
		return err
	}
	defer group.Close()

	// todo: reed topic from config as slice
	if err := a.con.Consume(ctx, group, []string{cfg.Kafka.Topic}); err != nil {
		return err
	}
	return nil
}
