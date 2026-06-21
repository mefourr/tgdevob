package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mefourr/tgdevob/tgbot/config"
	rqhandler "github.com/mefourr/tgdevob/tgbot/internal/bot"
	"github.com/mefourr/tgdevob/tgbot/internal/kafka/consumer"
	"github.com/mefourr/tgdevob/tgbot/internal/kafka/consumer/startup"
	"github.com/mefourr/tgdevob/tgbot/pkg/logger"
	"log/slog"
)

func main() {
	ctx := logger.Init()
	slog.InfoContext(ctx, "Logger initialized")

	// TODO: come up with smt better with config impl
	cfg := config.LoadConfig(ctx)
	slog.InfoContext(ctx, "Config loaded")

	bc := &startup.BaseConsumer{
		Ready:   make(chan struct{}),
		Brokers: []string{"localhost:9092"}, // TODO: must be replaced
		Group:   "tg-cons",
	}

	go func() {
		if err := consumer.New(bc, "invalidated_user_messages", struct{}{}).
			Consume(ctx); err != nil {
			slog.ErrorContext(ctx, "kafka consumer exited with error", "err", err)
			panic(err)
		}
	}()

	RunBot(ctx, *cfg)
}

func RunBot(ctx context.Context, cfg config.Config) {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create telegram bot", "err", err)
		panic(err)
	}
	slog.InfoContext(ctx, "telegram bot authenticated", "username", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Telegram.Timeout

	updates := bot.GetUpdatesChan(u)
	h := rqhandler.NewHandler(cfg)
	slog.InfoContext(ctx, "bot is listening for updates", "timeout", cfg.Telegram.Timeout)

	for update := range updates {
		slog.DebugContext(ctx, "update received", "update_id", update.UpdateID, "chat_id", update.Message.Chat.ID)
		h.ProcessUpdate(ctx, update, bot)
	}
}
