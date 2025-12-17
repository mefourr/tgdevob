package worker

import (
	"context"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/infra/kafka"
	"log/slog"
	"sync"
)

type App struct {
	con *kafka.Consumer
}

func New(con *kafka.Consumer) *App {
	return &App{con: con}
}

// MustRun must be wrapped by canceled context
func (a App) MustRun(ctx context.Context, cfg config.Config) error {
	if err := a.run(ctx, cfg); err != nil {
		panic(err)
	}
	return nil
}

func (a App) run(ctx context.Context, cfg config.Config) error {
	group, err := kafka.NewConsumerGroup(cfg)
	if err != nil {
		return err
	}
	defer group.Close()

	wg := &sync.WaitGroup{}
	errs := make(chan error, 1)

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err = a.con.Consume(ctx, group, cfg.Kafka.Topics); err != nil {
			errs <- err
			close(errs)
		}
	}()

	<-a.con.Ready
	slog.InfoContext(ctx, "Sarama consumer up and running!...")

	select {
	case <-ctx.Done():
		slog.InfoContext(ctx, "kafka.Consume: context cancelled")
	case err = <-errs:
		slog.InfoContext(ctx, "kafka.Consume: sarama consumer error:", err)
	}

	wg.Wait()
	return nil
}

func (a App) Shutdown(ctx context.Context, closed chan error) {
	<-closed
	slog.InfoContext(ctx, "App:Shutdown Sarama consumer successfully closed")
}
