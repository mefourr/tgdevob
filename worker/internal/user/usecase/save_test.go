package usecase

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"github.com/mefourr/tgdevob/worker/internal/dto"
	"github.com/mefourr/tgdevob/worker/internal/user/usecase/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_service_SaveUser_BaseCase(t *testing.T) {
	type args struct {
		u   *domain.RdUser
		msg dto.Message
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "MustValidated",
			args: args{
				msg: dto.Message{Request: &tgbotapi.Message{}},
				u: &domain.RdUser{
					MustValidated: true,
					LastRequest:   &tgbotapi.Message{},
				},
			},
		},
		{
			name: "not MustValidated",
			args: args{
				msg: dto.Message{Request: &tgbotapi.Message{}},
				u: &domain.RdUser{
					MustValidated: false,
					LastRequest:   &tgbotapi.Message{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCacheSaver := mocks.NewMockRedis(t)
			s := New(userCacheSaver, nil)

			userCacheSaver.
				On("Save", mock.Anything, mock.MatchedBy(func(u domain.RdUser) bool {
					return u.MustValidated == true
				})).Return(nil).
				Once()

			err := s.Save(nil, tt.args.u, tt.args.msg)
			require.NoError(t, err)
		})
	}
}

func Test_service_SaveUser_Error(t *testing.T) {
	type args struct {
		msg dto.Message
		u   *domain.RdUser
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "MustValidated",
			args: args{
				msg: dto.Message{Request: &tgbotapi.Message{}},
				u: &domain.RdUser{
					MustValidated: true,
					LastRequest:   &tgbotapi.Message{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCacheStore := mocks.NewMockRedis(t)
			s := New(userCacheStore, nil)

			userCacheStore.
				On("Save", mock.Anything, *tt.args.u).
				Return(errors.New("some errors")).
				Once()

			err := s.Save(nil, tt.args.u, tt.args.msg)
			require.Error(t, err)
		})
	}
}
