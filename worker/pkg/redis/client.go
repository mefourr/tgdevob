package redis

import (
	"context"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net"
)

type Client struct{}

func NewClient(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:               net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port),
		Password:           cfg.Redis.Password, // no password set
		Username:           cfg.Redis.Username,
		DB:                 cfg.Redis.Database, // use default DB
		DialerRetries:      cfg.Redis.DialerRetries,
		DialerRetryTimeout: cfg.Redis.DialerRetryTimeout,
	})
	c.FlushDB(ctx)
	ping := c.Ping(ctx)
	if _, err := ping.Result(); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "redis ping failed")
		return nil, err
	}
	return c, nil
}
