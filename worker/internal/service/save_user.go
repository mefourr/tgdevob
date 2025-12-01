package service

import (
	"context"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"log/slog"
)

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type UserCacheSaver interface {
	Save(context.Context, domain.RdUser) error
}

type SaveUserSvc interface {
	Execute(ctx context.Context, u *domain.RdUser, msg domain.Message) error
}

type saveUserSvc struct {
	userCacheSaver UserCacheSaver
}

func NewSaveUserSvc(userCacheSaver UserCacheSaver) SaveUserSvc {
	return &saveUserSvc{userCacheSaver: userCacheSaver}
}

func (s *saveUserSvc) Execute(ctx context.Context, u *domain.RdUser, msg domain.Message) error {
	saveLastRequest(u, msg)
	if !u.MustValidated {
		u.MustValidated = true
	}
	// TODO: if we got an error while saving and our kafka sends the same message what next?
	if err := s.userCacheSaver.Save(ctx, *u); err != nil {
		slog.WarnContext(ctx, "failed to save user to rediscache", "key", u.TgUserId, "error", err)
		return err
	}
	return nil
}

func saveLastRequest(u *domain.RdUser, msg domain.Message) {
	u.LastRequest = msg.Request
}
