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
		slog.DebugContext(ctx, "update contains no message, skipping", "update_id", update.UpdateID)
		return
	}

	if update.Message.Voice != nil {
		ctx = logger.WithLogUserName(ctx, update.Message.From.UserName)
		ctx = logger.WithLogUserID(ctx, update.Message.From.ID)
		ctx = logger.WithLogFileID(ctx, update.Message.Voice.FileID)
		slog.InfoContext(ctx, "voice message received", "duration_sec", update.Message.Voice.Duration, "file_size", update.Message.Voice.FileSize)
		h.handleVoiceMessage(ctx, update)
		return
	}

	if update.Message.Text != "" {
		ctx = logger.WithLogUserName(ctx, update.Message.From.UserName)
		ctx = logger.WithLogUserID(ctx, update.Message.From.ID)
		slog.InfoContext(ctx, "text message received", "chat_id", update.Message.Chat.ID)
		h.handleTextMessage(ctx, update, bot)
		return
	}

	slog.WarnContext(ctx, "unrecognized message type, skipping", "update_id", update.UpdateID, "chat_id", update.Message.Chat.ID)
}

func (h *Handler) handleVoiceMessage(ctx context.Context, update tgbotapi.Update) {
	producer := kafka.NewProducer("tg_requests", []string{"localhost:9092"})
	slog.DebugContext(ctx, "publishing voice message event to kafka", "topic", "tg_requests", "update_id", update.UpdateID)
	producer.ProduceVoiceMessage(ctx, update)
}

func (h *Handler) handleTextMessage(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	slog.DebugContext(ctx, "sending sticker response", "chat_id", update.Message.Chat.ID)

	stickerID := "CAACAgIAAxkBAAEPVdJoMI0T572X6QjdE0rIKPdp-uQgWwACYQADUomRI5wSPlG4RvGWNgQ"
	sticker := tgbotapi.NewSticker(update.Message.Chat.ID, tgbotapi.FileID(stickerID))

	if _, err := bot.Send(sticker); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to send sticker", "chat_id", update.Message.Chat.ID, "err", err)
		return
	}
	slog.InfoContext(ctx, "sticker sent successfully", "chat_id", update.Message.Chat.ID)
}
