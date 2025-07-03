package rqhandler

import (
	"context"
	"gihub.com/mefourr/tgdevob/telegram-api/config"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/broker"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		slog.InfoContext(ctx, "Update does not contain a message")
		return
	}

	if update.Message.Voice != nil {
		slog.InfoContext(ctx, "Update has a voice message")
		ctx = utils.WithLogUserName(ctx, update.Message.From.UserName)
		ctx = utils.WithLogUserID(ctx, update.Message.From.ID)
		ctx = utils.WithLogFileID(ctx, update.Message.Voice.FileID)
		h.handleVoiceMessage(ctx, update)
		return
	}

	if update.Message.Text != "" {
		slog.InfoContext(ctx, "Update has a text message")
		ctx = utils.WithLogUserName(ctx, update.Message.From.UserName)
		h.handleTextMessage(ctx, update, bot)
		return
	}

	slog.InfoContext(ctx, "Unrecognized message type", "update", update)
}

func (h *Handler) handleVoiceMessage(ctx context.Context, update tgbotapi.Update) {
	producer := broker.NewProducer("tg_requests", []string{"localhost:9092"})
	slog.DebugContext(ctx, "Ready to produce a msg")
	producer.ProduceMessage(ctx, update)
}

func (h *Handler) handleTextMessage(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	slog.InfoContext(ctx, "Received message", "from", update.Message.From.UserName, "message", update.Message.Text)

	stickerID := "CAACAgIAAxkBAAEPVdJoMI0T572X6QjdE0rIKPdp-uQgWwACYQADUomRI5wSPlG4RvGWNgQ"
	sticker := tgbotapi.NewSticker(update.Message.Chat.ID, tgbotapi.FileID(stickerID))

	if _, err := bot.Send(sticker); err != nil {
		slog.ErrorContext(utils.ErrorCtx(ctx, err), "Error sending sticker")
	}
}
