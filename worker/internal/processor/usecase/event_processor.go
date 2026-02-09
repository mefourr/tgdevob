package usecase

import (
	"context"
	"github.com/mefourr/tgdevob/worker/internal/dto"
	user_uc "github.com/mefourr/tgdevob/worker/internal/user/usecase"
	validator_uc "github.com/mefourr/tgdevob/worker/internal/validator/usecase"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"log/slog"
	"strconv"
	"time"
)

type EventProcessor interface {
	Execute(ctx context.Context, message *dto.Message) error
}

type eventProcessor struct {
	user      user_uc.User
	validator validator_uc.Validator
}

func New(user user_uc.User, validator validator_uc.Validator) EventProcessor {
	return &eventProcessor{user: user, validator: validator}
}

func (e *eventProcessor) Execute(ctx context.Context, message *dto.Message) error {
	// load and then cache user
	u, err := e.user.Load(ctx, *message, strconv.FormatInt(message.User.ID, 10))
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to load user")
		return err
	}

	if u.MustValidated {
		// TODO: idempotency guarantee
		if IsNew(u.LastRequest.MessageID, message.Request.MessageID) {
			slog.ErrorContext(logger.ErrorCtx(ctx, nil), "Message that we just got has been already processed", "prev_message_id", u.LastRequest.MessageID, "current_message_id", message.Request.MessageID)
			// TODO: add logic with sending prev recognition result from postgres
			return nil
		}
	}

	slog.InfoContext(ctx, "updating user in cache", "user", u)
	go func() {
		if err := e.user.Save(ctx, u, *message); err != nil {
			slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to save user in cache")
		}
	}()

	// TODO: write test
	slog.InfoContext(ctx, "validating user's voice message generator", "generator", message.Request.Voice.Duration)
	if err = e.validator.Validate(ctx, time.Duration(message.Request.Voice.Duration)); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to validate voice message")
		return nil
	}

	// TODO: S3-storage authT -> download/upload voice msg
	// storage for voice message. Im gonna use yandex MessageS3SaverSvc-storage object storage
	//_ = e.storeVoiceUC.Save(*m)

	// TODO: recognition authT
	// recognition eventProcessor. Im gonna use yandex stt eventProcessor
	//rec, _ := e.recognizeUC.Save(*m)

	// TODO: send msg to tgbot by kafka_consumer broker
	//_ = e.sendResponseUC.Save(ctx, rec)

	slog.InfoContext(ctx, "Successfully consume a msg", "from", message.Request.From.ID)
	return nil
}

func IsNew(prev, curr int) bool {
	return prev != curr
}
