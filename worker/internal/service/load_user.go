package service

import (
	"context"
	"errors"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type UserCacheFinder interface {
	Load(context.Context, string) (domain.RdUser, error)
}

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type UserFinder interface {
	FindByID(context.Context, string) (domain.PgUser, error)
}

type LoadUserSvc interface {
	Execute(ctx context.Context, msg domain.Message, key string) (*domain.RdUser, error)
}

type loadUserSvc struct {
	userCacheFinder UserCacheFinder
	userFinder      UserFinder
}

func NewLoadUserSvc(userCacheFinder UserCacheFinder, userFinder UserFinder) LoadUserSvc {
	return &loadUserSvc{
		userCacheFinder: userCacheFinder,
		userFinder:      userFinder,
	}
}

func (s *loadUserSvc) Execute(ctx context.Context, msg domain.Message, key string) (*domain.RdUser, error) {
	u, err := s.userCacheFinder.Load(ctx, key)
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

		u = domain.RdUser{
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
