package deserial

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/worker/internal/kafka/message"
)

type MessageParser interface {
	Parse(res *message.Message, cm *sarama.ConsumerMessage) error
}

type messageParser struct{}

func New() MessageParser {
	return &messageParser{}
}

func (mp *messageParser) Parse(res *message.Message, cm *sarama.ConsumerMessage) error {
	return json.Unmarshal(cm.Value, res)
}
