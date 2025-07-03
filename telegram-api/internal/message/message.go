package message

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Response struct {
	rs string
}

func NewResponse(rs string) *Response {
	return &Response{rs: rs}
}

type UserRequest struct {
	Request *tgbotapi.Message `json:"tg_request"`
}

func NewUserRequest(request *tgbotapi.Message) *UserRequest {
	return &UserRequest{Request: request}
}
