package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	App struct {
		Name    string
		Version string
	} `mapstructure:"app"`
	Kafka struct {
		BootstrapServers []string
		Group            string
		Topics           []string
	} `mapstructure:"kafka"`
	Postgres struct {
		Database string
		Username string
		Password string
		Host     string
		Port     string
		RetryNum int
		Delay    time.Duration
	} `mapstructure:"postgres"`
	Redis struct {
		Database           int
		Username           string
		Password           string
		Host               string
		Port               string
		DialerRetries      int
		DialerRetryTimeout time.Duration
	} `mapstructure:"redis"`
	Grpc struct {
		Port    string
		Host    string
		Timeout time.Duration
	} `mapstructure:"grpc"`
}

func MustLoadConfig() *Config {
	return LoadConfig("config")
}

func LoadConfig(path string) *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	return &cfg
}
