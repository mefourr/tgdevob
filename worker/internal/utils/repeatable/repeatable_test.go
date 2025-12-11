package repeatable

import (
	"errors"
	"github.com/mefourr/tgdevob/worker/config"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConnectBase(t *testing.T) {
	cfg := &config.Config{
		Postgres: struct {
			Database string
			Username string
			Password string
			Host     string
			Port     string
			RetryNum int
			Delay    time.Duration
		}{
			Database: "",
			Username: "",
			Password: "",
			Host:     "",
			Port:     "",
			RetryNum: 3,
			Delay:    1 * time.Second,
		},
	}
	t.Run("base test", func(t *testing.T) {
		err := Connect(func() error {
			return nil
		}, cfg)
		require.Nil(t, err)
		require.NoError(t, err)
	})
}

func TestConnectError(t *testing.T) {
	cfg := &config.Config{
		Postgres: struct {
			Database string
			Username string
			Password string
			Host     string
			Port     string
			RetryNum int
			Delay    time.Duration
		}{
			Database: "",
			Username: "",
			Password: "",
			Host:     "",
			Port:     "",
			RetryNum: 3,
			Delay:    1 * time.Second,
		},
	}
	t.Run("error occurs", func(t *testing.T) {
		err := Connect(func() error {
			return errors.New("error")
		}, cfg)
		require.Error(t, err)
	})
}
