package usecase

import (
	"context"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"github.com/mefourr/tgdevob/worker/internal/dto"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"log/slog"
)

func (u *user) Save(ctx context.Context, entity *domain.RdUser, msg dto.Message) error {
	saveLastRequest(entity, msg)
	if !entity.MustValidated {
		entity.MustValidated = true
	}
	slog.DebugContext(ctx, "saving user to cache", "tg_user_id", entity.TgUserId, "must_validated", entity.MustValidated)
	// TODO: if we got an error while saving and our kafka_consumer sends the same message what next?
	if err := u.redis.Save(ctx, *entity); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to save user to cache", "tg_user_id", entity.TgUserId, "err", err)
		return err
	}
	slog.DebugContext(ctx, "user saved to cache", "tg_user_id", entity.TgUserId)
	return nil
}

func saveLastRequest(entity *domain.RdUser, msg dto.Message) {
	entity.LastRequest = msg.Request
}
