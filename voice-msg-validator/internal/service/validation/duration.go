package validation

import (
	"context"
	"errors"
	"log/slog"
)

type Validator interface {
	Validate(ctx context.Context, duration float32) error
}

type VoiceMsgValidator struct{}

func New() Validator {
	return &VoiceMsgValidator{}
}

// TODO: move to config
const limit float32 = 2000

var ErrDurationIsTooLong = errors.New("duration limit (2000sec) exceeded")

func (v VoiceMsgValidator) Validate(ctx context.Context, duration float32) error {
	slog.InfoContext(ctx, "validating duration", "duration", duration, "limit", limit)
	if duration > limit {
		return ErrDurationIsTooLong
	}
	return nil
}
