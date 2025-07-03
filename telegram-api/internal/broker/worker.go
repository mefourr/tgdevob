package broker

import (
	"context"
	"encoding/json"
	"gihub.com/mefourr/tgdevob/telegram-api/internal/message"
	"github.com/IBM/sarama"
	"log/slog"
)

func Todo(ctx context.Context, msg *sarama.ConsumerMessage) {
	var m message.UserRequest
	err := json.Unmarshal(msg.Value, &m)
	if err != nil {
		return
	}
	slog.InfoContext(ctx, "Successfully consume a msg", "from", m.Request.From.UserName)
}
