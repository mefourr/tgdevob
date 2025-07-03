package rqhandler

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

func Todo(ctx context.Context, msg *sarama.ConsumerMessage) {
	type message struct {
		Request *tgbotapi.Message `json:"tg_request"`
	}
	m := message{}
	err := json.Unmarshal(msg.Value, &m)
	if err != nil {
		return
	}
	slog.InfoContext(ctx, "Successfully consume a msg", "from", m.Request.From.UserName)
}
