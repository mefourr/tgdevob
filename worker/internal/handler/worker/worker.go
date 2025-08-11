package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gihub.com/mefourr/tgdevob/worker/internal/storage/user/cache"
	"gihub.com/mefourr/tgdevob/worker/internal/utils/msgutil"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

type Worker interface {
	Process(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type userRequest struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) Worker {
	return &userRequest{rdb: rdb}
}

func (h *userRequest) Process(ctx context.Context, msg *sarama.ConsumerMessage) error {
	// TODO: idempotency guarantee
	m, err := msgutil.Parse(msg)
	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to unmarshal worker")
		return err
	}

	var u cache.User
	key := fmt.Sprintf("user:%d", m.User.ID)
	slog.InfoContext(ctx, "ready to search in redis", "id", key)

	res, err := h.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		fmt.Println("key not found")
		u = cache.User{
			Id:          0,
			TgUserId:    key,
			UserName:    m.User.UserName,
			FirstName:   m.User.FirstName,
			LastName:    m.User.LastName,
			UpdateId:    m.UpdateId,
			LastRequest: m.Request,
		}

		bytes, err := json.Marshal(&u)
		if err != nil {
			panic(err)
		}
		if err := h.rdb.Set(ctx, key, bytes, time.Minute).Err(); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	} else {
		if err := json.Unmarshal([]byte(res), &u); err != nil {
			panic(err)
		}
	}

	slog.InfoContext(ctx, "Work with: ", "user", u)
	// TODO: cache user
	// TODO: validate worker
	// TODO: s3 grpc
	// TODO: recognition grpc
	slog.InfoContext(ctx, "Successfully consume a msg", "from", m.Request.From.UserName)
	return nil
}
