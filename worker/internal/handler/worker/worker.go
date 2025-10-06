package worker

import (
	"context"
	"fmt"
	"gihub.com/mefourr/tgdevob/worker/internal/handler/userinfo"
	"gihub.com/mefourr/tgdevob/worker/internal/utils/deserial"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/IBM/sarama"
	"log/slog"
)

type Worker interface {
	Process(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type userRequest struct {
	loader userinfo.UserLoader
}

func New(loader userinfo.UserLoader) Worker {
	return &userRequest{loader: loader}
}

func (h *userRequest) Process(ctx context.Context, msg *sarama.ConsumerMessage) error {
	// TODO: idempotency guarantee
	m, err := deserial.ParseMessage(msg)
	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to unmarshal worker")
		return err
	}

	_, err = h.loader.LoadUser(ctx, m, fmt.Sprintf("user:%d", m.User.ID))
	if err != nil {
		return err
	}
	// TODO: cache user
	// TODO: validate worker
	// TODO: s3 grpc
	// TODO: recognition grpc
	slog.InfoContext(ctx, "Successfully consume a msg", "from", m.Request.From.UserName)
	return nil
}
