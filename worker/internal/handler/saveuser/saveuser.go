package saveuser

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/postgresql"
	"log/slog"
	"time"
)

type UserCreateRepository interface {
	Create(ctx context.Context, user *postgresql.User) error
}

type Welcome struct {
	repository UserCreateRepository
}

func (w *Welcome) SaveUser(ctx context.Context, data *tgbotapi.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u := &postgresql.User{
		TgUserId:  data.ID,
		UserName:  &data.UserName,
		FirstName: &data.FirstName,
		LastName:  &data.LastName,
	}

	if err := w.repository.Create(ctx, u); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			slog.WarnContext(ctx, "database request timed out", "ID", u.ID)
			return "", ctx.Err()
		}
		return "", err
	}

	// TODO: immediately add to cache
	//cachedUser := cache.SaveUser{
	//	Id:        u.ID,
	//	TgUserId:  u.TgUserId,
	//	UserName:  u.UserName,
	//	FirstName: u.FirstName,
	//	LastName:  u.LastName,
	//}
	//
	//go func() {
	//	if err := w.Store.Save(ctx, cachedUser); err != nil {
	//		slog.WarnContext(ctx, "failed to save user to cache", "tgUserId", u.ID, "error", err)
	//	}
	//}()

	slog.InfoContext(ctx, "user created successfully", "tgUserId", u.ID)
	return u.ID, nil
}
