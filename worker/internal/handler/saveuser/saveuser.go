package saveuser

import (
	"context"
	"gihub.com/mefourr/tgdevob/worker/internal/storage/user/cache"
	"gihub.com/mefourr/tgdevob/worker/internal/storage/user/posgresql"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"time"
)

type Welcome struct {
	Cache    cache.UserCache
	Postgres posgresql.UserRepository
}

type StoreUser interface {
	Create(ctx context.Context, user *posgresql.User) error
}

func (w *Welcome) SaveUser(ctx context.Context, data *tgbotapi.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u := &posgresql.User{
		TgUserId:  data.ID,
		UserName:  &data.UserName,
		FirstName: &data.FirstName,
		LastName:  &data.LastName,
	}

	if err := w.Postgres.Create(ctx, u); err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to create user in postgres", "tgUserId", data.ID)
		return "n/a", err
	}

	// TODO: immediately add to cache
	//cachedUser := cache.User{
	//	Id:        u.ID,
	//	TgUserId:  u.TgUserId,
	//	UserName:  u.UserName,
	//	FirstName: u.FirstName,
	//	LastName:  u.LastName,
	//}
	//
	//go func() {
	//	if err := w.Cache.Save(ctx, cachedUser); err != nil {
	//		slog.WarnContext(ctx, "failed to save user to cache", "tgUserId", u.ID, "error", err)
	//	}
	//}()

	slog.InfoContext(ctx, "user created successfully", "tgUserId", u.ID)
	return u.ID, nil
}
