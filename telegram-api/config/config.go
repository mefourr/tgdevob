package config

import (
	"context"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"os"
)

type Config struct {
	App struct {
		Name    string
		Version string
	} `mapstructure:"app"`
	Telegram struct {
		Token      string
		BaseUrl    string
		WebhookUrl string
		Timeout    int
	} `mapstructure:"telegram"`
	SpeechKit struct {
		Url     string
		Token   string
		Timeout int
	} `mapstructure:"speechKit"`
	Database struct {
		Host           string
		Port           string
		User           string
		Password       string
		Name           string
		MaxConnections int32
	} `mapstructure:"database"`
	Logging struct {
		level string
		file  string
	} `mapstructure:"logging"`
}

func LoadConfig(ctx context.Context) *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./telegram-api/config")

	viper.AutomaticEnv()

	viper.SetEnvPrefix("")
	viper.AllowEmptyEnv(true)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	slog.InfoContext(ctx, "File loaded successfully")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}
	slog.InfoContext(ctx, "File unmarshalled successfully")

	if cfg.Telegram.Token = os.Getenv(cfg.Telegram.Token); cfg.SpeechKit.Token == "" {
		log.Fatalf("SpeechKit token is empty")
	}
	if cfg.SpeechKit.Token = os.Getenv(cfg.SpeechKit.Token); cfg.Telegram.Token == "" {
		log.Fatalf("Telegram token is empty")
	}

	return &cfg
}
