package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	App  `mapstructure:"app"`
	Grpc `mapstructure:"grpc"`
}

func MustLoadConfig() *Config {
	return LoadConfig("config")
}

func LoadConfig(path string) *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("s3/" + path)

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
