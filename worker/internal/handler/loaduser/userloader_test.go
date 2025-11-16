package loaduser

import (
	"context"
	"github.com/mefourr/tgdevob/worker/internal/handler/loaduser/mocks"
	"github.com/mefourr/tgdevob/worker/internal/kafka/message"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/cache"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/postgresql"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	_ "reflect"
	"testing"
)

func Test_clientRetriever_LoadUser(t *testing.T) {
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

			s := &service{
				userCacheStore: userCacheStore,
				userRepository: userRepository,
			}

			userCacheStore.
				On("Load", tt.args.ctx, tt.args.key).
				Return(cache.User{}, redis.Nil).
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
