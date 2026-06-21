package consumer

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/tgbot/pkg/logger"
	"github.com/mefourr/tgdevob/tgbot/internal/kafka/consumer/startup"
	"golang.org/x/net/context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Consumer struct {
	bc      *startup.BaseConsumer
	Topic   string
	handler struct{} // handlers here
}

func New(bc *startup.BaseConsumer, topic string, handler struct{}) *Consumer {
	return &Consumer{
		bc:      bc,
		Topic:   topic,
		handler: handler,
	}
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.bc.Ready)
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				slog.InfoContext(session.Context(), "message channel closed, stopping consumer claim", "topic", claim.Topic(), "partition", claim.Partition())
				return nil
			}

			slog.DebugContext(session.Context(), "message claimed", "topic", message.Topic, "partition", message.Partition, "offset", message.Offset, "timestamp", message.Timestamp)

			//_ = c.controller.Process(context.Background(), message)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := newConsumerGroup(c.bc.Group, c.bc.Brokers)
	if err != nil {
		return fmt.Errorf("error creating consumer group client: %w", err)
	}
	slog.InfoContext(ctx, "consumer group created", "topic", c.Topic, "group", c.bc.Group)

	defer func() {
		if err := client.Close(); err != nil {
			slog.ErrorContext(logger.ErrorCtx(ctx, err), "error closing consumer group client")
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, []string{c.Topic}, c); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					slog.InfoContext(ctx, "consumer group closed", "topic", c.Topic)
					return
				}
				slog.ErrorContext(logger.ErrorCtx(ctx, err), "consumer group session error", "topic", c.Topic, "err", err)
			}

			if ctx.Err() != nil {
				slog.ErrorContext(logger.ErrorCtx(ctx, ctx.Err()), "consumer group closed by cancellation")
				return
			}

			c.bc.Ready = make(chan struct{})
		}
	}()

	// Wait until the consumer session is Ready
	<-c.bc.Ready
	slog.InfoContext(ctx, "kafka consumer ready", "topic", c.Topic, "group", c.bc.Group, "brokers", c.bc.Brokers)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigterm)

	select {
	case <-ctx.Done():
		slog.InfoContext(ctx, "consumer stopping: context cancelled", "reason", ctx.Err())
	case <-sigterm:
		slog.InfoContext(ctx, "consumer stopping: shutdown signal received")
	}

	cancel()
	wg.Wait()

	slog.InfoContext(ctx, "kafka consumer shutdown complete", "topic", c.Topic)
	return nil
}

func newConsumerGroup(group string, brokers []string) (sarama.ConsumerGroup, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	return sarama.NewConsumerGroup(brokers, group, cfg)
}
