package usecase

import "context"

type YandexCloud interface {
	CreateBucket(context.Context) (string, error)
}
