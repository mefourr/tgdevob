package config

type Grpc struct {
	Auth        `mapstructure:"auth"`
	YandexCloud `mapstructure:"yandex-cloud"`
	Server      `mapstructure:"server"`
}
