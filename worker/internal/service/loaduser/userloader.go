package loaduser

import (
	"context"
	"errors"
	"github.com/mefourr/tgdevob/worker/internal/kafka/message"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/postgresql"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/rediscache"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type UserCacheStore interface {
	Load(context.Context, string) (rediscache.User, error)
	Save(context.Context, rediscache.User) error
}

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type UserFinder interface {
	FindByID(context.Context, string) (postgresql.User, error)
}

type UserService interface {
	LoadUser(ctx context.Context, msg message.Message, key string) (*rediscache.User, error)
	SaveUser(ctx context.Context, u *rediscache.User, msg message.Message) error
}

type service struct {
	userCacheStore UserCacheStore
	userFinder     UserFinder
}

func New(userCacheStore UserCacheStore, userFinder UserFinder) UserService {
	return &service{userCacheStore: userCacheStore, userFinder: userFinder}
}

func (s *service) LoadUser(ctx context.Context, msg message.Message, key string) (*rediscache.User, error) {
	u, err := s.userCacheStore.Load(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to load user from userCacheStore", "key", key)
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		entity, err := s.userFinder.FindByID(ctx, key)
		if err != nil {
			slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to find user by ID", "key", key)
			return nil, err
		}

		u = rediscache.User{
			Id:            entity.ID,
			TgUserId:      entity.TgUserId,
			UserName:      entity.UserName,
			FirstName:     entity.FirstName,
			LastName:      entity.LastName,
			UpdateId:      msg.UpdateId,
			MustValidated: false,
		}
	}

	slog.DebugContext(ctx, "loaded user", "key", key, "user", u)
	return &u, nil
}

func (s *service) SaveUser(ctx context.Context, u *rediscache.User, msg message.Message) error {
	s.saveLastRequest(u, msg)
	if !u.MustValidated {
		u.MustValidated = true
	}
	// TODO: if we got an error while saving and our kafka sends the same message what next?
	if err := s.userCacheStore.Save(ctx, *u); err != nil {
		slog.WarnContext(ctx, "failed to save user to rediscache", "key", u.TgUserId, "error", err)
		return err
	}
	return nil
}

func (s *service) saveLastRequest(u *rediscache.User, msg message.Message) {
	u.LastRequest = msg.Request
}
