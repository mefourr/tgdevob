package integration

import (
	"github.com/mefourr/tgdevob/msg/voice/validator/test/integration/suite"
	"github.com/mefourr/tgdevob/proto/voice-msg-validator/pb/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	duration = 3
	expected = true
)

func TestValidVoiceMessage_HappyPath(t *testing.T) {
	ctx, s := suite.New(t)

	res, err := s.VVMLClient.Validate(ctx, &pb.VoiceMessageDataRq{
		Duration: float32(duration),
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, res.GetOk())
	assert.False(t, !res.GetOk())
}
