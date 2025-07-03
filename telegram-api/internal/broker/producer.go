package broker

import (
	"context"
	"encoding/json"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/message"
	"github.com/IBM/sarama"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type Producer struct {
	Topic   string
	Brokers []string
}

func NewProducer(topic string, brokers []string) *Producer {
	return &Producer{Topic: topic, Brokers: brokers}
}

// TODO: come up with smt better
func (p *Producer) ProduceMessage(ctx context.Context, rowMsg any) {
	switch msg := rowMsg.(type) {
	case tgbotapi.Update:
		slog.DebugContext(ctx, "Message type is tgbotapi.Update")
		userrq := message.NewUserRequest(msg.Message)
		p.produce(ctx, userrq)
	case string:
		rs := message.NewResponse(msg)
		p.produce(ctx, rs)
	}
}

func (p *Producer) produce(ctx context.Context, message any) {
	bytes, err := json.Marshal(message)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal voice message", "error", err)
	}

	err = p.pushRequestToQueue(ctx, bytes)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to push request to queue", "error", err)
	}
}

func (p *Producer) pushRequestToQueue(ctx context.Context, message []byte) error {
	client, err := createProducer(p.Brokers)
	if err != nil {
		return err
	}
	defer client.Close()

	msg := &sarama.ProducerMessage{
		Topic: p.Topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := client.SendMessage(msg)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Successfully sent message", "partition", partition, "offset", offset)
	return nil
}

func createProducer(brokers []string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, cfg)
}
