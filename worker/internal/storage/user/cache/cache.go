package cache

import (
	"encoding/json"
	"errors"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"log/slog"
)

type Cache struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Cache {
	return &Cache{rdb: rdb}
}

func (c *Cache) Load(ctx context.Context, id string) (*User, error) {
	res, err := c.rdb.Get(ctx, id).Result()

	if errors.Is(err, redis.Nil) {
		return nil, redis.Nil
	}

	var u User

	if err := json.Unmarshal([]byte(res), &u); err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "Unmarshalling user from Redis")
		return nil, err
	}

	return &u, nil
}

func (c *Cache) Save(user User) error {
	//TODO implement me
	panic("implement me")
}
