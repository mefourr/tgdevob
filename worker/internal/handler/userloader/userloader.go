package userloader

import (
	"context"
	"errors"
	"github.com/mefourr/tgdevob/worker/internal/kafka/message"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/cache"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/posgresql"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type Client interface {
	LoadUser(ctx context.Context, msg message.Message, key string) (*cache.User, error)
	SaveOrUpdateUser(ctx context.Context, u *cache.User)
	// SaveUserLastRequest just save Request as last before user caching
	SaveUserLastRequest(u *cache.User, msg message.Message)
}

type clientRetriever struct {
	Cache    cache.UserCache
	Postgres posgresql.UserRepository
}

func New(cache cache.UserCache, postgres posgresql.UserRepository) Client {
	return &clientRetriever{Cache: cache, Postgres: postgres}
}

type CachedUser interface {
	Load(context.Context, string) (cache.User, error)
	Save(context.Context, cache.User) error
}

type PgUserLoader interface {
	FindById(context.Context, string) (posgresql.User, error)
}

func (cr *clientRetriever) LoadUser(ctx context.Context, msg message.Message, key string) (*cache.User, error) {
	u, err := cr.Cache.Load(ctx, key)

	if err != nil && !errors.Is(err, redis.Nil) {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to load user from cache", "key", key)
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		entity, err := cr.Postgres.FindById(ctx, key)
		if err != nil {
			slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to find user by ID", "key", key)
			return nil, err
		}

		u = cache.User{
			Id:          entity.ID,
			TgUserId:    entity.TgUserId,
			UserName:    entity.UserName,
			FirstName:   entity.FirstName,
			LastName:    entity.LastName,
			UpdateId:    msg.UpdateId,
			LastRequest: msg.Request, // TODO: review if needed
			IsNonCached: true,
		}
		//go cr.SaveOrUpdateUser(ctx, u, key)
	}

	slog.DebugContext(ctx, "loaded user", "key", key, "user", u)
	return &u, nil
}

func (cr *clientRetriever) SaveOrUpdateUser(ctx context.Context, u *cache.User) {
	u.IsNonCached = false
	if err := cr.Cache.Save(ctx, *u); err != nil {
		slog.WarnContext(ctx, "failed to save user to cache", "key", u.TgUserId, "error", err)
		u.IsNonCached = true
	}
}

func (cr *clientRetriever) SaveUserLastRequest(u *cache.User, msg message.Message) {
	u.LastRequest = msg.Request
}
