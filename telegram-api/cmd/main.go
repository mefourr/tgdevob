package main

import (
	"context"
	"fmt"
	"gihub.com/mefourr/tgdevob/telegram-api/config"
	rqhandler "gihub.com/mefourr/tgdevob/telegram-api/internal/bot"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/kafka/consumer"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/kafka/consumer/startup"
	"gihub.com/mefourr/tgdevob/telegram-api/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

func main() {
	ctx := logging.Init()
	slog.InfoContext(ctx, "Logger initialized")

	// TODO: come up with smt better with config impl
	cfg := config.LoadConfig(ctx)
	slog.InfoContext(ctx, "Config loaded")

	bc := &startup.BaseConsumer{
		Ready:   make(chan struct{}),
		Brokers: []string{"localhost:9092"}, // TODO: must be replaced
		Group:   "tg-cons",
	}
	if err := consumer.New(bc, "invalidated_user_messages", struct{}{}).
		Consume(ctx); err != nil {
		panic(err)
	}
	if err := consumer.New(bc, "processed_user_messages", struct{}{}).
		Consume(ctx); err != nil {
		panic(err)
	}

	RunBot(ctx, *cfg)
}

func RunBot(ctx context.Context, cfg config.Config) {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.Telegram.Timeout

	updates := bot.GetUpdatesChan(u)
	h := rqhandler.NewHandler(cfg)
	slog.InfoContext(ctx, "Bot are listening")

	for update := range updates {
		slog.DebugContext(ctx, fmt.Sprintf("New Update: %+v", update))
		h.ProcessUpdate(ctx, update, bot)
	}
}
