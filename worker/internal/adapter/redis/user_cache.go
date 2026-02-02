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
	res, err := c.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return domain.RdUser{}, redis.Nil
	}

	var u domain.RdUser

	if err := json.Unmarshal([]byte(res), &u); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Unmarshalling user from Redis")
		return domain.RdUser{}, err
	}

	slog.DebugContext(ctx, "Unmarshalling user from Redis", "user", u)
	return u, nil
}

func (c *Cache) Save(ctx context.Context, user domain.RdUser) error {
	bytes, err := json.Marshal(&user)
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "error while marshalling user", "err:", err)
		return err
	}

	// TODO: add timeout for ctx
	key := strconv.FormatInt(user.TgUserId, 10)
	if err = c.rdb.Set(ctx, key, bytes, time.Minute).Err(); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "error while saving user to Redis")
		return err
	}

	slog.DebugContext(ctx, "Saving user to Redis", "key", key, "user", user)
	return nil
}
