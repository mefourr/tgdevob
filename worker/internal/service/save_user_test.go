package service

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"github.com/mefourr/tgdevob/worker/internal/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_service_SaveUser_BaseCase(t *testing.T) {
	type args struct {
		u   *domain.RdUser
		msg domain.Message
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "MustValidated",
			args: args{
				msg: domain.Message{Request: &tgbotapi.Message{}},
				u: &domain.RdUser{
					MustValidated: true,
					LastRequest:   &tgbotapi.Message{},
				},
			},
		},
		{
			name: "not MustValidated",
			args: args{
				msg: domain.Message{Request: &tgbotapi.Message{}},
				u: &domain.RdUser{
					MustValidated: false,
					LastRequest:   &tgbotapi.Message{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCacheSaver := mocks.NewMockUserCacheSaver(t)
			s := NewSaveUserSvc(userCacheSaver)

			userCacheSaver.
				On("Save", mock.Anything, mock.MatchedBy(func(u domain.RdUser) bool {
					return u.MustValidated == true
				})).Return(nil).
				Once()

			err := s.Execute(nil, tt.args.u, tt.args.msg)
			require.NoError(t, err)
		})
	}
}

func Test_service_SaveUser_Error(t *testing.T) {
	type args struct {
		msg domain.Message
		u   *domain.RdUser
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "MustValidated",
			args: args{
				msg: domain.Message{Request: &tgbotapi.Message{}},
				u: &domain.RdUser{
					MustValidated: true,
					LastRequest:   &tgbotapi.Message{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCacheStore := mocks.NewMockUserCacheSaver(t)
			s := NewSaveUserSvc(userCacheStore)

			userCacheStore.
				On("Save", mock.Anything, *tt.args.u).
				Return(errors.New("some errors")).
				Once()

			err := s.Execute(nil, tt.args.u, tt.args.msg)
			require.Error(t, err)
		})
	}
}
