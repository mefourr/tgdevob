package idem

import (
	"github.com/mefourr/tgdevob/worker/internal/kafka/message"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/cache"
)

// Validate returns true if tgbotapi.Message.MessageID mismatched
func Validate(u *cache.User, msg message.Message) bool {
	return u.LastRequest.MessageID != msg.Request.MessageID
}
