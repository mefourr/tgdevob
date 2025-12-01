package worker

import (
	"context"
	"github.com/mefourr/tgdevob/worker/internal/infra/kafka"
)

type App struct {
	con *kafka.Consumer
}

func New(con *kafka.Consumer) *App {
	return &App{con: con}
}

func (a App) MustRun(ctx context.Context) {
	if err := a.run(ctx); err != nil {
		panic(err)
	}
}

func (a App) Shutdown(_ context.Context) {
}

func (a App) run(ctx context.Context) error {
	if err := a.con.Consume(ctx); err != nil {
		return err
	}
	return nil
}
