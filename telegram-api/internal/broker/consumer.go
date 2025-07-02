package broker

import (
	"encoding/json"
	"github.com/IBM/sarama"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/context"
	"log/slog"
)

type Consumer struct {
	Ready chan bool
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				slog.InfoContext(context.TODO(), "message channel was closed")
				return nil
			}
			slog.InfoContext(context.TODO(), "Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			_ = todo(message)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
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
