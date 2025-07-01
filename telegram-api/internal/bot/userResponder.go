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
	consumer, err := co.ConsumePartition(c.topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}
	slog.InfoContext(c.ctx, "Consumer started")

	count := 0

	done := make(chan struct{})
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				slog.ErrorContext(utils.ErrorCtx(c.ctx, err), "error occurred")
			case _ = <-consumer.Messages():
				count++
				slog.InfoContext(c.ctx, "Consumer received message", "msg_count", count)
				//err = todo(msg)
				//if err != nil {
				//	slog.ErrorContext(utils.ErrorCtx(context.TODO(), err), "Failed to unmarshall consumed msg")
				//}
			case <-sigChannel:
				slog.InfoContext(c.ctx, "Consumer shutting down")
				done <- struct{}{}
			}
		}
	}()

	<-done
	slog.InfoContext(c.ctx, "Consumer shutting down")
	if err = consumer.Close(); err != nil {
		slog.ErrorContext(utils.ErrorCtx(c.ctx, err), "error occurred while closing consumer")
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
