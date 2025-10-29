package deserial

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/worker/internal/kafka/message"
)

func ParseMessage(msg *sarama.ConsumerMessage) (message.Message, error) {
	var m message.Message
	err := json.Unmarshal(msg.Value, &m)
	if err != nil {
		return message.Message{}, err
	}
	return m, nil
}
