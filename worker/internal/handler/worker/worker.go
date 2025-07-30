package worker

import (
	"context"
	"encoding/json"
	"gihub.com/mefourr/tgdevob/worker/internal/storage/user/cache"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"strconv"
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
	m, err := takeMessage(msg)
	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to unmarshal worker")
		return err
	}

	slog.InfoContext(ctx, "ready to search in redis")
	var u cache.User
	key := "user:" + strconv.FormatInt(m.User.ID, 10)
	slog.InfoContext(ctx, "user has", "id", key)
	val, err := h.rdb.Get(ctx, key).Result()

	if err != redis.Nil { // TODO: works bad
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to get user from redis")
		u = cache.User{
			Id:        0,
			TgUserId:  strconv.FormatInt(m.User.ID, 10),
			UserName:  m.User.UserName,
			FirstName: m.User.FirstName,
			LastName:  m.User.LastName,
			UpdateId:  m.UpdateId,
		}
		cmd := h.rdb.HSet(ctx, "user:"+strconv.FormatInt(m.User.ID, 10), u)
		i, err := cmd.Result()
		if err != nil {
			slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to save user to redis", "Result", i)
			return err
		}
		slog.InfoContext(ctx, "user's saved")
	} else {
		if err = json.Unmarshal([]byte(val), &u); err != nil {
			slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to save user to redis", "Unmarshal", val)
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

func takeMessage(msg *sarama.ConsumerMessage) (*Message, error) {
	var m *Message
	err := json.Unmarshal(msg.Value, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
