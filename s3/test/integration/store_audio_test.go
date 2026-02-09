package integration

import (
	"github.com/mefourr/tgdevob/proto/s3_storage/pb/s3_storage/v1"
	"github.com/mefourr/tgdevob/s3/test/integration/suite"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	expected = "Bucket has been created"
)

func TestStoreAudio_HappyPath(t *testing.T) {
	ctx, s := suite.New(t)

	res, err := s.S3StorageClient.StoreAudio(ctx, &s3_storage.TempRq{
		Smth: "something here",
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, res.Smth)
}
