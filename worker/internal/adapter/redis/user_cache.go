package rediscache

import (
	"encoding/json"
	"errors"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"log/slog"
	"strconv"
	"time"
)

type Cache struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Cache {
	return &Cache{rdb: rdb}
}

func (c *Cache) Load(ctx context.Context, key string) (domain.RdUser, error) {
	slog.DebugContext(ctx, "loading user from redis cache", "key", key)
	res, err := c.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		slog.DebugContext(ctx, "cache miss", "key", key)
		return domain.RdUser{}, redis.Nil
	}

	var u domain.RdUser
	if err := json.Unmarshal([]byte(res), &u); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to deserialize user from redis", "key", key, "err", err)
		return domain.RdUser{}, err
	}

	slog.DebugContext(ctx, "cache hit", "key", key, "tg_user_id", u.TgUserId)
	return u, nil
}

func (c *Cache) Save(ctx context.Context, user domain.RdUser) error {
	key := strconv.FormatInt(user.TgUserId, 10)
	slog.DebugContext(ctx, "saving user to redis cache", "key", key)

	bytes, err := json.Marshal(&user)
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to serialize user for redis", "key", key, "err", err)
		return err
	}

	// TODO: add timeout for ctx
	if err = c.rdb.Set(ctx, key, bytes, time.Minute).Err(); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to write user to redis", "key", key, "err", err)
		return err
	}

	slog.DebugContext(ctx, "user saved to redis cache", "key", key)
	return nil
}
