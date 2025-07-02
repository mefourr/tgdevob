package rqhandler

import (
	"context"
	"encoding/json"
	"gihub.com/mefourr/tgdevob/telegram-api/config"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/utils"
	"github.com/IBM/sarama"
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
		ctx = utils.WithLogUserName(ctx, update.Message.From.UserName)
		ctx = utils.WithLogUserID(ctx, update.Message.From.ID)
		ctx = utils.WithLogFileID(ctx, update.Message.Voice.FileID)
		h.handleVoiceMessage(ctx, update, bot)
		return
	}

	if update.Message.Text != "" {
		ctx = utils.WithLogUserName(ctx, update.Message.From.UserName)
		h.handleTextMessage(ctx, update, bot)
		return
	}

	slog.InfoContext(ctx, "Unrecognized message type", "update", update)
}

// stub
func (h *Handler) handleVoiceMessage(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	slog.InfoContext(ctx, "Received voice message", "from", update.Message.From.UserName)
	type message struct {
		Request *tgbotapi.Message `json:"tg_request"`
	}
	mBytes, err := json.Marshal(message{
		Request: update.Message,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal voice message", "error", err)
	}

	err = pushRequestToQueue("tg_requests", mBytes)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to push request to queue", "error", err)
	}
}

func connectProducer(brokers []string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, cfg)
}

func pushRequestToQueue(topic string, message []byte) error {
	// host of kafka container here
	brokers := []string{"localhost:9092"}
	producer, err := connectProducer(brokers)
	if err != nil {
		return err
	}
	defer producer.Close()

	prodMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(prodMessage)
	if err != nil {
		return err
	}
	slog.InfoContext(context.TODO(), "Successfully sent message", "partition", partition, "offset", offset)
	return nil
}

func (h *Handler) handleTextMessage(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	slog.InfoContext(ctx, "Received message", "from", update.Message.From.UserName, "message", update.Message.Text)

	stickerID := "CAACAgIAAxkBAAEPVdJoMI0T572X6QjdE0rIKPdp-uQgWwACYQADUomRI5wSPlG4RvGWNgQ"
	sticker := tgbotapi.NewSticker(update.Message.Chat.ID, tgbotapi.FileID(stickerID))

	if _, err := bot.Send(sticker); err != nil {
		slog.ErrorContext(utils.ErrorCtx(ctx, err), "Error sending sticker")
	}
}
