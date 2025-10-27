package consumer

import (
	"errors"
	"fmt"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/kafka/consumer/startup"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/IBM/sarama"
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
				slog.InfoContext(session.Context(), "worker channel was closed")
				return nil
			}

			slog.DebugContext(session.Context(), "Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

			//_ = c.handler.Process(context.Background(), message)
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
					slog.ErrorContext(logging.ErrorCtx(ctx, err), "consumer group closed by:", err)
					return
				}
				slog.ErrorContext(logging.ErrorCtx(ctx, err), "error from consumer")
			}

			if ctx.Err() != nil {
				slog.ErrorContext(logging.ErrorCtx(ctx, ctx.Err()), "consumer group closed by cancellation")
				return
			}

			c.bc.Ready = make(chan struct{})
		}
	}()

	// Wait until the consumer session is Ready
	<-c.bc.Ready
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
