package rqhandler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mefourr/tgdevob/tgbot/config"
	"github.com/mefourr/tgdevob/tgbot/internal/kafka"
	"github.com/mefourr/tgdevob/tgbot/pkg/logger"
	"log/slog"
)

type Handler struct {
	cfg config.Config
}

func NewHandler(cfg config.Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) ProcessUpdate(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message == nil {
		slog.InfoContext(ctx, "Update does not contain a worker")
		return
	}

	if update.Message.Voice != nil {
		slog.InfoContext(ctx, "Update has a voice worker")
		ctx = logger.WithLogUserName(ctx, update.Message.From.UserName)
		ctx = logger.WithLogUserID(ctx, update.Message.From.ID)
		ctx = logger.WithLogFileID(ctx, update.Message.Voice.FileID)
		h.handleVoiceMessage(ctx, update)
		return
	}

	if update.Message.Text != "" {
		slog.InfoContext(ctx, "Update has a text worker")
		ctx = logger.WithLogUserName(ctx, update.Message.From.UserName)
		h.handleTextMessage(ctx, update, bot)
		return
	}

	slog.InfoContext(ctx, "Unrecognized worker type", "update", update)
}

func (h *Handler) handleVoiceMessage(ctx context.Context, update tgbotapi.Update) {
	producer := kafka.NewProducer("tg_requests", []string{"localhost:9092"})
	slog.DebugContext(ctx, "Ready to produce a msg")
	producer.ProduceVoiceMessage(ctx, update)
}

// stub
func (h *Handler) handleTextMessage(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	slog.InfoContext(ctx, "Received worker", "from", update.Message.From.UserName, "worker", update.Message.Text)

	stickerID := "CAACAgIAAxkBAAEPVdJoMI0T572X6QjdE0rIKPdp-uQgWwACYQADUomRI5wSPlG4RvGWNgQ"
	sticker := tgbotapi.NewSticker(update.Message.Chat.ID, tgbotapi.FileID(stickerID))

	if _, err := bot.Send(sticker); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Error sending sticker")
	}
}
