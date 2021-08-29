package telegram

import (
	"net/url"
	"strconv"
)

type Request interface {
	GetParams() (url.Values, string, error)
}

type UpdatesRequest struct {
	Offset         int      `json:"offset"`
	Limit          int      `json:"limit"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

func (req *UpdatesRequest) GetParams() (url.Values, string, error) {
	values := url.Values{}
	if req.Offset != 0 {
		values.Add("offset", strconv.Itoa(req.Offset))
	}
	if req.Limit > 0 {
		values.Add("limit", strconv.Itoa(req.Limit))
	}
	if req.Timeout > 0 {
		values.Add("timeout", strconv.Itoa(req.Timeout))
	}
	for _, au := range req.AllowedUpdates {
		values.Add("allowed_updates", au)
	}
	return values, "getUpdates", nil
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
