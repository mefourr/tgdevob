package userinfo

import (
	"context"
	"errors"
	"gihub.com/mefourr/tgdevob/worker/internal/storage/user/cache"
	"github.com/redis/go-redis/v9"
)

type CachedUser interface {
	Load(context.Context, string) (*cache.User, bool)
}

func LoadUser(ctx context.Context, ca cache.Cache, key string) (*cache.User, error) {
	_, err := ca.Load(ctx, key)
	if errors.Is(err, redis.Nil) {
		//entity, err := posgresql.GetUserById(key)
		//if err != nil {return nil, err}
		//u = cache.User{
		//	entity.Id,
		//	entity.Name,
		//	entity.etc...
	} else if err != nil {
		return nil, err
	}
	//_ := ca.Store(ctx, key, u, time.Minute)
	panic("implement me")
}
