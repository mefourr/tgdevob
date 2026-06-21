package worker

import (
	"context"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/controller/kafka_consumer"
	"log/slog"
	"sync"
)

type App struct {
	con *kafka_consumer.Consumer
}

func New(con *kafka_consumer.Consumer) *App {
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
	group, err := kafka_consumer.NewConsumerGroup(cfg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create kafka consumer group", "brokers", cfg.Kafka.BootstrapServers, "group", cfg.Kafka.Group, "err", err)
		return err
	}
	defer group.Close()
	slog.InfoContext(ctx, "kafka consumer group created", "brokers", cfg.Kafka.BootstrapServers, "group", cfg.Kafka.Group, "topics", cfg.Kafka.Topics)

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
	slog.InfoContext(ctx, "kafka consumer ready", "topics", cfg.Kafka.Topics, "group", cfg.Kafka.Group)

	select {
	case <-ctx.Done():
		slog.InfoContext(ctx, "consumer stopping: context cancelled")
	case err = <-errs:
		slog.ErrorContext(ctx, "consumer stopped with error", "err", err)
	}

	wg.Wait()
	return nil
}

func (a App) Shutdown(ctx context.Context, closed chan error) {
	<-closed
	slog.InfoContext(ctx, "kafka consumer shutdown complete")
}
