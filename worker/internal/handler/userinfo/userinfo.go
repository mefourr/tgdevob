package userinfo

import (
	"context"
	"errors"
	"gihub.com/mefourr/tgdevob/worker/internal/kafka/message"
	"gihub.com/mefourr/tgdevob/worker/internal/storage/user/cache"
	"gihub.com/mefourr/tgdevob/worker/internal/storage/user/posgresql"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type UserLoader struct {
	Cache    cache.UserCache
	Postgres posgresql.UserRepository
}

type CachedUser interface {
	Load(context.Context, string) (cache.User, error)
}

type LoadedUser interface {
	FindById(context.Context, string) (posgresql.User, error)
}

func (ui *UserLoader) LoadUser(ctx context.Context, msg message.Message, key string) (*cache.User, error) {
	u, err := ui.Cache.Load(ctx, key)

	if err != nil && !errors.Is(err, redis.Nil) {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to load user from cache", "key", key)
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		entity, err := ui.Postgres.FindById(ctx, key)
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
		}

		go func() {
			if err := ui.Cache.Save(ctx, u); err != nil {
				slog.WarnContext(ctx, "failed to save user to cache", "key", key, "error", err)
			}
		}()
	}

	slog.DebugContext(ctx, "loaded user", "key", key, "user", u)
	return &u, nil
}
