package config

import "time"

type Server struct {
	Host    string
	Port    int
	Timeout time.Duration
}
