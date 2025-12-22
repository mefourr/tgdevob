package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/worker/internal/interfaces/kafka"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"log/slog"
)

type Consumer struct {
	handler *kafka.EventHandler
	Ready   chan struct{}
}

func NewConsumer(handler *kafka.EventHandler) *Consumer {
	return &Consumer{
		handler: handler,
		Ready:   make(chan struct{}),
	}
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.Ready)
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

			slog.DebugContext(ctx, "consume message", "timestamp", message.Timestamp, "value", string(message.Value), "topic", message.Topic, "key", message.Key)

			if err := c.handler.Handle(ctx, message); err != nil {
				slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to process message", "error", err)
			}

			session.MarkMessage(message, "")
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *Consumer) Consume(ctx context.Context, group sarama.ConsumerGroup, topics []string) error {
	for {
		if err := group.Consume(ctx, topics, c); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				slog.ErrorContext(logging.ErrorCtx(ctx, err), "consumer group closed by", "err", err)
				return err
			}
			slog.ErrorContext(logging.ErrorCtx(ctx, err), "error from consumer")
		}

		if ctx.Err() != nil {
			slog.InfoContext(ctx, "consumer group closed by cancellation")
			return ctx.Err()
		}

		c.Ready = make(chan struct{})
	}
}
