package kafka

import (
	"context"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"github.com/mefourr/tgdevob/worker/pkg/logging"
	"log/slog"
	"strconv"
	"time"
)

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type VoiceDurationValidator interface {
	Execute(ctx context.Context, duration time.Duration) error
}

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type IdempotencyCheckerSvc interface {
	Execute(prev, curr int) bool
}

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type UserRequestParserSvc interface {
	Execute(data []byte) (*domain.Message, error)
}

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type LoadUserSvc interface {
	Execute(ctx context.Context, msg domain.Message, key string) (*domain.RdUser, error)
}

//go:generate go run github.com/vektra/mockery/v3@v3.6.0
type SaveUserSvc interface {
	Execute(ctx context.Context, u *domain.RdUser, msg domain.Message) error
}

type EventProcessor interface {
	Execute(ctx context.Context, e domain.Event) error
}

type eventProcessor struct {
	userRequestParser      UserRequestParserSvc
	loadUser               LoadUserSvc
	saveUser               SaveUserSvc
	idempotencyChecker     IdempotencyCheckerSvc
	voiceDurationValidator VoiceDurationValidator
}

func NewEventProcessor(userRequestParser UserRequestParserSvc, loadUser LoadUserSvc, saveUser SaveUserSvc, idempotencyChecker IdempotencyCheckerSvc, voiceDurationValidator VoiceDurationValidator) *eventProcessor {
	return &eventProcessor{userRequestParser: userRequestParser, loadUser: loadUser, saveUser: saveUser, idempotencyChecker: idempotencyChecker, voiceDurationValidator: voiceDurationValidator}
}

func (e *eventProcessor) Execute(ctx context.Context, event domain.Event) error {
	m, err := e.userRequestParser.Execute(event.Value)
	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to unmarshal eventProcessor")
		return err
	}

	// load and then cache user
	u, err := e.loadUser.Execute(ctx, *m, strconv.FormatInt(m.User.ID, 10))
	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to load user")
		return err
	}

	if u.MustValidated {
		// TODO: idempotency guarantee
		if ok := e.idempotencyChecker.Execute(u.LastRequest.MessageID, m.Request.MessageID); !ok {
			slog.ErrorContext(logging.ErrorCtx(ctx, nil), "Message that we just got has been already processed", "prev_message_id", u.LastRequest.MessageID, "current_message_id", m.Request.MessageID)
			// TODO: add logic with sending prev recognition result from postgres
			return nil
		}
	}

	slog.InfoContext(ctx, "updating user in cache", "user", u)
	go func() {
		if err := e.saveUser.Execute(ctx, u, *m); err != nil {
			slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to save user in cache")
		}
	}()

	// TODO: write test
	slog.InfoContext(ctx, "validating user's voice message duration", "duration", m.Request.Voice.Duration)
	if err = e.voiceDurationValidator.Execute(ctx, time.Duration(m.Request.Voice.Duration)); err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to validate voice duration")
		return nil
	}

	// TODO: S3-storage grpcapp -> download/upload voice msg
	// storage for voice message. Im gonna use yandex MessageS3SaverSvc-storage object storage
	//_ = e.storeVoiceUC.Execute(*m)

	// TODO: recognition grpcapp
	// recognition eventProcessor. Im gonna use yandex stt eventProcessor
	//rec, _ := e.recognizeUC.Execute(*m)

	// TODO: send msg to tgbot by kafka broker
	//_ = e.sendResponseUC.Execute(ctx, rec)

	slog.InfoContext(ctx, "Successfully consume a msg", "from", m.Request.From.ID)
	return nil
}
