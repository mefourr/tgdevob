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
	UpdateId int               `json:"update_id"`
	User     *tgbotapi.User    `json:"sent_from"`
	Request  *tgbotapi.Message `json:"worker"`
}

func NewProducer(topic string, brokers []string) *Producer {
	return &Producer{Topic: topic, Brokers: brokers}
}

func (p *Producer) ProduceVoiceMessage(ctx context.Context, update tgbotapi.Update) {
	userrq := UserRequest{
		UpdateId: update.UpdateID,
		User:     update.SentFrom(),
		Request:  update.Message,
	}
	slog.DebugContext(ctx, "preparing voice message event", "update_id", update.UpdateID, "user_id", update.SentFrom().ID, "topic", p.Topic)
	p.produce(ctx, userrq)
}

func (p *Producer) produce(ctx context.Context, message any) {
	bytes, err := json.Marshal(message)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal message payload", "err", err)
		return
	}
	slog.DebugContext(ctx, "message payload serialized", "bytes", len(bytes))

	if err = p.pushRequestToQueue(ctx, bytes); err != nil {
		slog.ErrorContext(ctx, "failed to push message to kafka", "topic", p.Topic, "err", err)
	}
}

func (p *Producer) pushRequestToQueue(ctx context.Context, message []byte) error {
	client, err := createProducer(p.Brokers)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create kafka producer", "brokers", p.Brokers, "err", err)
		return err
	}
	defer client.Close()

	msg := &sarama.ProducerMessage{
		Topic: p.Topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := client.SendMessage(msg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send message to kafka", "topic", p.Topic, "err", err)
		return err
	}

	slog.InfoContext(ctx, "voice message event published", "topic", p.Topic, "partition", partition, "offset", offset)
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
