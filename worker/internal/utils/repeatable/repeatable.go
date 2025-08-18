package repeatable

import (
	"context"
	"gihub.com/mefourr/tgdevob/worker/config"
	"log/slog"
	"time"
)

func Connect(fn func() error, conf config.StorageConfig) (err error) {
	for conf.RetryNum > 0 {
		if err = fn(); err != nil {
			time.Sleep(conf.Delay)
			conf.RetryNum--
			slog.WarnContext(context.Background(), "try to connect to psql by", "retry", conf.RetryNum)
			continue
		}
		return nil
	}
	return
}
