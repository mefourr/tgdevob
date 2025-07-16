package handler

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type iHandler interface {
	HandleTgUserMessage(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type UserRequest struct {
	Request *tgbotapi.Message `json:"tg_request"`
}

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleTgUserMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	// TODO: idempotency guarantee
	// TODO: cache user
	// TODO: validate message
	// TODO: s3 grpc
	// TODO: recognition grpc
	var m UserRequest
	err := json.Unmarshal(msg.Value, &m)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "Successfully consume a msg", "from", m.Request.From.UserName)
	return nil
}
