package utils

import (
	"context"
	"gihub.com/mefourr/tgdevob/worker/config"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
	"log/slog"
	"time"
)

func TryConn(fn func() error, conf config.StorageConfig) (err error) {
	for conf.RetryNum > 0 {
		if err = fn(); err != nil {
			slog.ErrorContext(logging.ErrorCtx(context.TODO(), err), "cant connect. trt again", "retries", conf.RetryNum)
			time.Sleep(conf.Delay)
			conf.RetryNum--
			continue
		}
		return nil
	}
	return
}
