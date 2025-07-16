package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type Producer struct {
	Topic   string
	Brokers []string
}

type UserRequest struct {
	Request *tgbotapi.Message `json:"tg_request"`
}

func NewProducer(topic string, brokers []string) *Producer {
	return &Producer{Topic: topic, Brokers: brokers}
}

// ProduceVoiceMessage TODO: come up with smt better
func (p *Producer) ProduceVoiceMessage(ctx context.Context, update tgbotapi.Update) {
	userrq := UserRequest{
		Request: update.Message,
	}
	p.produce(ctx, userrq)
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

// TODO: do once
func createProducer(brokers []string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, cfg)
}
