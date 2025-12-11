package repeatable

import (
	"context"
	"github.com/mefourr/tgdevob/worker/config"
	"log/slog"
	"time"
)

func Connect(fn func() error, conf *config.Config) (err error) {
	for conf.Postgres.RetryNum > 0 {
		if err = fn(); err != nil {
			time.Sleep(conf.Postgres.Delay)
			conf.Postgres.RetryNum--
			slog.WarnContext(context.Background(), "try to connect to psql by", "retry", conf.Postgres.RetryNum)
			continue
		}
		return nil
	}
	return
}
