package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/internal/interfaces/kafka"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"log/slog"
	"sync"
)

type Consumer struct {
	cfg     *config.Config
	ready   chan struct{}
	handler *kafka.EventHandler
}

func NewConsumer(cfg *config.Config, handler *kafka.EventHandler) *Consumer {
	return &Consumer{
		cfg:     cfg,
		ready:   make(chan struct{}),
		handler: handler,
	}
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				slog.InfoContext(ctx, "worker channel was closed")
				return nil
			}

			slog.DebugContext(ctx, "consume message", "timestamp", message.Timestamp, "value", string(message.Value), "topic", c.cfg.Kafka.Topic, "group", c.cfg.Kafka.Group)

			if err := c.handler.Handle(ctx, message); err != nil {
				slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to process message", "error", err)
			}

			session.MarkMessage(message, "")
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	client, err := newConsumerGroup(c.cfg.Kafka.Group, c.cfg.Kafka.Bootstraps)
	if err != nil {
		return err
	}
	defer client.Close()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, []string{c.cfg.Kafka.Topic}, c); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					slog.ErrorContext(logging.ErrorCtx(ctx, err), "consumer group closed by", "err", err)
					return
				}
				slog.ErrorContext(logging.ErrorCtx(ctx, err), "error from consumer")
			}

			if ctx.Err() != nil {
				slog.InfoContext(ctx, "consumer group closed by cancellation")
				return
			}

			c.ready = make(chan struct{})
		}
	}()

	// Wait until the consumer session is ready
	<-c.ready
	slog.InfoContext(ctx, "Sarama consumer up and running!...")

	select {
	case <-ctx.Done():
		slog.InfoContext(ctx, "kafka.Consume: context cancelled")
	}

	wg.Wait()

	slog.InfoContext(ctx, "Sarama consumer successfully closed")
	return nil
}

func newConsumerGroup(group string, brokers []string) (sarama.ConsumerGroup, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	return sarama.NewConsumerGroup(brokers, group, cfg)
}
