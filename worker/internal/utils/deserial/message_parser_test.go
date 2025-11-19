package deserial

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mefourr/tgdevob/worker/internal/kafka/message"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_messageParser_Parse_BaseCase(t *testing.T) {
	payload, err := json.Marshal(struct {
		UpdateId int
		User     *tgbotapi.User
		Request  *tgbotapi.Message
	}{
		UpdateId: 1,
		User: &tgbotapi.User{
			ID:           1231231,
			IsBot:        false,
			FirstName:    "Ivan",
			LastName:     "Ivanov",
			UserName:     "@ivanIvanov",
			LanguageCode: "ru",
		},
		Request: &tgbotapi.Message{
			MessageID: 1,
			Text:      "",
			Audio:     nil,
			Document:  nil,
			Photo:     nil,
			Sticker:   nil,
			Video:     nil,
			Voice:     nil,
		},
	})
	require.NoError(t, err)

	type args struct {
		res  *message.Message
		data []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "parse broker message: base test",
			args: args{
				res:  &message.Message{},
				data: payload,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := &messageParser{}
			err := mp.Parse(tt.args.res, tt.args.data)
			require.NoError(t, err)
		})
	}
}

func Test_messageParser_Parse_Errors(t *testing.T) {
	type args struct {
		res  *message.Message
		data []byte
	}
	tests := []struct {
		name   string
		args   args
		expErr string
	}{
		{
			name: "no bytes",
			args: args{
				res:  &message.Message{},
				data: make([]byte, 0),
			},
			expErr: ErrNoData.Error(),
		},
		{
			name: "invalid bytes",
			args: args{
				res:  &message.Message{},
				data: []byte("invalid"),
			},
			expErr: "invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := New()
			err := mp.Parse(tt.args.res, tt.args.data)
			require.ErrorContains(t, err, tt.expErr)
		})
	}
}
