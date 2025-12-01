package kafka

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/worker/internal/interfaces/kafka"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"golang.org/x/net/context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Consumer struct {
	ready   chan struct{}
	Brokers []string
	Topic   string
	Group   string
	handler *kafka.EventHandler
}

func NewConsumer(brokers []string, topic string, group string, handler *kafka.EventHandler) *Consumer {
	return &Consumer{
		ready:   make(chan struct{}),
		Brokers: brokers,
		Topic:   topic,
		Group:   group,
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

			slog.DebugContext(ctx, "consume message", "timestamp", message.Timestamp, "value", string(message.Value), "topic", c.Topic, "group", c.Group)

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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := newConsumerGroup(c.Group, c.Brokers)
	if err != nil {
		return fmt.Errorf("error creating consumer group client: %w", err)
	}
	slog.InfoContext(ctx, "consumer group created")

	defer func() {
		if err := client.Close(); err != nil {
			slog.ErrorContext(logging.ErrorCtx(ctx, err), "error closing consumer group client")
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, []string{c.Topic}, c); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					slog.ErrorContext(logging.ErrorCtx(ctx, err), "consumer group closed by", "err", err)
					return
				}
				slog.ErrorContext(logging.ErrorCtx(ctx, err), "error from consumer")
			}

			if ctx.Err() != nil {
				slog.ErrorContext(logging.ErrorCtx(ctx, ctx.Err()), "consumer group closed by cancellation")
				return
			}

			c.ready = make(chan struct{})
		}
	}()

	// Wait until the consumer session is ready
	<-c.ready
	slog.InfoContext(ctx, "Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigterm)

	select {
	case <-ctx.Done():
		slog.InfoContext(ctx, "terminating: context cancelled", "ctx", ctx.Err())
	case <-sigterm:
		slog.InfoContext(ctx, "terminating: via signal")
	}

	cancel() // Trigger goroutine shutdown
	wg.Wait()

	slog.InfoContext(ctx, "Sarama consumer shut down")
	return nil
}

func newConsumerGroup(group string, brokers []string) (sarama.ConsumerGroup, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	return sarama.NewConsumerGroup(brokers, group, cfg)
}
