package worker

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Message struct {
	UpdateId int               `json:"update_id"`
	User     *tgbotapi.User    `json:"sent_from"`
	Request  *tgbotapi.Message `json:"worker"`
}
