package cache

import (
	"encoding/json"
	"errors"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
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

func (c *Cache) Load(ctx context.Context, key string) (*User, error) {
	res, err := c.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, redis.Nil
	}

	var u User

	if err := json.Unmarshal([]byte(res), &u); err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "Unmarshalling user from Redis")
		return nil, err
	}

	slog.DebugContext(ctx, "Unmarshalling user from Redis", "user", u)
	return &u, nil
}

func (c *Cache) Save(ctx context.Context, user User) {
	bytes, err := json.Marshal(&user)
	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "error while marshalling user", err)
		panic(err)
	}

	// TODO: add timeout for ctx
	key := strconv.FormatInt(user.TgUserId, 10)
	if err = c.rdb.Set(ctx, key, bytes, time.Minute).Err(); err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "error while saving user to Redis")
		panic(err)
	}

	slog.DebugContext(ctx, "Saving user to Redis", "key", key, "user", user)
}
