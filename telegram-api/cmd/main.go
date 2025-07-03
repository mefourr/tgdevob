package main

import (
	"context"
	"fmt"
	"gihub.com/mefourr/tgdevob/telegram-api/config"
	rqhandler "gihub.com/mefourr/tgdevob/telegram-api/internal/bot"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/broker"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

func main() {
	ctx := utils.Init()
	slog.InfoContext(ctx, "Logger initialized")

	cfg := config.LoadConfig(ctx)
	slog.InfoContext(ctx, "Config loaded")

	go func() {
		consumer := broker.NewConsumer(
			[]string{"localhost:9092"},
			"tg_requests",
			"example",
		)
		if err := consumer.Consume(ctx); err != nil {
			panic(err)
		}
	}()

	//	_ = repository.NewDatabase(cfg.Database)
	//	slog.InfoContext(ctx, "Database initialized")

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
	handler := rqhandler.NewHandler(cfg)
	slog.InfoContext(ctx, "Bot are listening")

	for update := range updates {
		slog.DebugContext(ctx, fmt.Sprintf("Update: %+v", update))
		handler.ProcessUpdate(ctx, update, bot)
	}
}
