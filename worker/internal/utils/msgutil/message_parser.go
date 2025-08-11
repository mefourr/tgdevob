package msgutil

import (
	"encoding/json"
	"gihub.com/mefourr/tgdevob/worker/internal/kafka/message"
	"github.com/IBM/sarama"
)

func Parse(msg *sarama.ConsumerMessage) (message.Message, error) {
	var m message.Message
	err := json.Unmarshal(msg.Value, &m)
	if err != nil {
		return message.Message{}, err
	}
	return m, nil
}
