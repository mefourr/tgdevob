package kafka_consumer

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"log/slog"
)

type Consumer struct {
	handler *EventHandler
	Ready   chan struct{}
}

func New(handler *EventHandler) *Consumer {
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
				slog.InfoContext(ctx, "message channel closed, stopping claim", "topic", claim.Topic(), "partition", claim.Partition())
				return nil
			}

			slog.DebugContext(ctx, "message claimed", "topic", message.Topic, "partition", message.Partition, "offset", message.Offset, "timestamp", message.Timestamp)

			if err := c.handler.Handle(ctx, message); err != nil {
				slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to process message", "topic", message.Topic, "offset", message.Offset, "err", err)
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
				slog.InfoContext(ctx, "consumer group closed", "topics", topics)
				return err
			}
			slog.ErrorContext(logger.ErrorCtx(ctx, err), "consumer group session error", "topics", topics, "err", err)
			return err
		}

		if ctx.Err() != nil {
			slog.InfoContext(ctx, "consumer stopping: context cancelled", "topics", topics)
			return ctx.Err()
		}

		c.Ready = make(chan struct{})
	}
}
