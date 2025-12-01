package service

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mefourr/tgdevob/worker/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_messageParser_Parse_BaseCase(t *testing.T) {
	expected := domain.Message{
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
	}
	payload, err := json.Marshal(expected)
	require.NoError(t, err)

	type args struct {
		res  *domain.Message
		data []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "parse broker message: base test",
			args: args{
				res:  &domain.Message{},
				data: payload,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParserSvc()
			actual, err := p.Execute(tt.args.data)
			require.NoError(t, err)
			require.Equal(t, expected.UpdateId, actual.UpdateId)
		})
	}
}

func Test_messageParser_Parse_Errors(t *testing.T) {
	type args struct {
		res  *domain.Message
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
				res:  &domain.Message{},
				data: make([]byte, 0),
			},
			expErr: ErrNoData.Error(),
		},
		{
			name: "invalid bytes",
			args: args{
				res:  &domain.Message{},
				data: []byte("invalid"),
			},
			expErr: "invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParserSvc()
			_, err := p.Execute(tt.args.data)
			require.ErrorContains(t, err, tt.expErr)
		})
	}
}
