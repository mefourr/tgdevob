package broker

import (
	"errors"
	"fmt"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/bot"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/utils"
	"github.com/IBM/sarama"
	"golang.org/x/net/context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Consumer struct {
	ready   chan bool
	Brokers []string
	Topic   string
	Group   string
}

func NewConsumer(brokers []string, topic, group string) *Consumer {
	return &Consumer{
		ready:   make(chan bool),
		Brokers: brokers,
		Topic:   topic,
		Group:   group,
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
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				slog.InfoContext(session.Context(), "message channel was closed")
				return nil
			}
			slog.DebugContext(session.Context(), "Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

			rqhandler.Todo(session.Context(), message)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
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
	defer func() {
		if err := client.Close(); err != nil {
			slog.ErrorContext(utils.ErrorCtx(ctx, err), "error closing consumer group client")
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, []string{c.Topic}, c); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					slog.ErrorContext(utils.ErrorCtx(ctx, err), "consumer group closed by:", err)
					return
				}
				slog.ErrorContext(utils.ErrorCtx(ctx, err), "error from consumer")
			}
			if ctx.Err() != nil {
				slog.ErrorContext(utils.ErrorCtx(ctx, ctx.Err()), "consumer group closed by cancellation or deadline")
				return
			}

			c.ready = make(chan bool)
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
