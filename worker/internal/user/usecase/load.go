package usecase

import (
	"context"
	"errors"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"github.com/mefourr/tgdevob/worker/internal/dto"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

func (s *user) Load(ctx context.Context, msg dto.Message, key string) (*domain.RdUser, error) {
	u, err := s.redis.Load(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to load user from cache", "key", key, "err", err)
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		slog.DebugContext(ctx, "user not in cache, falling back to postgres", "key", key)
		entity, err := s.postgres.FindByID(ctx, key)
		if err != nil {
			slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to find user in postgres", "key", key, "err", err)
			return nil, err
		}
		slog.DebugContext(ctx, "user loaded from postgres", "key", key, "user_id", entity.ID)

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

	slog.DebugContext(ctx, "user loaded", "key", key, "must_validated", u.MustValidated)
	return &u, nil
}
