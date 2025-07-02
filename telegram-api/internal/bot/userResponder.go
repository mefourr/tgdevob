package rqhandler

import (
	"context"
	"encoding/json"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/utils"
	"github.com/IBM/sarama"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
	"os/signal"
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

func connectConsumer(brokers []string) (sarama.Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true

	return sarama.NewConsumer(brokers, cfg)
}

func (c *Consumer) Listen() error {
	co, err := connectConsumer([]string{"localhost:9092"})
	if err != nil {
		return err
	}
	partitionConsumer, err := co.ConsumePartition(c.topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}
	slog.InfoContext(c.ctx, "Consumer started")

	count := 0

	done := make(chan struct{})
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer close(done)
		for {
			select {
			case err := <-partitionConsumer.Errors():
				slog.ErrorContext(utils.ErrorCtx(c.ctx, err), "error occurred")
				return

			case _, ok := <-partitionConsumer.Messages():
				if !ok {
					slog.ErrorContext(utils.ErrorCtx(c.ctx, err), "error occurred while claiming a message")
					return
				}
				count++
				slog.InfoContext(c.ctx, "Consumer received message", "msg_count", count)

			case <-sigChannel:
				slog.InfoContext(c.ctx, "Consumer shutting down by signal Ctrl+c")
				return
			}
		}
	}()

	<-done
	slog.InfoContext(c.ctx, "Consumer shutting down")
	if err = partitionConsumer.Close(); err != nil {
		slog.ErrorContext(utils.ErrorCtx(c.ctx, err), "error occurred while closing partitionConsumer")
	}

	return nil
}

func todo(msg *sarama.ConsumerMessage) error {
	type message struct {
		Request *tgbotapi.Message `json:"tg_request"`
	}
	m := message{}
	err := json.Unmarshal(msg.Value, &m)
	if err != nil {
		return err
	}
	slog.InfoContext(context.TODO(), "Successfully consume a msg", "from", m.Request.From.UserName)
	return nil
}
