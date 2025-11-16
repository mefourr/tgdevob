package loaduser

import (
	"context"
	"errors"
	"github.com/mefourr/tgdevob/worker/internal/kafka/message"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/cache"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/postgresql"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type UserCacheStore interface {
	Load(context.Context, string) (cache.User, error)
	Save(context.Context, cache.User) error
}

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type UserRepository interface {
	FindByID(context.Context, string) (postgresql.User, error)
}

type UserService interface {
	LoadUser(ctx context.Context, msg message.Message, key string) (*cache.User, error)
	SaveUser(ctx context.Context, u *cache.User, msg message.Message) error
}

type service struct {
	userCacheStore UserCacheStore
	userRepository UserRepository
}

func New(userCacheStore UserCacheStore, userRepository UserRepository) UserService {
	return &service{userCacheStore: userCacheStore, userRepository: userRepository}
}

func (s *service) LoadUser(ctx context.Context, msg message.Message, key string) (*cache.User, error) {
	u, err := s.userCacheStore.Load(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to load user from userCacheStore", "key", key)
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		entity, err := s.userRepository.FindByID(ctx, key)
		if err != nil {
			slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to find user by ID", "key", key)
			return nil, err
		}

		u = cache.User{
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

func (s *service) SaveUser(ctx context.Context, u *cache.User, msg message.Message) error {
	s.saveLastRequest(u, msg)
	if !u.MustValidated {
		u.MustValidated = true
	}
	// TODO: if we got an error while saving and our kafka sends the same message what next?
	if err := s.userCacheStore.Save(ctx, *u); err != nil {
		slog.WarnContext(ctx, "failed to save user to cache", "key", u.TgUserId, "error", err)
		return err
	}
	return nil
}

func (s *service) saveLastRequest(u *cache.User, msg message.Message) {
	u.LastRequest = msg.Request
}
