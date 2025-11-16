package integration

import (
	"github.com/mefourr/tgdevob/msg/voice/validator/test/integration/suite"
	"github.com/mefourr/tgdevob/proto/voice-msg-validator/pb/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	"testing"
)

const (
	duration = 3
	fileSize = 1024
	expected = false
)

func TestValidVoiceMessage_HappyPath(t *testing.T) {
	ctx, s := suite.New(t)

	res, err := s.VVMLClient.ValidateVMLength(ctx, &pb.VoiceMessageDataRq{
		FileSize: fileSize,
		Duration: duration,
		Du:       durationpb.New(duration),
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, res.GetIsValidated())
	assert.True(t, !res.GetIsValidated())
}
