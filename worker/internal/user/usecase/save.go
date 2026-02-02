package usecase

import (
	"context"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"github.com/mefourr/tgdevob/worker/internal/dto"
	"log/slog"
)

func (u *user) Save(ctx context.Context, entity *domain.RdUser, msg dto.Message) error {
	saveLastRequest(entity, msg)
	if !entity.MustValidated {
		entity.MustValidated = true
	}
	// TODO: if we got an error while saving and our kafka_consumer sends the same message what next?
	if err := u.redis.Save(ctx, *entity); err != nil {
		slog.WarnContext(ctx, "failed to save user to rediscache", "key", entity.TgUserId, "error", err)
		return err
	}
	return nil
}

func saveLastRequest(entity *domain.RdUser, msg dto.Message) {
	entity.LastRequest = msg.Request
}
