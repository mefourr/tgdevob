package cache

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type User struct {
	Id          string            `redis:"id"`
	TgUserId    int               `redis:"tg_user_id"`
	UserName    string            `redis:"tg_user_name"`
	FirstName   string            `redis:"tg_user_first_name"`
	LastName    string            `redis:"tg_user_last_name"`
	UpdateId    int               `redis:"update_id"`
	LastRequest *tgbotapi.Message `redis:"last_request"`
}
