package validation

import (
	"context"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Validator interface {
	Validate(ctx context.Context, duration *durationpb.Duration, fileSize int64) (err error, res bool)
}

type VoiceMsgValidator struct{}

func New() Validator {
	return &VoiceMsgValidator{}
}

func (v *VoiceMsgValidator) Validate(ctx context.Context, duration *durationpb.Duration, fileSize int64) (err error, res bool) {
	if fileSize < 2000 {
		return nil, false
	}
	return nil, true
}
