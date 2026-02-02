package usecase

import (
	"context"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"github.com/mefourr/tgdevob/worker/internal/dto"
)

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type Redis interface {
	Load(context.Context, string) (domain.RdUser, error)
	Save(context.Context, domain.RdUser) error
}

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type Postgres interface {
	FindByID(context.Context, string) (domain.PgUser, error)
}

type User interface {
	Save(ctx context.Context, entity *domain.RdUser, msg dto.Message) error
	Load(ctx context.Context, msg dto.Message, key string) (*domain.RdUser, error)
}

type user struct {
	redis    Redis
	postgres Postgres
}

func New(redis Redis, postgres Postgres) User {
	return &user{redis: redis, postgres: postgres}
}
