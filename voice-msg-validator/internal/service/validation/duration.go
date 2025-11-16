package validation

import (
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Validator interface {
	Validate(ctx context.Context, duration *durationpb.Duration, fileSize int64) (res bool, err error)
}

type VoiceMsgValidator struct{}

func New() Validator {
	return &VoiceMsgValidator{}
}

func (v *VoiceMsgValidator) Validate(ctx context.Context, duration *durationpb.Duration, fileSize int64) (res bool, err error) {
	if fileSize > 2000 {
		return false, errors.New("invalid file size")
	}
	return true, nil
}
