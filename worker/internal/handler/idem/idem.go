package idem

import (
	"gihub.com/mefourr/tgdevob/worker/internal/kafka/message"
	"gihub.com/mefourr/tgdevob/worker/internal/storage/user/cache"
)

func Validate(u *cache.User, msg message.Message) bool {
	if u.LastRequest.MessageID != msg.Request.MessageID {
		u.LastRequest = msg.Request
		return true
	}
	return false
}
