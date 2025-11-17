package idem

import (
	"github.com/mefourr/tgdevob/worker/internal/kafka/message"
	"github.com/mefourr/tgdevob/worker/internal/storage/user/rediscache"
)

type IdempotencyChecker interface {
	Check(u *rediscache.User, msg message.Message) bool
}

type DefaultChecker struct{}

// Check returns true if tgbotapi.Message.MessageID mismatched
func (DefaultChecker) Check(u *rediscache.User, msg message.Message) bool {
	return u.LastRequest.MessageID != msg.Request.MessageID
}
