package usecase

import (
	"context"
	"errors"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"github.com/mefourr/tgdevob/worker/internal/dto"
	mocks2 "github.com/mefourr/tgdevob/worker/internal/user/usecase/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_service_LoadUser_BaseCase(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
		key string
	}

	var (
		id        = "911e2171-9a61-4184-9fd8-ab1df149cb9b"
		tgUserId  = 123456789
		userName  = "Ivan1234"
		firstName = "Ivan"
	)
	tests := []struct {
		name    string
		args    args
		want    domain.PgUser
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
			userCacheFinder := mocks2.NewMockRedis(t)
			userFinder := mocks2.NewMockPostgres(t)
			u := New(
				userCacheFinder,
				userFinder,
			)

			userCacheFinder.
				On("Load", tt.args.ctx, tt.args.key).
				Return(domain.RdUser{}, redis.Nil).
				Once()
			userFinder.
				On("FindByID", tt.args.ctx, tt.args.key).
				Return(domain.PgUser{
					ID:        id,
					TgUserId:  int64(tgUserId),
					UserName:  &userName,
					FirstName: &firstName,
					LastName:  nil,
				}, nil).
				Once()

			got, err := u.Load(
				tt.args.ctx,
				dto.Message{UpdateId: 1, Request: nil},
				tt.args.key,
			)

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
		want domain.PgUser
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
			userCacheFinder := mocks2.NewMockRedis(t)
			s := New(userCacheFinder, nil)

			userCacheFinder.
				On("Load", tt.args.ctx, tt.args.key).
				Return(domain.RdUser{}, errors.New("some errors")).
				Once()

			_, err := s.Load(
				tt.args.ctx,
				dto.Message{UpdateId: 1, Request: nil},
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
		want domain.PgUser
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
			userCacheFinder := mocks2.NewMockRedis(t)
			userFinder := mocks2.NewMockPostgres(t)
			s := New(
				userCacheFinder,
				userFinder,
			)

			userCacheFinder.
				On("Load", tt.args.ctx, tt.args.key).
				Return(domain.RdUser{}, redis.Nil).
				Once()
			userFinder.
				On("FindByID", tt.args.ctx, tt.args.key).
				Return(domain.PgUser{}, errors.New("some errors")).
				Once()

			_, err := s.Load(
				tt.args.ctx,
				dto.Message{UpdateId: 1, Request: nil},
				tt.args.key,
			)
			require.Error(t, err)
		})
	}
}
