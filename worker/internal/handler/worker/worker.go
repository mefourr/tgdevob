package worker

import (
	"context"
	"gihub.com/mefourr/tgdevob/worker/internal/handler/userloader"
	"gihub.com/mefourr/tgdevob/worker/internal/kafka/idem"
	"gihub.com/mefourr/tgdevob/worker/internal/utils/deserial"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/IBM/sarama"
	"log/slog"
	"strconv"
)

type Worker interface {
	Process(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type rqWorker struct {
	client userloader.Client
}

func New(client userloader.Client) Worker {
	return &rqWorker{client: client}
}

func (rqw *rqWorker) Process(ctx context.Context, msg *sarama.ConsumerMessage) error {
	m, err := deserial.ParseMessage(msg)
	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to unmarshal worker")
		return err
	}

	// TODO: cache user
	u, err := rqw.client.LoadUser(ctx, m, strconv.FormatInt(m.User.ID, 10))
	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to load user")
		return err
	}

	if !u.IsNonCached {
		// TODO: idempotency guarantee
		if ok := idem.Validate(u, m); !ok {
			slog.ErrorContext(logging.ErrorCtx(ctx, nil), "Message that we just got has been already processed", "prev_message_id", u.LastRequest.MessageID, "current_message_id", m.Request.MessageID)
			// TODO: add logic with sending prev recognition result from postgres
			return nil
		}
	}

	slog.InfoContext(ctx, "updating user cache", "user", u)
	rqw.client.SaveUserLastRequest(u, m)
	go rqw.client.SaveOrUpdateUser(ctx, u)

	// TODO: validate worker
	// means i have to check voice message time maybe smth else

	// TODO: s3 grpc
	// storage for voice message. Im gonna use yandex s3 object storage

	// TODO: recognition grpc
	// recognition service. Im gonna use yandex stt service
	slog.InfoContext(ctx, "Successfully consume a msg", "from", m.Request.From.ID)
	return nil
}
