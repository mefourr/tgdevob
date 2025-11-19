package deserial

import (
	"encoding/json"
	"errors"
	"github.com/mefourr/tgdevob/worker/internal/kafka/message"
)

var ErrNoData = errors.New("no data in consumed message")

type MessageParser interface {
	Parse(res *message.Message, data []byte) error
}

type messageParser struct{}

func New() MessageParser {
	return &messageParser{}
}

func (mp *messageParser) Parse(res *message.Message, data []byte) error {
	if len(data) == 0 {
		return ErrNoData
	}
	return json.Unmarshal(data, res)
}
