package config

import "time"

type Config struct {
	Username string
	Password string
	Hostname string
	Port     string
	Database string
	RetryNum int
	Delay    time.Duration
}

//
//type Config struct {
//	App struct {
//		Name    string
//		Version string
//	} `mapstructure:"app"`
//	Grpc struct {
//		Port    int
//		Timeout time.Duration
//	} `mapstructure:"grpc"`
//}
//
//func MustLoadConfig() *Config {
//	return LoadConfig("config")
//}
//
//func LoadConfig(path string) *Config {
//	viper.SetConfigName("config")
//	viper.SetConfigType("yaml")
//	viper.AddConfigPath(path)
//
//	viper.AutomaticEnv()
//
//	if err := viper.ReadInConfig(); err != nil {
//		log.Fatalf("Error reading config file: %v", err)
//	}
//
//	var cfg Config
//	if err := viper.Unmarshal(&cfg); err != nil {
//		log.Fatalf("Error unmarshaling config: %v", err)
//	}
//
//	return &cfg
//}
