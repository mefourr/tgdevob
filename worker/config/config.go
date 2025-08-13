package config

import "time"

type PostgresConfig struct {
	Username string
	Password string
	Hostname string
	Port     string
	Database string
	RetryNum int
	Delay    time.Duration
}
