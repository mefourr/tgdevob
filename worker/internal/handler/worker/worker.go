package worker

import (
	"context"
	s "gihub.com/mefourr/tgdevob/worker/internal/handler/userloader"
	"gihub.com/mefourr/tgdevob/worker/internal/utils/deserial"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/IBM/sarama"
	"log/slog"
	"strconv"
)

type Worker interface {
	Process(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type userRequest struct {
	loader userloader.UserLoader
}

func New(loader userloader.UserLoader) Worker {
	return &userRequest{loader: loader}
}

func (ur *userRequest) Process(ctx context.Context, msg *sarama.ConsumerMessage) error {
	// TODO: idempotency guarantee
	m, err := deserial.ParseMessage(msg)
	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to unmarshal worker")
		return err
	}

	// TODO: cache user
	// TODO: noticed that messages continue to be sent if panic occurs - maybe it can overload a cache handler
	// --- throw a panic while caching to see behavior
	_, err = ur.loader.LoadUser(ctx, m, strconv.FormatInt(m.User.ID, 10))
	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to load user")
		return err
	}

	// TODO: validate worker
	// TODO: s3 grpc
	// TODO: recognition grpc
	slog.InfoContext(ctx, "Successfully consume a msg", "from", m.Request.From.ID)
	return nil
}
