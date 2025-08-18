package config

import "time"

type StorageConfig struct {
	Username string
	Password string
	Hostname string
	Port     string
	Database string
	RetryNum int
	Delay    time.Duration
}
