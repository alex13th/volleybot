package telegram

import (
	"net/url"
	"strconv"
)

type Request interface {
	GetParams() (url.Values, string, error)
}

type MessageRequest struct {
	ChatId                   int
	Text                     string
	ParseMode                string
	Entities                 []MessageEntity
	DisableWebPagePreview    bool
	DisableNotification      bool
	ReplyToMessageId         int
	AllowSendingWithoutReply bool
	ReplyMarkup              interface{}
}

func (req *MessageRequest) GetParams() (url.Values, string, error) {
	values := url.Values{}
	values.Add("chat_id", strconv.Itoa(req.ChatId))
	values.Add("text", req.Text)
	return values, "sendMessage", nil
}
