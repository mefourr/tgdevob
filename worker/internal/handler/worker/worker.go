package worker

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/proto/voice-msg-validator/v1/pb"
	"github.com/mefourr/tgdevob/worker/internal/handler/userloader"
	"github.com/mefourr/tgdevob/worker/internal/kafka/idem"
	"github.com/mefourr/tgdevob/worker/internal/utils/deserial"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"log/slog"
	"strconv"
	"sync"
	"time"
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
	// audio length and error to user bout validating error
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		conn, err := grpc.NewClient("localhost:5353", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewValidateVMLengthClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		slog.InfoContext(ctx, "Try to send msg to validator")

		r, err := c.ValidateVMLength(ctx, &pb.VoiceMessageDataRq{
			FileSize: int64(m.Request.Voice.FileSize),
			Duration: int64(m.Request.Voice.Duration),
		})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Result: %t", r.GetIsValidated())
	}()
	wg.Wait()

	// TODO: s3-storage grpc -> download/upload voice msg
	// storage for voice message. Im gonna use yandex s3-storage object storage

	// TODO: recognition grpc
	// recognition service. Im gonna use yandex stt service
	slog.InfoContext(ctx, "Successfully consume a msg", "from", m.Request.From.ID)
	return nil
}
