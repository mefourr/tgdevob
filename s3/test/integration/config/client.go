package config

import "time"

type Client struct {
	Host    string
	Port    int
	Timeout time.Duration
}
