package rqhandler

import (
	"context"
	"errors"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/broker"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/utils"
	"github.com/IBM/sarama"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Consumer struct {
	topic string
	ctx   context.Context
}

func NewConsumer(ctx context.Context, topic string) *Consumer {
	return &Consumer{
		topic: topic,
		ctx:   ctx,
	}
}

func createConsumerGroup(brokers []string) (sarama.ConsumerGroup, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	return sarama.NewConsumerGroup(brokers, "example", cfg)
}

func (c *Consumer) Listen() error {
	keepRunning := true

	ctx, cancel := context.WithCancel(c.ctx)
	client, err := createConsumerGroup([]string{"localhost:9092"})
	if err != nil {
		cancel()
		return errors.New("error creating consumer group client")
	}

	consumer := broker.Consumer{
		Ready: make(chan bool),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, []string{"tg_requests"}, &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				slog.ErrorContext(utils.ErrorCtx(ctx, err), "Error from consumer")
			}
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	slog.InfoContext(ctx, "Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			slog.InfoContext(ctx, "terminating: via signal")
			keepRunning = false
		}
	}

	cancel()
	wg.Wait()
	slog.InfoContext(ctx, "Sarama consumer shutting down")
	if err = client.Close(); err != nil {
		slog.ErrorContext(utils.ErrorCtx(ctx, err), "Error closing client")
	}

	return nil
}
