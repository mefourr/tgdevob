package worker

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/worker/internal/kafka/message"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/rediscache"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"log/slog"
	"strconv"
)

type MessageHandler interface {
	Process(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type ResponseProducer interface {
	Produce(context.Context, string) error
}

type MessageRecognizer interface {
	Recognize(msg message.Message) (string, error)
}

type MessageS3Saver interface {
	Save(message.Message) error
}

type VoiceMessageValidator interface {
	Validate(duration, fileSize int) bool
}

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type IdempotencyChecker interface {
	Check(u *rediscache.User, msg message.Message) bool
}

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type MessageParser interface {
	Parse(res *message.Message, cm *sarama.ConsumerMessage) error
}

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type UserService interface {
	LoadUser(ctx context.Context, msg message.Message, key string) (*rediscache.User, error)
	SaveUser(ctx context.Context, u *rediscache.User, msg message.Message) error
}

type handler struct {
	userService           UserService
	messageParser         MessageParser
	idempotencyChecker    IdempotencyChecker
	voiceMessageValidator VoiceMessageValidator
	s3                    MessageS3Saver
	messageRecognizer     MessageRecognizer
	responseProducer      ResponseProducer
}

func New(userService UserService, messageParser MessageParser, idempotencyChecker IdempotencyChecker, voiceMessageValidator VoiceMessageValidator, s3 MessageS3Saver, messageRecognizer MessageRecognizer, responseProducer ResponseProducer) *handler {
	return &handler{userService: userService, messageParser: messageParser, idempotencyChecker: idempotencyChecker, voiceMessageValidator: voiceMessageValidator, s3: s3, messageRecognizer: messageRecognizer, responseProducer: responseProducer}
}

func (h *handler) Process(ctx context.Context, msg *sarama.ConsumerMessage) error {
	m := message.Message{}
	if err := h.messageParser.Parse(&m, msg); err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to unmarshal handler")
		return err
	}

	// load and then cache user
	u, err := h.userService.LoadUser(ctx, m, strconv.FormatInt(m.User.ID, 10))
	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to load user")
		return err
	}

	if u.MustValidated {
		// TODO: idempotency guarantee
		if ok := h.idempotencyChecker.Check(u, m); !ok {
			slog.ErrorContext(logging.ErrorCtx(ctx, nil), "Message that we just got has been already processed", "prev_message_id", u.LastRequest.MessageID, "current_message_id", m.Request.MessageID)
			// TODO: add logic with sending prev recognition result from postgres
			return nil
		}
	}

	slog.InfoContext(ctx, "updating user in cache", "user", u)
	go func() {
		_ = h.userService.SaveUser(ctx, u, m)
	}()

	// TODO: validate handler
	// audio length and error to user bout validating error
	_ = h.voiceMessageValidator.Validate(0, 0)

	// TODO: S3-storage grpcapp -> download/upload voice msg
	// storage for voice message. Im gonna use yandex MessageS3Saver-storage object storage
	_ = h.s3.Save(m)

	// TODO: recognition grpcapp
	// recognition service. Im gonna use yandex stt service
	rec, _ := h.messageRecognizer.Recognize(m)

	// TODO: send msg to tgbot by kafka broker
	_ = h.responseProducer.Produce(ctx, rec)

	slog.InfoContext(ctx, "Successfully consume a msg", "from", m.Request.From.ID)
	return nil
}
