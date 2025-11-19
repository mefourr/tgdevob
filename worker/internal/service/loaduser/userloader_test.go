package loaduser

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mefourr/tgdevob/worker/internal/kafka/message"
	"github.com/mefourr/tgdevob/worker/internal/service/loaduser/mocks"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/postgresql"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/rediscache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	_ "reflect"
	"testing"
)

func Test_service_LoadUser_BaseCase(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
		key string
	}

	var (
		id        = "901e2676-9a60-4284-9fd8-ab1df149cb9b"
		tgUserId  = 367466842
		userName  = "slavaoli"
		firstName = "Jeff"
	)
	tests := []struct {
		name    string
		args    args
		want    postgresql.User
		wantErr bool
	}{
		{
			name: "base test",
			args: args{
				ctx: ctx,
				key: "367466842",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepository := mocks.NewMockUserRepository(t)
			userCacheStore := mocks.NewMockUserCacheStore(t)

			s := New(userCacheStore, userRepository)

			userCacheStore.
				On("Load", tt.args.ctx, tt.args.key).
				Return(rediscache.User{}, redis.Nil).
				Once()
			userRepository.
				On("FindByID", tt.args.ctx, tt.args.key).
				Return(postgresql.User{
					ID:        id,
					TgUserId:  int64(tgUserId),
					UserName:  &userName,
					FirstName: &firstName,
					LastName:  nil,
				}, nil).
				Once()

			got, err := s.LoadUser(tt.args.ctx, message.Message{UpdateId: 1, Request: nil}, tt.args.key)

			assert.Nil(t, err)
			assert.Equal(t, id, got.Id)
			assert.Equal(t, int64(tgUserId), got.TgUserId)
			assert.Equal(t, userName, *got.UserName)
			assert.Equal(t, firstName, *got.FirstName)
			assert.Equal(t, 1, got.UpdateId)
			assert.Nil(t, got.LastRequest)
			assert.False(t, got.MustValidated)
		})
	}
}

func Test_service_LoadUser_UnexpectedError(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name string
		args args
		want postgresql.User
	}{
		{
			name: "unexpected error",
			args: args{
				ctx: nil,
				key: "333",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCacheStore := mocks.NewMockUserCacheStore(t)
			s := New(userCacheStore, nil)

			userCacheStore.
				On("Load", tt.args.ctx, tt.args.key).
				Return(rediscache.User{}, errors.New("some errors")).
				Once()

			_, err := s.LoadUser(
				tt.args.ctx,
				message.Message{UpdateId: 1, Request: nil},
				tt.args.key,
			)
			require.Error(t, err)
		})
	}
}

func Test_service_LoadUser_FailedFindByID(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name string
		args args
		want postgresql.User
	}{
		{
			name: "unexpected error",
			args: args{
				ctx: nil,
				key: "333",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCacheStore := mocks.NewMockUserCacheStore(t)
			userRepository := mocks.NewMockUserRepository(t)

			s := New(userCacheStore, userRepository)

			userCacheStore.
				On("Load", tt.args.ctx, tt.args.key).
				Return(rediscache.User{}, redis.Nil).
				Once()
			userRepository.
				On("FindByID", tt.args.ctx, tt.args.key).
				Return(postgresql.User{}, errors.New("some errors")).
				Once()

			_, err := s.LoadUser(
				tt.args.ctx,
				message.Message{UpdateId: 1, Request: nil},
				tt.args.key,
			)
			require.Error(t, err)
		})
	}
}

func Test_service_SaveUser_BaseCase(t *testing.T) {
	type args struct {
		msg message.Message
		u   *rediscache.User
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "MustValidated",
			args: args{
				msg: message.Message{Request: &tgbotapi.Message{}},
				u: &rediscache.User{
					MustValidated: true,
					LastRequest:   &tgbotapi.Message{},
				},
			},
		},
		{
			name: "not MustValidated",
			args: args{
				msg: message.Message{Request: &tgbotapi.Message{}},
				u: &rediscache.User{
					MustValidated: false,
					LastRequest:   &tgbotapi.Message{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCacheStore := mocks.NewMockUserCacheStore(t)
			s := New(userCacheStore, nil)

			userCacheStore.
				On("Save", mock.Anything, mock.MatchedBy(func(u rediscache.User) bool {
					return u.MustValidated == true
				})).Return(nil).
				Once()

			err := s.SaveUser(nil, tt.args.u, tt.args.msg)
			require.NoError(t, err)
		})
	}
}

func Test_service_SaveUser_Error(t *testing.T) {
	type args struct {
		msg message.Message
		u   *rediscache.User
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "MustValidated",
			args: args{
				msg: message.Message{Request: &tgbotapi.Message{}},
				u: &rediscache.User{
					MustValidated: true,
					LastRequest:   &tgbotapi.Message{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCacheStore := mocks.NewMockUserCacheStore(t)
			s := New(userCacheStore, nil)

			userCacheStore.
				On("Save", mock.Anything, *tt.args.u).
				Return(errors.New("some errors")).
				Once()

			err := s.SaveUser(nil, tt.args.u, tt.args.msg)
			require.Error(t, err)
		})
	}
}
