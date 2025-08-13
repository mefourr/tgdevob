package utils

import (
	"gihub.com/mefourr/tgdevob/worker/config"
	"time"
)

func TryConn(fn func() error, conf config.PostgresConfig) (err error) {
	for conf.RetryNum > 0 {
		if err = fn(); err != nil {
			time.Sleep(conf.Delay)
			conf.RetryNum--
			continue
		}
		return nil
	}
	return
}
